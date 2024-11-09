package dispatch

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"time"
	"turbo-mailer-server/internal/models"
	"turbo-mailer-server/internal/query"
	"turbo-mailer-server/internal/render"
	"turbo-mailer-server/internal/schema"
	"turbo-mailer-server/internal/utils/ptr"

	"github.com/bsm/redislock"
	"github.com/labstack/echo/v4"
	cmap "github.com/orcaman/concurrent-map/v2"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type dispatcher struct {
	poolCache *cmap.ConcurrentMap[uint, *models.Pool]
}

func newDispatcher() (*dispatcher, error) {
	poolCache := cmap.NewWithCustomShardingFunction[uint, *models.Pool](func(key uint) uint32 {
		return uint32(key)
	})

	return &dispatcher{
		poolCache: &poolCache,
	}, nil
}

func (d *dispatcher) Run(ctx context.Context) {
	interval := viper.GetDuration("dispatch.interval")
	if interval == 0 {
		interval = time.Minute
	}

	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(interval):
			log.Info().Msg("trigger dispatch tasks")
			if err := d.processTasks(ctx); err != nil {
			}
		}
	}
}

func (d *dispatcher) processTasks(ctx context.Context) error {
	tasks, err := query.Task.WithContext(ctx).
		Preload(query.Task.Pools).
		Where(query.Task.State.Eq(models.TaskStatePending)).
		Where(query.Task.ScheduleAt.Gte(time.Now().Add(-time.Hour * 24))).
		Order(query.Task.CreatedAt.Asc()).
		Limit(100).Find()

	conn, channel, err := d.createChannel()
	if err != nil {
		log.Error().Err(err).Msg("failed to create channel")
		return err
	}
	defer d.closeChannel(conn, channel)

	for _, task := range tasks {
		if err := d.dispatchTask(ctx, task, channel); err != nil {
			log.Error().Err(err).Msgf("failed to dispatch task %d", task.ID)
			return err
		}
	}

	return nil
}

func (d *dispatcher) dispatchTask(ctx context.Context, task *models.Task, channel *amqp.Channel) error {
	lockKey := query.Redis.Key("lock", "task", fmt.Sprintf("%d", task.ID))
	lock, err := redislock.New(query.Redis.Client).Obtain(ctx, lockKey, 10*time.Second, nil)
	if err != nil && errors.Is(err, redislock.ErrNotObtained) {
		log.Info().Msgf("failed to obtain lock for task %d", task.ID)
		return nil
	}
	defer lock.Release(context.Background())

	if err != nil {
		log.Error().Err(err).Msgf("failed to create lock for task %d", task.ID)
		return err
	}

	lastDispatchAt := task.LastDispatchAt

	defer func() {
		task.LastDispatchAt = ptr.Ptr(time.Now())
		if err := query.DB.WithContext(ctx).Save(task).Error; err != nil {
			log.Error().Err(err).Msgf("failed to save task %d", task.ID)
		}
	}()

	var remaining int

	if lastDispatchAt != nil && task.MaxDispatchPreHour > 0 {
		var count int64
		err := query.DB.WithContext(ctx).
			Where(query.TaskLog.TaskID.Eq(task.ID)).
			Where(query.TaskLog.CreatedAt.Gte(time.Now().Add(-time.Hour))).
			Count(&count).Error

		if err != nil {
			log.Error().Err(err).Msgf("failed to count task log for task %d", task.ID)
			return err
		}

		if count >= int64(task.MaxDispatchPreHour) {
			log.Info().Msgf("task %d reached max dispatch pre hour", task.ID)
			return nil
		}

		remaining = task.MaxDispatchPreHour - int(count)

		log.Info().
			Uint("task_id", task.ID).
			Int("remaining", remaining).
			Msg("task check remaining for limit dispatch policy")
	} else {
		if task.MaxDispatchPreHour > 0 {
			remaining = task.MaxDispatchPreHour
		} else {
			remaining = math.MaxInt
		}
	}

	count := 0
	for _, receiver := range task.Receivers.Val() {
		var n int64
		err := query.DB.WithContext(ctx).
			Where(query.TaskLog.TaskID.Eq(task.ID)).
			Where(query.TaskLog.Receiver.Eq(receiver)).
			Count(&n).Error
		if err != nil {
			log.Error().Err(err).Msgf("failed to count task log for task %d", task.ID)
			return err
		}
		if n > 0 {
			// skip if already dispatched
			continue
		}

		log.Info().
			Uint("task_id", task.ID).
			Str("receiver", receiver).
			Int("remaining", remaining-count).
			Msg("dispatch task")

		if count >= remaining {
			log.Info().Msgf("task %d reached max dispatch pre hour", task.ID)
			return nil
		}

		email, sender, err := d.buildEmail(ctx, task, receiver)
		if err != nil {
			log.Error().Err(err).Msgf("failed to build email for task %d", task.ID)
			return err
		}

		if err := d.send(channel, email); err != nil {
			log.Error().Err(err).Msgf("failed to send email for task %d", task.ID)
			return err
		}

		// save task log
		if err := d.saveTaskLog(ctx, task.ID, sender.PoolID, sender.FromEmail, receiver, models.TaskLogStatePending, ""); err != nil {
			log.Error().Err(err).Msgf("failed to save task log for task %d", task.ID)
			return err
		}
		count++
	}

	return nil
}

func (d *dispatcher) buildEmail(ctx context.Context, task *models.Task, receivers ...string) (*schema.Email, *models.PoolSender, error) {
	sender, err := d.selectSenderByTask(ctx, task)
	if err != nil {
		log.Error().Err(err).Msgf("failed to select sender for task %d", task.ID)
		return nil, nil, err
	}

	from := fmt.Sprintf("%s <%s>", sender.FromName, sender.FromEmail)

	// email template parameters
	params := map[string]any{
		"receivers": receivers,
		"sender":    sender,
	}

	if task.Metadata != nil {
		metadata := make(map[string]any)
		if err := json.Unmarshal(*task.Metadata, &metadata); err != nil {
			log.Error().Err(err).Msgf("failed to unmarshal metadata for task %d", task.ID)
			return nil, nil, err
		}
		params["metadata"] = metadata
	}

	subject, err := render.Render(task.Subject, params)
	if err != nil {
		log.Error().Err(err).Msgf("failed to render email subject for task %d", task.ID)
		return nil, nil, err
	}

	content, err := render.Render(task.Content, params)
	if err != nil {
		log.Error().Err(err).Msgf("failed to render email content for task %d", task.ID)
		return nil, nil, err
	}

	email := &schema.Email{
		Subject:     subject,
		From:        from,
		Content:     content,
		ContentType: task.ContentType,
		Receivers:   receivers,
	}

	return email, sender, nil
}

func (d *dispatcher) selectSenderByTask(ctx context.Context, task *models.Task) (*models.PoolSender, error) {
	poolID, err := d.selectPool(task.Pools)
	if err != nil {
		log.Error().Err(err).Msgf("failed to select pool for task %d", task.ID)
		return nil, err
	}

	sender, err := d.selectSender(ctx, poolID)
	if err != nil {
		log.Error().Err(err).Msgf("failed to select sender for task %d", task.ID)
		return nil, err
	}

	return sender, nil
}

func (d *dispatcher) saveTaskLog(ctx context.Context, taskID uint, poolID uint, sender string, receiver string, state string, message string) error {
	return query.DB.WithContext(ctx).Create(&models.TaskLog{
		TaskID:    taskID,
		PoolID:    poolID,
		Sender:    sender,
		Receiver:  receiver,
		State:     state,
		Message:   message,
		CreatedAt: time.Now(),
	}).Error
}

func (d *dispatcher) send(channel *amqp.Channel, email *schema.Email) error {
	body, err := email.ToJSONBytes()
	if err != nil {
		return err
	}

	return channel.Publish(
		"",
		d.queueName(),
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  echo.MIMEApplicationJSON,
			Body:         body,
		},
	)

}

func (d *dispatcher) selectPool(pools []*models.TaskPool) (poolID uint, err error) {
	totalWeight := 0
	for _, pool := range pools {
		if pool.Weight <= 0 {
			pool.Weight = 1
		}
		totalWeight += pool.Weight
	}

	if totalWeight == 0 {
		return 0, fmt.Errorf("total weight is zero")
	}

	randomWeight := rand.Intn(totalWeight)
	cumulativeWeight := 0

	for _, pool := range pools {
		cumulativeWeight += pool.Weight
		if randomWeight < cumulativeWeight {
			return pool.PoolID, nil
		}
	}

	return 0, fmt.Errorf("no pool selected")
}

func (d *dispatcher) selectSender(ctx context.Context, poolID uint) (sender *models.PoolSender, err error) {
	var pool *models.Pool

	if p, ok := d.poolCache.Get(poolID); ok {
		pool = p
	} else {
		pool, err = query.Q.WithContext(ctx).Pool.Where(query.Pool.ID.Eq(poolID)).First()
		if err != nil {
			log.Error().Err(err).Msgf("failed to get pool %d", poolID)
			return
		}
		d.poolCache.Set(poolID, pool)
	}

	senders := pool.Senders
	if len(senders) == 0 {
		err = fmt.Errorf("no sender found for pool %d", poolID)
		return
	}

	randomIndex := rand.Intn(len(senders))
	sender = senders[randomIndex]

	return
}

func (d *dispatcher) queueName() string {
	queue := viper.GetString("rabbitmq.queue")
	if queue == "" {
		queue = "task_queue"
	}
	return queue
}

func (d *dispatcher) createChannel() (*amqp.Connection, *amqp.Channel, error) {
	conn, err := amqp.Dial(viper.GetString("rabbitmq.url"))
	if err != nil {
		return nil, nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		d.closeChannel(conn, nil)
		return nil, nil, err
	}

	queue := d.queueName()
	_, err = channel.QueueDeclare(
		queue,
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Error().Err(err).Msg("failed to declare queue")
		d.closeChannel(conn, channel)
		return nil, nil, err
	}

	return conn, channel, nil
}

func (d *dispatcher) closeChannel(conn *amqp.Connection, channel *amqp.Channel) {
	if conn != nil {
		_ = conn.Close()
	}
	if channel != nil {
		_ = channel.Close()
	}
}
