package webhook

import (
	"fmt"
	"net/http"
	"time"

	"github.com/urfave/negroni"
	"go.uber.org/zap"
)

// NewZapMiddleware creates a negroni middleware to log every request.
func NewZapMiddleware(logger *zap.Logger) negroni.Handler {
	logger = logger.WithOptions(zap.WithCaller(false))
	return negroni.HandlerFunc(func(writer http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
		writer = negroni.NewResponseWriter(writer)

		start := time.Now()
		logger = logger.With(
			zap.String("http_referer", req.Referer()),
			zap.String("http_user_agent", req.UserAgent()),
			zap.String("remote_addr", req.RemoteAddr),
			zap.String("remote_user", req.Header.Get("X-Remote-User")),
		)

		next(writer, req)

		logger = logger.With(
			zap.Int("status", writer.(negroni.ResponseWriter).Status()),
			zap.Int("body_bytes_sent", writer.(negroni.ResponseWriter).Size()),
			zap.Float32("duration", float32(time.Since(start))/float32(time.Second)),
		)

		switch {
		case writer.(negroni.ResponseWriter).Status() >= 500:
			logger.Error(req.URL.String())
		case writer.(negroni.ResponseWriter).Status() >= 400:
			logger.Warn(req.URL.String())
		default:
			logger.Info(req.URL.String())
		}
	})
}

// NewGitlabSecretTokenMiddleware rejects all requests using the wrong Gitlab secret token.
func NewGitlabSecretTokenMiddleware(token string) negroni.Handler {
	return negroni.HandlerFunc(func(writer http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
		xtoken := req.Header.Get("X-Gitlab-Token")
		if xtoken != token {
			// NOTE: this returns 500 in order to notify Gitlab to disable this webhook. See
			//			 https://docs.gitlab.com/ee/user/project/integrations/webhooks.html#configure-your-webhook-receiver-endpoint
			//			 for more details.
			writer.WriteHeader(http.StatusInternalServerError)
			_, _ = writer.Write(errToJson("invalid Gitlab webhook secret token"))
			return
		}

		next(writer, req)
	})
}

// NewRecoverMiddleware recovers from any panic and notify Gitlab through a 400.
func NewRecoverMiddleware(logger *zap.Logger) negroni.Handler {
	logger = logger.WithOptions(zap.AddStacktrace(zap.ErrorLevel))

	return negroni.HandlerFunc(func(writer http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
		defer func() {
			if err := recover(); err != nil {
				writer.WriteHeader(http.StatusInternalServerError)
				_, _ = writer.Write(errToJson(err))

				logger.Error("panic recovered", zap.String("error", fmt.Sprintf("%s", err)))
			}
		}()

		next(writer, req)
	})
}
