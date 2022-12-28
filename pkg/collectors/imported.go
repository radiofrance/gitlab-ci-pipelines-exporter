// Package collectors exists to make this exporter compatible with mvisonneau/gitlab-ci-pipelines-exporter.
// See https://github.com/mvisonneau/gitlab-ci-pipelines-exporter/blob/v0.5.4/pkg/controller/collectors.go for more
// details.
package collectors

import "github.com/prometheus/client_golang/prometheus"

var (
	// topics and variables are removed (too much cardinality)
	defaultLabels = []string{"project", "kind", "ref"}
	// runner_description is removed (too much cardinality on job metrics)
	jobLabels    = []string{"stage", "job_name"}
	statusLabels = []string{"status"}
)

// NewCollectorDurationSeconds returns a new collector for the gitlab_ci_pipeline_duration_seconds metric.
func NewCollectorDurationSeconds() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gitlab_ci_pipeline_duration_seconds",
			Help: "Duration in seconds of the most recent pipeline",
		},
		defaultLabels,
	)
}

// NewCollectorQueuedDurationSeconds returns a new collector for the gitlab_ci_pipeline_queued_duration_seconds metric.
func NewCollectorQueuedDurationSeconds() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gitlab_ci_pipeline_queued_duration_seconds",
			Help: "Duration in seconds the most recent pipeline has been queued before starting",
		},
		defaultLabels,
	)
}

// NewCollectorID returns a new collector for the gitlab_ci_pipeline_id metric.
func NewCollectorID() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gitlab_ci_pipeline_id",
			Help: "ID of the most recent pipeline",
		},
		defaultLabels,
	)
}

// NewCollectorJobDurationSeconds returns a new collector for the gitlab_ci_pipeline_job_duration_seconds metric.
func NewCollectorJobDurationSeconds() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gitlab_ci_pipeline_job_duration_seconds",
			Help: "Duration in seconds of the most recent job",
		},
		append(defaultLabels, jobLabels...),
	)
}

// NewCollectorJobID returns a new collector for the gitlab_ci_pipeline_job_id metric.
func NewCollectorJobID() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gitlab_ci_pipeline_job_id",
			Help: "ID of the most recent job",
		},
		append(defaultLabels, jobLabels...),
	)
}

// NewCollectorJobQueuedDurationSeconds returns a new collector for the gitlab_ci_pipeline_job_queued_duration_seconds metric.
func NewCollectorJobQueuedDurationSeconds() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gitlab_ci_pipeline_job_queued_duration_seconds",
			Help: "Duration in seconds the most recent job has been queued before starting",
		},
		append(defaultLabels, jobLabels...),
	)
}

// NewCollectorJobRunCount returns a new collector for the gitlab_ci_pipeline_job_run_count metric.
func NewCollectorJobRunCount() *prometheus.CounterVec {
	return prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "gitlab_ci_pipeline_job_run_count",
			Help: "Number of executions of a job",
		},
		append(defaultLabels, jobLabels...),
	)
}

// NewCollectorJobStatus returns a new collector for the gitlab_ci_pipeline_job_status metric.
func NewCollectorJobStatus() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gitlab_ci_pipeline_job_status",
			Help: "Status of the most recent job",
		},
		append(defaultLabels, append(jobLabels, statusLabels...)...),
	)
}

// NewCollectorJobTimestamp returns a new collector for the gitlab_ci_pipeline_job_timestamp metric.
func NewCollectorJobTimestamp() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gitlab_ci_pipeline_job_timestamp",
			Help: "Creation date timestamp of the the most recent job",
		},
		append(defaultLabels, jobLabels...),
	)
}

// NewCollectorStatus returns a new collector for the gitlab_ci_pipeline_status metric.
func NewCollectorStatus() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gitlab_ci_pipeline_status",
			Help: "Status of the most recent pipeline",
		},
		append(defaultLabels, "status"),
	)
}

// NewCollectorTimestamp returns a new collector for the gitlab_ci_pipeline_timestamp metric.
func NewCollectorTimestamp() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gitlab_ci_pipeline_timestamp",
			Help: "Timestamp of the last update of the most recent pipeline",
		},
		defaultLabels,
	)
}

// NewCollectorRunCount returns a new collector for the gitlab_ci_pipeline_run_count metric.
func NewCollectorRunCount() *prometheus.CounterVec {
	return prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "gitlab_ci_pipeline_run_count",
			Help: "Number of executions of a pipeline",
		},
		defaultLabels,
	)
}
