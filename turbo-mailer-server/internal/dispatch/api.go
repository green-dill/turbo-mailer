package dispatch

import (
	"context"
	"turbo-mailer-server/internal/models"
	"turbo-mailer-server/internal/schema"

	"github.com/openzipkin/zipkin-go/reporter/amqp"
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
func TestSend(ctx context.Context, email *schema.Email) error {
	d, conn, channel, err := buildApiDispatcher()
	if err != nil {
		return err
	}
	defer d.closeChannel(conn, channel)
	return d.send(channel, email)
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
