package webhook

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	gcpehttp "github.com/radiofrance/gitlab-ci-pipelines-exporter/pkg/http"
	"github.com/radiofrance/gitlab-ci-pipelines-exporter/pkg/metrics"
	"github.com/urfave/negroni"
	"go.uber.org/zap"
)

type (
	Webhook struct {
		collectors *metrics.PrometheusCollectors
		mux        *http.ServeMux
		log        *zap.Logger
	}
)

// NewWebhook creates a new instance of the Gitlab event webhook handler.
func NewWebhook(secretToken string, collectors *metrics.PrometheusCollectors, opts ...Option) *Webhook {
	webhook := &Webhook{
		collectors: collectors,
		mux:        http.NewServeMux(),
		log:        zap.Must(zap.NewProduction()),
	}

	for _, opt := range opts {
		opt(webhook)
	}

	webhook.mux.Handle("/healthz", http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
		writer.WriteHeader(200)
	}))

	root := negroni.New(gcpehttp.NewRecoverMiddleware(webhook.log), gcpehttp.NewZapMiddleware(webhook.log))
	gitlab := root.With(NewGitlabSecretTokenMiddleware(secretToken))
	webhook.mux.Handle("/pipeline", gitlab.With(negroni.Wrap(processHandler(webhook.handlePipelineEvent))))
	webhook.mux.Handle("/job", gitlab.With(negroni.Wrap(processHandler(webhook.handleJobEvent))))

	return webhook
}

func (webhook Webhook) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	webhook.mux.ServeHTTP(w, r)
}

// processHandler wraps the event handlers to simplify their usage.
func processHandler[T any](handler func(event T) error) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		// NOTE: reject all non-POST requests
		if request.Method != http.MethodPost {
			writer.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		body := &bytes.Buffer{}
		_, err := io.Copy(body, request.Body)
		if err != nil {
			// NOTE: this should not arrive, but it will be handled by
			//			 the Recover middleware.
			panic(err)
		}
		_ = request.Body.Close()

		var payload T
		err = json.Unmarshal(body.Bytes(), &payload)
		if err != nil {
			// NOTE: we return this status code as recommended by Gitlab
			writer.WriteHeader(http.StatusBadRequest)
			gcpehttp.WriteError(writer, err)
			return
		}

		err = handler(payload)
		if err != nil {
			// NOTE: we return this status code as recommended by Gitlab
			writer.WriteHeader(http.StatusBadRequest)
			gcpehttp.WriteError(writer, err)
			return
		}

		writer.WriteHeader(http.StatusOK)
	}
}
