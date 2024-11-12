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
	smtpHost := fmt.Sprintf("postfix-%s", hashStr[:8])
	smtpPort := 587

	from := fmt.Sprintf("noreply@%s", domain)
	receiver := "yinheli+test@gmail.com"

	message := gomail.NewMessage()
	message.SetAddressHeader("From", from, "noreply")
	message.SetAddressHeader("To", receiver, "")
	message.SetHeader("Subject", "李白 将进酒")
	message.SetBody("text/plain", `君不见黄河之水天上来，奔流到海不复回。
君不见高堂明镜悲白发，朝如青丝暮成雪。`)

	smtpHost = "127.0.0.1"
	smtpPort = 2525

	dialer := gomail.NewDialer(smtpHost, smtpPort, "", "")
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	err := dialer.DialAndSend(message)
	if err != nil {
		t.Fatal(err)
	}
}

func TestSend_via_postfix_dkim(t *testing.T) {
	domain := "mail.turbomx.org"
	from := fmt.Sprintf("noreply@%s", domain)
	// https://www.appmaildev.com/cn/dkim
	receiver := "test-a44a1609@appmaildev.com"

	//
	// receiver := "ping@tools.mxtoolbox.com"

	message := gomail.NewMessage()
	message.SetAddressHeader("From", from, "noreply")
	message.SetAddressHeader("To", receiver, "")
	message.SetHeader("Subject", "李白 将进酒")
	message.SetBody("text/plain", `君不见黄河之水天上来，奔流到海不复回。
君不见高堂明镜悲白发，朝如青丝暮成雪。`)

	smtpHost := "127.0.0.1"
	smtpPort := 2525

	dialer := gomail.NewDialer(smtpHost, smtpPort, "", "")
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	err := dialer.DialAndSend(message)
	if err != nil {
		t.Fatal(err)
	}
}
