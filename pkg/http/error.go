package http

import (
	"encoding/json"
	"fmt"
	"io"
)

func WriteError(writer io.Writer, err any) {
	_, _ = writer.Write(errToJSON(err))
}

// errToJSON converts anything into json object representing an error.
func errToJSON(err any) []byte {
	bytes, _ := json.Marshal(map[string]string{"error": fmt.Sprintf("%s", err)}) //nolint:errchkjson
	return bytes
}
