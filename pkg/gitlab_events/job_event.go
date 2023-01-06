package gitlab_events

import (
	"strings"
)

type (
	// JobEvent contains all status information about a Gitlab Job event.
	JobEvent struct {
		FullRef             string          `json:"ref,omitempty"`
		Tag                 bool            `json:"tag,omitempty"`
		BuildID             int             `json:"build_id,omitempty"`
		BuildName           string          `json:"build_name,omitempty"`
		BuildStage          string          `json:"build_stage,omitempty"`
		BuildStatus         status          `json:"build_status,omitempty"`
		BuildDuration       Option[float64] `json:"build_duration,omitempty"`
		BuildQueuedDuration Option[float64] `json:"build_queued_duration,omitempty"`
		PipelineID          int             `json:"pipeline_id,omitempty"`

		Runner struct {
			Description string `json:"description"`
		} `json:"runner,omitempty"`

		SpacedProjectName string `json:"project_name,omitempty"`
	}
)

// ProjectName returns the Gitlab project name.
func (j JobEvent) ProjectName() string {
	return strings.ReplaceAll(j.SpacedProjectName, " / ", "/")
}

// RefKind returns what kind of ref as generated the event.
func (j JobEvent) RefKind() refKind {
	switch {
	case j.Tag:
		return TagKind
	case strings.HasPrefix(j.FullRef, "refs/merge-requests/"):
		return MergeRequestKind
	default:
		return BranchKind
	}
}

// Ref returns the reference that triggers this event.
func (j JobEvent) Ref() string {
	switch j.RefKind() {
	case MergeRequestKind:
		ref := j.FullRef
		ref = strings.TrimPrefix(ref, "refs/merge-requests/")
		ref = strings.TrimSuffix(ref, "/merge")

		return ref
	case TagKind, BranchKind:
		fallthrough
	default:
		return j.FullRef
	}
}

// Stage returns the stage name of the current job.
func (j JobEvent) Stage() string {
	return j.BuildStage
}

// JobName returns the current job name.
func (j JobEvent) JobName() string {
	return j.BuildName
}

// RunnerDescription returns the runner description that ran this job.
func (j JobEvent) RunnerDescription() string {
	return j.Runner.Description
}
