package api

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

func MarschalWithIndent(data any, indent string) ([]byte, error) {
	var b []byte
	var marschalErr error

	if i, err := strconv.Atoi(indent); err == nil && i > 0 && i < 10 {
		b, marschalErr = json.MarshalIndent(data, "", strings.Repeat(" ", i))
	} else {
		b, marschalErr = json.Marshal(data)
	}

	if marschalErr != nil {
		return nil, fmt.Errorf("unable to marschal data: %w", marschalErr)
	}
	return b, nil
}