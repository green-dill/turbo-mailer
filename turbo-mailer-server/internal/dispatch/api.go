package dispatch

import (
	"context"
	"turbo-mailer-server/internal/models"
	"turbo-mailer-server/internal/query"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog/log"
)

// DispatchImmediately dispatch task immediately
// for api call to trigger dispatch immediately
func DispatchImmediately(ctx context.Context, task *models.Task) error {
	d, conn, channel, err := buildApiDispatcher()
	if err != nil {
		return err
	}
	defer d.closeChannel(conn, channel)
	return d.dispatchTask(ctx, task, channel)
}

// TestSend is used for testing
func TestSend(ctx context.Context, taskID uint, to string) error {
	task, err := query.Task.WithContext(ctx).Where(query.Task.ID.Eq(taskID)).First()
	if err != nil {
		return err
	}
	d, conn, channel, err := buildApiDispatcher()
	if err != nil {
		return err
	}

	email, sender, err := d.buildEmail(ctx, task, to)
	if err != nil {
		return err
	}

	log.Info().Msgf("test send email from %s to %s", sender.FromEmail, to)

	err = d.send(channel, email)
	if err != nil {
		return err
	}

	defer d.closeChannel(conn, channel)
	return err
}

func buildApiDispatcher() (*dispatcher, *amqp.Connection, *amqp.Channel, error) {
	d, err := newDispatcher()
	if err != nil {
		return nil, nil, nil, err
	}
	conn, channel, err := d.createChannel()
	if err != nil {
		d.closeChannel(conn, channel)
		return nil, nil, nil, err
	}
	return d, conn, channel, nil
}
