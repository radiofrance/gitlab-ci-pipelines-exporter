package gitlab

import (
	"strings"

	"github.com/xanzy/go-gitlab"
)

type (
	// JobEvent contains all status information about a Gitlab Job event.
	JobEvent struct {
		gitlab.JobEvent

		BuildQueuedDuration float64 `json:"build_queued_duration,omitempty"`
	}
)

// ProjectName returns the Gitlab project name.
func (j JobEvent) ProjectName() string {
	return strings.ReplaceAll(j.JobEvent.ProjectName, " / ", "/")
}

// RefKind returns what kind of ref as generated the event.
func (j JobEvent) RefKind() Kind {
	switch {
	case j.Tag:
		return KindTag
	case strings.HasPrefix(j.JobEvent.Ref, "refs/merge-requests/"):
		return KindMergeRequest
	default:
		return KindBranch
	}
}

// Ref returns the reference that triggers this event.
func (j JobEvent) Ref() string {
	switch j.RefKind() {
	case KindMergeRequest:
		ref := j.JobEvent.Ref
		ref = strings.TrimPrefix(ref, "refs/merge-requests/")
		ref = strings.TrimSuffix(ref, "/merge")

		return ref
	case KindTag, KindBranch:
		fallthrough
	default:
		return j.JobEvent.Ref
	}
}
