package rabbitmq

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/spf13/viper"
)

func Dial() (*amqp.Connection, error) {
	return amqp.Dial(viper.GetString("rabbitmq.url"))
}

func Channel(conn *amqp.Connection) (*amqp.Channel, error) {
	return conn.Channel()
}

func TaskChannel(conn *amqp.Connection) (*amqp.Channel, string, error) {
	channel, err := conn.Channel()
	if err != nil {
		return nil, "", err
	}

	queue := viper.GetString("rabbitmq.queue")
	if queue == "" {
		queue = "task_queue"
	}

	_, err = channel.QueueDeclare(
		queue,
		true,
		false,
		false,
		false,
		amqp.Table{
			amqp.QueueMessageTTLArg: 1000 * 60 * 60 * 24, // 24 hours
		},
	)

	if err != nil {
		return nil, "", err
	}

	return channel, queue, nil
}
