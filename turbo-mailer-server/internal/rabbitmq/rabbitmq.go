package rabbitmq

import (
	"crypto/tls"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/spf13/viper"
)

func Dial() (*amqp.Connection, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	return amqp.DialConfig(viper.GetString("rabbitmq.url"), amqp.Config{
		Heartbeat: 10 * time.Second,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		Properties: amqp.Table{
			"connection_name": hostname,
		},
	})
}

func Channel(conn *amqp.Connection) (*amqp.Channel, error) {
	return conn.Channel()
}

func TaskQueueName() string {
	queue := viper.GetString("rabbitmq.queue")
	if queue == "" {
		queue = "task_queue"
	}
	return queue
}

func TaskChannel(conn *amqp.Connection) (*amqp.Channel, string, error) {
	channel, err := conn.Channel()
	if err != nil {
		return nil, "", err
	}

	queue := TaskQueueName()

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
