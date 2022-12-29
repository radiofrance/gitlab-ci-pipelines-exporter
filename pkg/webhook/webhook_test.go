package webhook

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_processHandler(t *testing.T) {
	var status atomic.Bool
	dummy := func(event bool) error {
		status.Store(event)

		if event {
			return fmt.Errorf("dummy error")
		}
		return nil
	}

	tcases := map[string]struct {
		method string
		event  string

		expectedStatus     bool
		expectedReturnCode int
		expectedBody       string
	}{
		"event:false": {
			method:             http.MethodPost,
			event:              `false`,
			expectedStatus:     false,
			expectedReturnCode: http.StatusOK,
		},
		"event:true": {
			method:             http.MethodPost,
			event:              `true`,
			expectedStatus:     true,
			expectedReturnCode: http.StatusBadRequest,
			expectedBody:       `{"error":"dummy error"}`,
		},
		"event:invalid method": {
			method:             http.MethodGet,
			expectedReturnCode: http.StatusMethodNotAllowed,
		},
		"event:invalid event": {
			method:             http.MethodPost,
			event:              `not a json`,
			expectedReturnCode: http.StatusBadRequest,
			expectedBody:       `{"error":"invalid character 'o' in literal null (expecting 'u')"}`,
		},
	}

	for name, tcase := range tcases {
		t.Run(name, func(t *testing.T) {
			status.Store(false)
			request := httptest.NewRequest(tcase.method, "http://:0", bytes.NewBuffer([]byte(tcase.event)))
			writer := httptest.NewRecorder()

			processHandler(dummy).ServeHTTP(writer, request)

			assert.Equal(t, tcase.expectedStatus, status.Load())
			assert.Equal(t, tcase.expectedReturnCode, writer.Code)
			assert.Equal(t, tcase.expectedBody, writer.Body.String())
		})
	}
}
