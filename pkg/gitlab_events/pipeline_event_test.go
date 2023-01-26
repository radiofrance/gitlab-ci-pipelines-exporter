package gitlab_events

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPipelineEvent_Ref(t *testing.T) {
	t.Run("ref:branch", func(t *testing.T) {
		event := PipelineEvent{}
		event.Project.PathWithNamespace = "aaa/bbb"
		event.ObjectAttributes.Ref = "master"

		assert.Equal(t, "master", event.Ref())
	})
	t.Run("ref:tag", func(t *testing.T) {
		event := PipelineEvent{}
		event.Project.PathWithNamespace = "aaa/bbb"
		event.ObjectAttributes.Ref = "v1.0.0"
		event.ObjectAttributes.Tag = true

		assert.Equal(t, "v1.0.0", event.Ref())
	})
	t.Run("ref:merge_request", func(t *testing.T) {
		event := PipelineEvent{}
		event.Project.PathWithNamespace = "aaa/bbb"
		event.MergeRequest.IID = 1029

		assert.Equal(t, "1029", event.Ref())
	})
}

func TestPipelineEvent_RefKind(t *testing.T) {
	t.Run("ref:branch", func(t *testing.T) {
		event := PipelineEvent{}
		event.ObjectAttributes.Ref = "master"

		assert.Equal(t, KindBranch, event.RefKind())
	})
	t.Run("ref:tag", func(t *testing.T) {
		event := PipelineEvent{}
		event.ObjectAttributes.Ref = "v1.0.0"
		event.ObjectAttributes.Tag = true

		assert.Equal(t, KindTag, event.RefKind())
	})
	t.Run("ref:merge_request", func(t *testing.T) {
		event := PipelineEvent{}
		event.MergeRequest.IID = 1029

		assert.Equal(t, KindMergeRequest, event.RefKind())
	})
}
