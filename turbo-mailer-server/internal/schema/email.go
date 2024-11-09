package schema

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

type Email struct {
	LogID       uint     `json:"log_id"`
	Subject     string   `json:"subject"`
	From        string   `json:"from"`
	Domain      string   `json:"domain"`
	Content     string   `json:"content"`
	ContentType string   `json:"content_type"`
	Receivers   []string `json:"receivers"`
}

func (e *Email) ToJSON() (string, error) {
	b, err := e.ToJSONBytes()
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (e *Email) ToJSONBytes() ([]byte, error) {
	return json.Marshal(e)
}

func (e *Email) GetHashKey() string {
	sha := sha256.New()
	sha.Write([]byte(e.From))
	sha.Write([]byte(e.Domain))
	sha.Write([]byte(e.Subject))
	sha.Write([]byte(e.Content))
	sha.Write([]byte(e.ContentType))
	for _, receiver := range e.Receivers {
		sha.Write([]byte(receiver))
	}
	return hex.EncodeToString(sha.Sum(nil))
}
