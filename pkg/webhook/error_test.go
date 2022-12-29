package webhook

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_errToJson(t *testing.T) {
	tcases := []struct {
		err      any
		expected string
	}{
		{"error", `{"error":"error"}`},
		{false, `{"error":"%!s(bool=false)"}`},
		{net.ErrClosed, `{"error":"use of closed network connection"}`},
		{nil, `{"error":"%!s(\u003cnil\u003e)"}`},
	}

	for _, tcase := range tcases {
		assert.Equal(t, tcase.expected, string(errToJson(tcase.err)))
	}
}
