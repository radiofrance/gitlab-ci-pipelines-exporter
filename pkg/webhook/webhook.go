package webhook

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/xunleii/gitlab-ci-pipelines-exporter/pkg/collectors"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/urfave/negroni"
	"go.uber.org/zap"
)

type (
	// Collectors groups all Prometheus collectors used to exporter Gitlab CI metrics.
	Collectors struct {
		IDCollector                       *prometheus.GaugeVec
		DurationSecondsCollector          *prometheus.GaugeVec
		QueuedDurationSecondsCollector    *prometheus.GaugeVec
		RunCountCollector                 *prometheus.CounterVec
		StatusCollector                   *prometheus.GaugeVec
		TimestampCollector                *prometheus.GaugeVec
		JobIDCollector                    *prometheus.GaugeVec
		JobDurationSecondsCollector       *prometheus.GaugeVec
		JobQueuedDurationSecondsCollector *prometheus.GaugeVec
		JobRunCountCollector              *prometheus.CounterVec
		JobStatusCollector                *prometheus.GaugeVec
		JobTimestampCollector             *prometheus.GaugeVec
	}

	Webhook struct {
		http.Handler
		Collectors

		log *zap.Logger
	}
)

// NewWebhook creates a new instance of the Gitlab event webhook handler.
func NewWebhook(pattern, secretToken string, opts ...Option) *Webhook {
	webhook := &Webhook{
		Handler: http.NewServeMux(),
		Collectors: Collectors{
			IDCollector:                       collectors.NewCollectorID(),
			DurationSecondsCollector:          collectors.NewCollectorDurationSeconds(),
			QueuedDurationSecondsCollector:    collectors.NewCollectorQueuedDurationSeconds(),
			RunCountCollector:                 collectors.NewCollectorRunCount(),
			StatusCollector:                   collectors.NewCollectorStatus(),
			TimestampCollector:                collectors.NewCollectorTimestamp(),
			JobIDCollector:                    collectors.NewCollectorJobID(),
			JobDurationSecondsCollector:       collectors.NewCollectorJobDurationSeconds(),
			JobQueuedDurationSecondsCollector: collectors.NewCollectorJobQueuedDurationSeconds(),
			JobRunCountCollector:              collectors.NewCollectorJobRunCount(),
			JobStatusCollector:                collectors.NewCollectorJobStatus(),
			JobTimestampCollector:             collectors.NewCollectorJobTimestamp(),
		},

		log: zap.Must(zap.NewProduction()),
	}

	prometheus.MustRegister(
		webhook.IDCollector,
		webhook.DurationSecondsCollector,
		webhook.QueuedDurationSecondsCollector,
		webhook.RunCountCollector,
		webhook.StatusCollector,
		webhook.TimestampCollector,
		webhook.JobIDCollector,
		webhook.JobDurationSecondsCollector,
		webhook.JobQueuedDurationSecondsCollector,
		webhook.JobRunCountCollector,
		webhook.JobStatusCollector,
		webhook.JobTimestampCollector,
	)

	for _, opt := range opts {
		opt(webhook)
	}

	mux := webhook.Handler.(*http.ServeMux)
	mux.Handle("/healthz", http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
		writer.WriteHeader(200)
	}))

	root := negroni.New(NewRecoverMiddleware(webhook.log), NewZapMiddleware(webhook.log))
	mux.Handle(pattern, root.With(negroni.Wrap(promhttp.Handler())))

	gitlab := root.With(NewGitlabSecretTokenMiddleware(secretToken))
	mux.Handle("/pipeline", gitlab.With(negroni.Wrap(processHandler(webhook.handlePipelineEvent))))
	mux.Handle("/job", gitlab.With(negroni.Wrap(processHandler(webhook.handleJobEvent))))

	return webhook
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
			_, _ = writer.Write(errToJson(err))
			return
		}

		err = handler(payload)
		if err != nil {
			// NOTE: we return this status code as recommended by Gitlab
			writer.WriteHeader(http.StatusBadRequest)
			_, _ = writer.Write(errToJson(err))
			return
		}

		writer.WriteHeader(http.StatusOK)
	}
}
