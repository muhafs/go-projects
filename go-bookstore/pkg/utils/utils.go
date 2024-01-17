package utils

import (
	"encoding/json"
	"io"
	"net/http"
)

func ParseBody(r *http.Request, data interface{}) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal([]byte(body), data)
	if err != nil {
		return
	}
}
