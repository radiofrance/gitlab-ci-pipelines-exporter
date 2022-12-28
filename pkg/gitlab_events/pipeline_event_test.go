package gitlab_events

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPipelineEvent_ProjectName(t *testing.T) {
	event := PipelineEvent{}
	event.Project.PathWithNamespace = "aaa/bbb"

	assert.Equal(t, "aaa/bbb", event.ProjectName())
}

func TestPipelineEvent_Ref(t *testing.T) {
	t.Run("ref:branch", func(t *testing.T) {
		event := PipelineEvent{}
		event.ObjectAttributes.Ref = "master"

		assert.Equal(t, "master", event.Ref())
	})
	t.Run("ref:tag", func(t *testing.T) {
		event := PipelineEvent{}
		event.ObjectAttributes.Ref = "v1.0.0"
		event.ObjectAttributes.Tag = true

		assert.Equal(t, "v1.0.0", event.Ref())
	})
	t.Run("ref:merge_request", func(t *testing.T) {
		event := PipelineEvent{}
		event.MergeRequest = Some(struct {
			IID int `json:"iid"`
		}{IID: 1029})

		assert.Equal(t, "1029", event.Ref())
	})
}

func TestPipelineEvent_RefKind(t *testing.T) {
	t.Run("ref:branch", func(t *testing.T) {
		event := PipelineEvent{}
		event.ObjectAttributes.Ref = "master"

		assert.Equal(t, BranchKind, event.RefKind())
	})
	t.Run("ref:tag", func(t *testing.T) {
		event := PipelineEvent{}
		event.ObjectAttributes.Ref = "v1.0.0"
		event.ObjectAttributes.Tag = true

		assert.Equal(t, TagKind, event.RefKind())
	})
	t.Run("ref:merge_request", func(t *testing.T) {
		event := PipelineEvent{}
		event.MergeRequest = Some(struct {
			IID int `json:"iid"`
		}{IID: 1029})

		assert.Equal(t, MergeRequestKind, event.RefKind())
	})
}
