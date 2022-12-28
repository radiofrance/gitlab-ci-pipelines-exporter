package webhook

import (
	"encoding/json"
	"fmt"
)

// errToJson convert anything into json object representing an error.
func errToJson(err any) []byte {
	bytes, _ := json.Marshal(map[string]string{"error": fmt.Sprintf("%s", err)})
	return bytes
}
