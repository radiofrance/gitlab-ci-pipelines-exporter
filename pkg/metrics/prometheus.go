package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/radiofrance/gitlab-ci-pipelines-exporter/pkg/collectors"
)

// PrometheusCollectors groups all Prometheus collectors used to exporter Gitlab CI metrics.
type PrometheusCollectors struct {
	ID                       *prometheus.GaugeVec
	DurationSeconds          *prometheus.GaugeVec
	QueuedDurationSeconds    *prometheus.GaugeVec
	RunCount                 *prometheus.CounterVec
	Status                   *prometheus.GaugeVec
	Timestamp                *prometheus.GaugeVec
	JobID                    *prometheus.GaugeVec
	JobDurationSeconds       *prometheus.GaugeVec
	JobQueuedDurationSeconds *prometheus.GaugeVec
	JobRunCount              *prometheus.CounterVec
	JobStatus                *prometheus.GaugeVec
	JobTimestamp             *prometheus.GaugeVec
}

// NewPrometheusCollectors creates a new PrometheusCollectors instance.
func NewPrometheusCollectors() *PrometheusCollectors {
	return &PrometheusCollectors{
		ID:                       collectors.NewCollectorID(),
		DurationSeconds:          collectors.NewCollectorDurationSeconds(),
		QueuedDurationSeconds:    collectors.NewCollectorQueuedDurationSeconds(),
		RunCount:                 collectors.NewCollectorRunCount(),
		Status:                   collectors.NewCollectorStatus(),
		Timestamp:                collectors.NewCollectorTimestamp(),
		JobID:                    collectors.NewCollectorJobID(),
		JobDurationSeconds:       collectors.NewCollectorJobDurationSeconds(),
		JobQueuedDurationSeconds: collectors.NewCollectorJobQueuedDurationSeconds(),
		JobRunCount:              collectors.NewCollectorJobRunCount(),
		JobStatus:                collectors.NewCollectorJobStatus(),
		JobTimestamp:             collectors.NewCollectorJobTimestamp(),
	}
}

// MustRegister registers the Prometheus collectors and panics if any error occurs.
func (c PrometheusCollectors) MustRegister() {
	prometheus.MustRegister(
		c.ID,
		c.DurationSeconds,
		c.QueuedDurationSeconds,
		c.RunCount,
		c.Status,
		c.Timestamp,
		c.JobID,
		c.JobDurationSeconds,
		c.JobQueuedDurationSeconds,
		c.JobRunCount,
		c.JobStatus,
		c.JobTimestamp,
	)
}
