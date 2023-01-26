package webhook

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/radiofrance/gitlab-ci-pipelines-exporter/pkg/gitlab_events"
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
	webhook.mux.Handle("/hooks", gitlab.With(negroni.WrapFunc(webhook.handleHook)))

	return webhook
}

func (webhook Webhook) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	webhook.mux.ServeHTTP(w, r)
}

func (webhook Webhook) handleHook(writer http.ResponseWriter, request *http.Request) {
	// NOTE: reject all non-POST requests
	if request.Method != http.MethodPost {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(request.Body)
	if err != nil {
		// NOTE: this should not arrive, but it will be handled by
		//			 the Recover middleware.
		panic(err)
	}
	_ = request.Body.Close()

	var payload struct {
		ObjectKind string `json:"object_kind"`
	}
	err = json.Unmarshal(body, &payload)
	if err != nil {
		// NOTE: we return this status code as recommended by Gitlab
		writer.WriteHeader(http.StatusBadRequest)
		gcpehttp.WriteError(writer, err)
		return
	}

	switch payload.ObjectKind {
	case "pipeline":
		event := gitlab_events.PipelineEvent{}
		if err := unmarshallEvent(body, &event, writer); err != nil {
			return
		}
		webhook.handlePipelineEvent(event)
	case "build":
		event := gitlab_events.JobEvent{}
		if err := unmarshallEvent(body, &event, writer); err != nil {
			return
		}
		webhook.handleJobEvent(event)
	default:
		writer.WriteHeader(http.StatusNotImplemented)
		return
	}

	writer.WriteHeader(http.StatusOK)
}

// unmarshallEvent unmarshalls the JSON event data into the target object.
// The obj argument is a pointer to struct with json tags.
// If an error occurs it writes the error to the ResponseWriter and returns the error.
func unmarshallEvent(data []byte, obj any, writer http.ResponseWriter) error {
	err := json.Unmarshal(data, obj)
	if err != nil {
		// NOTE: we return this status code as recommended by Gitlab
		writer.WriteHeader(http.StatusBadRequest)
		gcpehttp.WriteError(writer, fmt.Errorf("failed to unmarshall json into %T", obj))
	}

	return err
}
