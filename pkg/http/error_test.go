package http_test

import (
	"bytes"
	"net"
	"testing"

	"github.com/radiofrance/gitlab-ci-pipelines-exporter/pkg/http"
	"github.com/stretchr/testify/assert"
)

func Test_WriteError(t *testing.T) {
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
		var buf bytes.Buffer
		http.WriteError(&buf, tcase.err)

		assert.Equal(t, tcase.expected, buf.String())
	}
}
