package schema

import "encoding/json"

type Email struct {
	Subject     string   `json:"subject"`
	From        string   `json:"from"`
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
