package eml

import (
	"io"

	"github.com/jhillyerd/enmime"
)

type Email struct {
	From        string
	To          string
	Subject     string
	ContentType string
	Body        string
}

func Parse(eml io.Reader) (*Email, error) {
	envelope, err := enmime.ReadEnvelope(eml)
	if err != nil {
		return nil, err
	}

	email := &Email{
		From:        envelope.GetHeader("From"),
		To:          envelope.GetHeader("To"),
		Subject:     envelope.GetHeader("Subject"),
		ContentType: "text/html",
		Body:        envelope.HTML,
	}

	return email, nil
}
