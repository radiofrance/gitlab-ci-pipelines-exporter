package webhook

import (
	"time"

	"github.com/radiofrance/gitlab-ci-pipelines-exporter/pkg/gitlab_events"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	// btof allows to convert quickly bool into float value
	btof = map[bool]float64{false: 0., true: 1.}

	// timestamp returns the timestamp that should be used by the timestamp* collector.
	// Putting it in a variable allows us to simulate it during the tests.
	timestamp = func() float64 { return float64(time.Now().Unix()) }
)

// handlePipelineEvent handles Gitlab pipeline events.
func (w Webhook) handlePipelineEvent(event gitlab_events.PipelineEvent) {
	kind := event.RefKind()
	labels := prometheus.Labels{
		"project": event.Project.PathWithNamespace,
		"ref":     event.Ref(),
		"kind":    kind.String(),
	}

	w.collectors.ID.With(labels).Set(float64(event.ObjectAttributes.ID))
	w.collectors.Timestamp.With(labels).Set(timestamp())

	if event.ObjectAttributes.QueuedDuration != 0 {
		w.collectors.QueuedDurationSeconds.With(labels).Set(float64(event.ObjectAttributes.QueuedDuration))
	}
	if event.ObjectAttributes.Duration != 0 {
		w.collectors.DurationSeconds.With(labels).Set(float64(event.ObjectAttributes.Duration))
	}

	if event.ObjectAttributes.Status == gitlab_events.StatusRunning.String() {
		w.collectors.RunCount.With(labels).Inc()
	}

	for _, status := range gitlab_events.Statuses[1:] {
		labels["status"] = status
		w.collectors.Status.With(labels).Set(btof[event.ObjectAttributes.Status == status])
	}
}

// handleJobEvent handles Gitlab job events.
func (w Webhook) handleJobEvent(event gitlab_events.JobEvent) {
	labels := prometheus.Labels{
		"project":  event.ProjectName(),
		"ref":      event.Ref(),
		"kind":     event.RefKind().String(),
		"stage":    event.BuildStage,
		"job_name": event.BuildName,
	}

	w.collectors.JobID.With(labels).Set(float64(event.BuildID))
	w.collectors.JobTimestamp.With(labels).Set(timestamp())

	if event.BuildQueuedDuration != 0 {
		w.collectors.JobQueuedDurationSeconds.With(labels).Set(event.BuildQueuedDuration)
	}
	if event.BuildDuration != 0 {
		w.collectors.JobDurationSeconds.With(labels).Set(event.BuildDuration)
	}

	if event.BuildStatus == gitlab_events.StatusRunning.String() {
		w.collectors.JobRunCount.With(labels).Inc()
	}

	for _, status := range gitlab_events.Statuses[1:] {
		labels["status"] = status
		w.collectors.JobStatus.With(labels).Set(btof[event.BuildStatus == status])
	}
}
