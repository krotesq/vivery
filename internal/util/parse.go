package util

import (
	"encoding/json"
	"errors"
	"io"
)

const maxBodySize = 1 << 20 // 1 MB

// ParseBody decodes a JSON body into dst.
func ParseBody[T any](body io.Reader, dst *T) error {
	if body == nil {
		return errors.New("request body is empty")
	}
	limited := io.LimitReader(body, maxBodySize)
	dec := json.NewDecoder(limited)
	dec.DisallowUnknownFields()
	return dec.Decode(dst)
}