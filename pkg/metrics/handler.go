package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/radiofrance/gitlab-ci-pipelines-exporter/pkg/collectors"
	gcpehttp "github.com/radiofrance/gitlab-ci-pipelines-exporter/pkg/http"
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

	Handler struct {
		Collectors

		mux *http.ServeMux
		log *zap.Logger
	}
)

// NewHandler creates a new instance of a Handler object.
func NewHandler(pattern string, collectors Collectors, opts ...Option) *Handler {
	handler := &Handler{
		Collectors: collectors,

		mux: http.NewServeMux(),
		log: zap.Must(zap.NewProduction()),
	}

	prometheus.MustRegister(
		handler.IDCollector,
		handler.DurationSecondsCollector,
		handler.QueuedDurationSecondsCollector,
		handler.RunCountCollector,
		handler.StatusCollector,
		handler.TimestampCollector,
		handler.JobIDCollector,
		handler.JobDurationSecondsCollector,
		handler.JobQueuedDurationSecondsCollector,
		handler.JobRunCountCollector,
		handler.JobStatusCollector,
		handler.JobTimestampCollector,
	)

	for _, opt := range opts {
		opt(handler)
	}

	root := negroni.New(gcpehttp.NewRecoverMiddleware(handler.log), gcpehttp.NewZapMiddleware(handler.log))
	handler.mux.Handle(pattern, root.With(negroni.Wrap(promhttp.Handler())))

	return handler
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}

func AllCollectors() Collectors {
	return Collectors{
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
	}
}
