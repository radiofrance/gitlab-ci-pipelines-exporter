package http

import (
	"encoding/json"
	"fmt"
	"io"
)

func WriteError(writer io.Writer, err any) {
	_, _ = writer.Write(errToJson(err))
}

// errToJson convert anything into json object representing an error.
func errToJson(err any) []byte {
	bytes, _ := json.Marshal(map[string]string{"error": fmt.Sprintf("%s", err)})
	return bytes
}
