package utils

import (
	"encoding/json"
	"io"
	"strings"
)

func JsonStringReaderFor(o interface{}) io.Reader {
	bytes, _ := json.Marshal(o)
	return strings.NewReader(string(bytes))
}
