package gitlab_events

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRefKind_String(t *testing.T) {
	assert.Equal(t, "branch", KindBranch.String())
	assert.Equal(t, "tag", KindTag.String())
	assert.Equal(t, "merge_request", KindMergeRequest.String())
}

func TestStatusFromString(t *testing.T) {
	tcases := map[string]Status{
		"unknown":              StatusUnknown,
		"created":              StatusCreated,
		"waiting_for_resource": StatusWaitingForResource,
		"preparing":            StatusPreparing,
		"pending":              StatusPending,
		"running":              StatusRunning,
		"success":              StatusSuccess,
		"failed":               StatusFailed,
		"canceled":             StatusCanceled,
		"skipped":              StatusSkipped,
		"manual":               StatusManual,
		"scheduled":            StatusScheduled,
	}

	for str, status := range tcases {
		assert.Equal(t, status, StatusFromString(str))
	}
}

func TestStatus_String(t *testing.T) {
	assert.Equal(t, "unknown", StatusUnknown.String())
	assert.Equal(t, "created", StatusCreated.String())
	assert.Equal(t, "waiting_for_resource", StatusWaitingForResource.String())
	assert.Equal(t, "preparing", StatusPreparing.String())
	assert.Equal(t, "pending", StatusPending.String())
	assert.Equal(t, "running", StatusRunning.String())
	assert.Equal(t, "success", StatusSuccess.String())
	assert.Equal(t, "failed", StatusFailed.String())
	assert.Equal(t, "canceled", StatusCanceled.String())
	assert.Equal(t, "skipped", StatusSkipped.String())
	assert.Equal(t, "manual", StatusManual.String())
	assert.Equal(t, "scheduled", StatusScheduled.String())
}
