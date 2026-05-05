package util

import (
	"encoding/json"
	"errors"
	"io"
)

// ParseBody decodes a JSON body into dst.
func ParseBody[T any](body io.Reader, dst *T) error {
	if body == nil {
		return errors.New("request body is empty")
	}
	dec := json.NewDecoder(body)
	dec.DisallowUnknownFields()
	return dec.Decode(dst)
}