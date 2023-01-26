package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/urfave/negroni"
	"go.uber.org/zap"
)

// NewZapMiddleware creates a negroni middleware to log every request.
func NewZapMiddleware(logger *zap.Logger) negroni.HandlerFunc {
	logger = logger.WithOptions(zap.WithCaller(false))
	return func(writer http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
		writer = negroni.NewResponseWriter(writer)

		start := time.Now()
		logger := logger.With(
			zap.String("http_referer", req.Referer()),
			zap.String("http_user_agent", req.UserAgent()),
			zap.String("remote_addr", req.RemoteAddr),
			zap.String("remote_user", req.Header.Get("X-Remote-User")),
		)

		next(writer, req)

		negroniWriter, hasStatus := writer.(negroni.ResponseWriter)
		if hasStatus {
			logger = logger.With(
				zap.Int("status", negroniWriter.Status()),
				zap.Int("body_bytes_sent", negroniWriter.Size()),
			)
		}

		logger = logger.With(
			zap.Float32("duration", float32(time.Since(start))/float32(time.Second)),
		)

		if !hasStatus {
			logger.Info(req.URL.String())
			return
		}

		switch {
		case negroniWriter.Status() >= http.StatusInternalServerError:
			logger.Error(req.URL.String())
		case negroniWriter.Status() >= http.StatusBadRequest:
			logger.Warn(req.URL.String())
		default:
			logger.Info(req.URL.String())
		}
	}
}

// NewRecoverMiddleware recovers from any panic and notify Gitlab through a 400.
func NewRecoverMiddleware(logger *zap.Logger) negroni.HandlerFunc {
	logger = logger.WithOptions(zap.AddStacktrace(zap.ErrorLevel))

	return func(writer http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
		defer func() {
			if err := recover(); err != nil {
				writer.WriteHeader(http.StatusInternalServerError)
				WriteError(writer, err)

				logger.Error("panic recovered", zap.String("error", fmt.Sprintf("%s", err)))
			}
		}()

		next(writer, req)
	}
}
