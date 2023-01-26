package webhook

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/radiofrance/gitlab-ci-pipelines-exporter/pkg/gitlab"
)

// btof allows to convert quickly bool into float value.
var btof = map[bool]float64{false: 0., true: 1.}

// handlePipelineEvent handles Gitlab pipeline events.
func (wh Webhook) handlePipelineEvent(event gitlab.PipelineEvent) {
	kind := event.RefKind()
	labels := prometheus.Labels{
		"project": event.Project.PathWithNamespace,
		"ref":     event.Ref(),
		"kind":    kind.String(),
	}

	wh.collectors.ID.With(labels).Set(float64(event.ObjectAttributes.ID))
	wh.collectors.Timestamp.With(labels).Set(float64(wh.timestamp()))

	if event.ObjectAttributes.QueuedDuration != 0 {
		wh.collectors.QueuedDurationSeconds.With(labels).Set(float64(event.ObjectAttributes.QueuedDuration))
	}
	if event.ObjectAttributes.Duration != 0 {
		wh.collectors.DurationSeconds.With(labels).Set(float64(event.ObjectAttributes.Duration))
	}

	if event.ObjectAttributes.Status == gitlab.StatusRunning.String() {
		wh.collectors.RunCount.With(labels).Inc()
	}

	for _, status := range gitlab.Statuses[1:] {
		labels["status"] = status
		wh.collectors.Status.With(labels).Set(btof[event.ObjectAttributes.Status == status])
	}
}

// handleJobEvent handles Gitlab job events.
func (wh Webhook) handleJobEvent(event gitlab.JobEvent) {
	labels := prometheus.Labels{
		"project":  event.ProjectName(),
		"ref":      event.Ref(),
		"kind":     event.RefKind().String(),
		"stage":    event.BuildStage,
		"job_name": event.BuildName,
	}

	wh.collectors.JobID.With(labels).Set(float64(event.BuildID))
	wh.collectors.JobTimestamp.With(labels).Set(float64(wh.timestamp()))

	if event.BuildQueuedDuration != 0 {
		wh.collectors.JobQueuedDurationSeconds.With(labels).Set(event.BuildQueuedDuration)
	}
	if event.BuildDuration != 0 {
		wh.collectors.JobDurationSeconds.With(labels).Set(event.BuildDuration)
	}

	if event.BuildStatus == gitlab.StatusRunning.String() {
		wh.collectors.JobRunCount.With(labels).Inc()
	}

	for _, status := range gitlab.Statuses[1:] {
		labels["status"] = status
		wh.collectors.JobStatus.With(labels).Set(btof[event.BuildStatus == status])
	}
}
