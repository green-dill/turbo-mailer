package worker

import (
	"crypto/sha1"
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"testing"

	gomail "gopkg.in/mail.v2"
)

func TestSend_via_postfix(t *testing.T) {
	domain := "mail.turbomx.org"
	hash := sha1.New()
	hash.Write([]byte(domain))
	hashStr := hex.EncodeToString(hash.Sum(nil))
	smtpHost := fmt.Sprintf("postfix.%s", hashStr[:8])
	smtpPort := 587

	from := fmt.Sprintf("test@%s", domain)
	receiver := fmt.Sprintf("test@%s", domain)

	message := gomail.NewMessage()
	message.SetAddressHeader("From", from, from)
	message.SetAddressHeader("To", receiver, receiver)
	message.SetHeader("Subject", "test")
	message.SetBody("text/plain", "test")

	smtpHost = "127.0.0.1"
	smtpPort = 2525

	dialer := gomail.NewDialer(smtpHost, smtpPort, "", "")
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	err := dialer.DialAndSend(message)
	if err != nil {
		t.Fatal(err)
	}
}
