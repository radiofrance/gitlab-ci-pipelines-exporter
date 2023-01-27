package gitlab_test

import (
	"testing"

	"github.com/radiofrance/gitlab-ci-pipelines-exporter/pkg/gitlab"
	"github.com/stretchr/testify/assert"
)

func TestRefKind_String(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "branch", gitlab.KindBranch.String())
	assert.Equal(t, "tag", gitlab.KindTag.String())
	assert.Equal(t, "merge_request", gitlab.KindMergeRequest.String())
}

func TestStatusFromString(t *testing.T) {
	t.Parallel()

	tcases := map[string]gitlab.Status{
		"unknown":              gitlab.StatusUnknown,
		"created":              gitlab.StatusCreated,
		"waiting_for_resource": gitlab.StatusWaitingForResource,
		"preparing":            gitlab.StatusPreparing,
		"pending":              gitlab.StatusPending,
		"running":              gitlab.StatusRunning,
		"success":              gitlab.StatusSuccess,
		"failed":               gitlab.StatusFailed,
		"canceled":             gitlab.StatusCanceled,
		"skipped":              gitlab.StatusSkipped,
		"manual":               gitlab.StatusManual,
		"scheduled":            gitlab.StatusScheduled,
	}

	for str, status := range tcases {
		assert.Equal(t, status, gitlab.StatusFromString(str))
	}
}

func TestStatus_String(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "unknown", gitlab.StatusUnknown.String())
	assert.Equal(t, "created", gitlab.StatusCreated.String())
	assert.Equal(t, "waiting_for_resource", gitlab.StatusWaitingForResource.String())
	assert.Equal(t, "preparing", gitlab.StatusPreparing.String())
	assert.Equal(t, "pending", gitlab.StatusPending.String())
	assert.Equal(t, "running", gitlab.StatusRunning.String())
	assert.Equal(t, "success", gitlab.StatusSuccess.String())
	assert.Equal(t, "failed", gitlab.StatusFailed.String())
	assert.Equal(t, "canceled", gitlab.StatusCanceled.String())
	assert.Equal(t, "skipped", gitlab.StatusSkipped.String())
	assert.Equal(t, "manual", gitlab.StatusManual.String())
	assert.Equal(t, "scheduled", gitlab.StatusScheduled.String())
}
