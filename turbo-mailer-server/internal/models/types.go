package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type JSON[T any] struct {
	V T
}

func (j *JSON[T]) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	bytes, ok := value.(string)
	if !ok {
		return errors.New("failed to unmarshal JSON value")
	}

	return json.Unmarshal([]byte(bytes), &j.V)
}

func (j JSON[T]) Value() (driver.Value, error) {
	return json.Marshal(j.V)
}

func (j JSON[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(j.V)
}

func (j *JSON[T]) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &j.V)
}

func NewJSON[T any](v T) JSON[T] {
	return JSON[T]{V: v}
}

func (j JSON[T]) Val() T {
	return j.V
}

type StringArray = JSON[[]string]
