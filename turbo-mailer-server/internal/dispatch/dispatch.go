package dispatch

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
	"turbo-mailer-server/internal/models"
	"turbo-mailer-server/internal/query"
	"turbo-mailer-server/internal/utils/ptr"

	"github.com/bsm/redislock"
	"github.com/labstack/echo"
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
	rows, err := query.DB.WithContext(ctx).
		Where(query.Task.State.Eq(models.TaskStatePending)).
		Order(query.Task.CreatedAt.Desc()).
		Limit(1000).
		Rows()
	if err != nil {
		log.Error().Err(err).Msg("failed to get tasks")
		return err
	}

	defer rows.Close()

	conn, err := amqp.Dial(viper.GetString("rabbitmq.url"))
	if err != nil {
		return err
	}
	defer conn.Close()

	channel, err := conn.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

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
		return err
	}

	for rows.Next() {
		var task models.Task
		if err := rows.Scan(&task); err != nil {
			log.Error().Err(err).Msg("failed to scan task")
			return err
		}

		if err := d.dispatchTask(ctx, &task, channel); err != nil {
			log.Error().Err(err).Msgf("failed to dispatch task %d", task.ID)
			return err
		}
	}

	return nil
}

func (d *dispatcher) dispatchTask(ctx context.Context, task *models.Task, channel *amqp.Channel) error {
	// redis distributed lock
	lockKey := query.Redis.Key("lock", "task", strconv.FormatUint(uint64(task.ID), 120))
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

	remaining := task.MaxDispatchPreHour

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
	}

	count := 0
	for _, receiver := range task.Receivers {
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

		subject := task.Subject
		content := task.Content

		poolID, err := d.selectPool(task.Pools)
		if err != nil {
			log.Error().Err(err).Msgf("failed to select pool for task %d", task.ID)
			return err
		}

		sender, err := d.selectSender(ctx, poolID)
		if err != nil {
			log.Error().Err(err).Msgf("failed to select sender for task %d", task.ID)
			return err
		}

		subject = strings.Replace(subject, "{name}", receiver, -1)
		content = strings.Replace(content, "{name}", receiver, -1)

		// todo template render

		from := fmt.Sprintf("%s <%s>", sender.FromName, sender.FromEmail)

		queue := viper.GetString("rabbitmq.queue")

		// send to rabbitmq
		err = channel.Publish(
			"",
			queue,
			false,
			false,
			amqp.Publishing{
				DeliveryMode: amqp.Persistent,
				ContentType:  echo.MIMEApplicationJSON,
				Body:         []byte(fmt.Sprintf("Subject: %s\nFrom: %s\n\n%s", subject, from, content)),
			},
		)
		if err != nil {
			log.Error().Err(err).Msgf("failed to publish task %d", task.ID)
			return err
		}

		// save task log
		query.DB.WithContext(ctx).Create(&models.TaskLog{
			PoolID:    poolID,
			Sender:    sender.FromEmail,
			Receiver:  receiver,
			State:     models.TaskLogStatePending,
			CreatedAt: time.Now(),
		})

		count++
	}

	return nil
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

func (d *dispatcher) selectSender(ctx context.Context, poolID uint) (sender models.PoolSender, err error) {
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
