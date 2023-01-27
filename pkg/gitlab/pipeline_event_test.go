package gitlab_test

import (
	"testing"

	"github.com/radiofrance/gitlab-ci-pipelines-exporter/pkg/gitlab"
	"github.com/stretchr/testify/assert"
)

func TestPipelineEvent_Ref(t *testing.T) {
	t.Parallel()

	t.Run("ref:branch", func(t *testing.T) {
		t.Parallel()

		event := gitlab.PipelineEvent{}
		event.ObjectAttributes.Ref = "master"

		assert.Equal(t, "master", event.Ref())
	})
	t.Run("ref:tag", func(t *testing.T) {
		t.Parallel()

		event := gitlab.PipelineEvent{}
		event.ObjectAttributes.Ref = "v1.0.0"
		event.ObjectAttributes.Tag = true

		assert.Equal(t, "v1.0.0", event.Ref())
	})
	t.Run("ref:merge_request", func(t *testing.T) {
		t.Parallel()

		event := gitlab.PipelineEvent{}
		event.MergeRequest.IID = 1029

		assert.Equal(t, "1029", event.Ref())
	})
}

func TestPipelineEvent_RefKind(t *testing.T) {
	t.Parallel()

	t.Run("ref:branch", func(t *testing.T) {
		t.Parallel()

		event := gitlab.PipelineEvent{}
		event.ObjectAttributes.Ref = "master"

		assert.Equal(t, gitlab.KindBranch, event.RefKind())
	})
	t.Run("ref:tag", func(t *testing.T) {
		t.Parallel()

		event := gitlab.PipelineEvent{}
		event.ObjectAttributes.Ref = "v1.0.0"
		event.ObjectAttributes.Tag = true

		assert.Equal(t, gitlab.KindTag, event.RefKind())
	})
	t.Run("ref:merge_request", func(t *testing.T) {
		t.Parallel()

		event := gitlab.PipelineEvent{}
		event.MergeRequest.IID = 1029

		assert.Equal(t, gitlab.KindMergeRequest, event.RefKind())
	})
}
