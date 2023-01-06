package gitlab_events

import "strconv"

type (
	// PipelineEvent contains all status information about a Gitlab Pipeline event.
	PipelineEvent struct {
		ObjectAttributes struct {
			ID             int             `json:"id,omitempty"`
			Ref            string          `json:"ref,omitempty"`
			Tag            bool            `json:"tag,omitempty"`
			Status         status          `json:"status,omitempty"`
			Duration       Option[float64] `json:"duration,omitempty"`
			QueuedDuration Option[float64] `json:"queued_duration,omitempty"`
		} `json:"object_attributes,omitempty"`

		MergeRequest Option[struct {
			IID int `json:"iid"`
		}] `json:"merge_request,omitempty"`

		Project struct {
			PathWithNamespace string `json:"path_with_namespace,omitempty"`
		} `json:"project,omitempty"`
	}
)

// ProjectName returns the Gitlab project name.
func (p PipelineEvent) ProjectName() string {
	return p.Project.PathWithNamespace
}

// RefKind returns what kind of ref as generated the event.
func (p PipelineEvent) RefKind() refKind {
	switch {
	case p.ObjectAttributes.Tag:
		return TagKind
	case p.MergeRequest.IsSome():
		return MergeRequestKind
	default:
		return BranchKind
	}
}

// Ref returns the reference that triggers this event.
func (p PipelineEvent) Ref() string {
	switch p.RefKind() {
	case MergeRequestKind:
		return strconv.Itoa(p.MergeRequest.Unwrap().IID)
	case TagKind, BranchKind:
		fallthrough
	default:
		return p.ObjectAttributes.Ref
	}
}
