package http_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	gcpehttp "github.com/radiofrance/gitlab-ci-pipelines-exporter/pkg/http"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestNewZapMiddleware(t *testing.T) {
	buffer := bytes.NewBuffer(nil)
	logger := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(zapcore.EncoderConfig{
			LevelKey:       "L",
			NameKey:        "N",
			CallerKey:      "C",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "M",
			StacktraceKey:  "S",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeTime:     zapcore.RFC3339NanoTimeEncoder,
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}),
		zapcore.AddSync(buffer),
		zap.InfoLevel,
	))

	next := http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) { writer.WriteHeader(http.StatusOK) })
	mw := gcpehttp.NewZapMiddleware(logger)

	req, _ := http.NewRequest(http.MethodPost, "https://:::0", nil)
	req.Header.Set("Referer", "go-test")
	req.Header.Set("User-Agent", "go-test")
	req.Header.Set("X-Remote-User", "go-test")
	req.RemoteAddr = "127.0.0.1"
	w := httptest.NewRecorder()
	mw.ServeHTTP(w, req, next)
	assert.Equal(t, http.StatusOK, w.Code)

	var log struct {
		Level   string `json:"L"`
		Message string `json:"M"`

		BodyBytesSent int     `json:"body_bytes_sent"`
		Duration      float32 `json:"duration"`
		HttpReferer   string  `json:"http_referer"`
		HttpUserAgent string  `json:"http_user_agent"`
		RemoteAddr    string  `json:"remote_addr"`
		RemoteUser    string  `json:"remote_user"`
		Status        int     `json:"status"`
		TimeLocal     string  `json:"time_local"`
	}
	require.NoError(t, json.Unmarshal(buffer.Bytes(), &log))
	assert.Equal(t, "INFO", log.Level)
	assert.Equal(t, "https://:::0", log.Message)

	assert.Equal(t, 0, log.BodyBytesSent)
	assert.Less(t, log.Duration*float32(time.Second), float32(time.Millisecond))
	assert.Equal(t, "go-test", log.HttpReferer)
	assert.Equal(t, "go-test", log.HttpUserAgent)
	assert.Equal(t, "127.0.0.1", log.RemoteAddr)
	assert.Equal(t, "go-test", log.RemoteUser)
	assert.Equal(t, http.StatusOK, log.Status)
}

func TestNewRecoverMiddleware(t *testing.T) {
	buffer := bytes.NewBuffer(nil)
	logger := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(zapcore.EncoderConfig{
			LevelKey:       "L",
			NameKey:        "N",
			CallerKey:      "C",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "M",
			StacktraceKey:  "S",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}),
		zapcore.AddSync(buffer),
		zap.InfoLevel,
	))

	next := http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) { panic("STONK !!!") })
	mw := gcpehttp.NewRecoverMiddleware(logger)

	req, _ := http.NewRequest(http.MethodPost, "https://:::0", nil)
	w := httptest.NewRecorder()
	mw.ServeHTTP(w, req, next)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, `{"error":"STONK !!!"}`, w.Body.String())

	var log struct {
		Level      string `json:"L"`
		Message    string `json:"M"`
		StackTrace string `json:"S"`
		Error      string `json:"error"`
	}
	require.NoError(t, json.Unmarshal(buffer.Bytes(), &log))
	assert.Equal(t, "ERROR", log.Level)
	assert.Equal(t, "panic recovered", log.Message)
	assert.Equal(t, "STONK !!!", log.Error)
	assert.NotEmpty(t, log.StackTrace)
}
