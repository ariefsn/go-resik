package helper

import (
	"encoding/json"
	"net/http"
)

func ParsePayload(r *http.Request, target interface{}) error {
	return json.NewDecoder(r.Body).Decode(&target)
}
