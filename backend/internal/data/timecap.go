package data

import (
	"errors"
	"strconv"
	"strings"
)

var ErrInvalidTimeCapFormat = errors.New("invalid TimeCap format")

type TimeCap int32

func (tc *TimeCap) UnmarshalJSON(jsonValue []byte) error {
	unqoutedJSONValue, err := strconv.Unquote(string(jsonValue))
	if err != nil {
		return ErrInvalidTimeCapFormat
	}

	parts := strings.Split(unqoutedJSONValue, " ")

	if len(parts) != 2 || parts[1] != "mins" {
		return ErrInvalidTimeCapFormat
	}

	i, err := strconv.ParseInt(parts[0], 10, 32)
	if err != nil {
		return ErrInvalidTimeCapFormat
	}

	*tc = TimeCap(i)

	return nil
}
