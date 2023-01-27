package gitlab

import (
	"strconv"

	"github.com/xanzy/go-gitlab"
)

type (
	// PipelineEvent contains all status information about a Gitlab Pipeline event.
	PipelineEvent struct {
		gitlab.PipelineEvent
	}
)

// RefKind returns what kind of ref as generated the event.
func (p PipelineEvent) RefKind() Kind {
	switch {
	case p.ObjectAttributes.Tag:
		return KindTag
	case p.MergeRequest.IID != 0:
		return KindMergeRequest
	default:
		return KindBranch
	}
}

// Ref returns the reference that triggers this event.
func (p PipelineEvent) Ref() string {
	switch p.RefKind() {
	case KindMergeRequest:
		return strconv.Itoa(p.MergeRequest.IID)
	case KindTag, KindBranch:
		fallthrough
	default:
		return p.ObjectAttributes.Ref
	}
}
