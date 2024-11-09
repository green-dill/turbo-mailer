package worker

import (
	"context"
	"crypto/sha1"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
	"turbo-mailer-server/internal/query"
	"turbo-mailer-server/internal/rabbitmq"
	"turbo-mailer-server/internal/schema"

	"github.com/rs/zerolog/log"
	gomail "gopkg.in/mail.v2"
)

func run(ctx context.Context) {
	defer func() {
		if r := recover(); r != nil {
			log.Error().Msgf("panic: %v", r)
			time.Sleep(time.Second * 10)
			run(ctx)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			consumeAndSend()
		}

		time.Sleep(time.Second * 10)
	}
}

func consumeAndSend() {
	conn, err := rabbitmq.Dial()
	if err != nil {
		log.Error().Err(err).Msg("failed to dial")
		return
	}

	channel, queue, err := rabbitmq.TaskChannel(conn)
	if err != nil {
		log.Error().Err(err).Msg("failed to create channel")
		return
	}

	hostname, err := os.Hostname()
	if err != nil {
		log.Error().Err(err).Msg("failed to get hostname")
		return
	}

	deliveries, err := channel.Consume(
		queue,
		hostname,
		false,
		false,
		false,
		false,
		nil,
	)

	for delivery := range deliveries {
		var email schema.Email
		err := json.Unmarshal(delivery.Body, &email)
		if err != nil {
			log.Error().Err(err).Msg("failed to unmarshal email")
			delivery.Reject(false)
			continue
		}

		key := query.Redis.Key(email.GetHashKey())
		r := query.Redis.Incr(context.Background(), key)
		if r.Val() > 10 {
			log.Error().Msg("send too many times")
			delivery.Reject(false)
			continue
		}

		success, err := send(&email)
		if err != nil {
			log.Error().Err(err).Msg("failed to send email")
			delivery.Nack(false, true)
			continue
		}

		if success {
			delivery.Ack(false)
		} else {
			delivery.Nack(false, true)
		}
	}
}

func send(email *schema.Email) (bool, error) {
	message := gomail.NewMessage()
	// parse from like `"test@example.com" <test@example.com>`
	fromName := email.From
	fromEmail := email.From
	if strings.Contains(email.From, "<") {
		fromParts := strings.Split(email.From, "<")
		fromName = strings.TrimSpace(fromParts[0])
		fromName = strings.Trim(fromName, "\"")
		fromEmail = strings.Trim(fromParts[1], ">")
	}
	message.SetAddressHeader("From", fromEmail, fromName)

	for _, receiver := range email.Receivers {
		message.SetAddressHeader("To", receiver, receiver)
	}
	message.SetHeader("Subject", email.Subject)
	message.SetBody(email.ContentType, email.Content)

	hasher := sha1.New()
	hasher.Write([]byte(email.Domain))
	hash := hex.EncodeToString(hasher.Sum(nil))
	smtpHost := fmt.Sprintf("postfix.%s", hash[:8])

	log.Debug().Str("domain", email.Domain).Str("smtp_host", smtpHost).Msg("send email")

	dialer := gomail.NewDialer(smtpHost, 25, "", "")
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	err := dialer.DialAndSend(message)
	if err != nil {
		return false, err
	}

	return true, nil
}
