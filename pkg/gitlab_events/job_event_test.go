package gitlab_events

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJobEvent_ProjectName(t *testing.T) {
	tcases := map[string]struct {
		event    JobEvent
		expected string
	}{
		"simple":    {JobEvent{SpacedProjectName: "aaa / bbb"}, "aaa/bbb"},
		"formatted": {JobEvent{SpacedProjectName: "aaa/bbb"}, "aaa/bbb"},
		"multiple":  {JobEvent{SpacedProjectName: "aaa / bbb / ccc / ddd"}, "aaa/bbb/ccc/ddd"},
		"random":    {JobEvent{SpacedProjectName: "aaa / bbb/ccc / ddd"}, "aaa/bbb/ccc/ddd"},
	}

	for name, tcase := range tcases {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tcase.expected, tcase.event.ProjectName())
		})
	}
}

func TestJobEvent_Ref(t *testing.T) {
	tcases := map[string]struct {
		event    JobEvent
		expected string
	}{
		"ref:branch":        {JobEvent{FullRef: "master"}, "master"},
		"ref:tag":           {JobEvent{FullRef: "v1.0.0", Tag: true}, "v1.0.0"},
		"ref:merge_request": {JobEvent{FullRef: "refs/merge-requests/1029/merge"}, "1029"},
	}

	for name, tcase := range tcases {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tcase.expected, tcase.event.Ref())
		})
	}
}

func TestJobEvent_RefKind(t *testing.T) {
	tcases := map[string]struct {
		event    JobEvent
		expected refKind
	}{
		"ref:branch":        {JobEvent{FullRef: "master"}, BranchKind},
		"ref:tag":           {JobEvent{FullRef: "v1.0.0", Tag: true}, TagKind},
		"ref:merge_request": {JobEvent{FullRef: "refs/merge-requests/1029/merge"}, MergeRequestKind},
	}

	for name, tcase := range tcases {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tcase.expected, tcase.event.RefKind())
		})
	}
}

func TestJobEvent_Stage(t *testing.T) {
	assert.Equal(t, "quality", JobEvent{BuildStage: "quality"}.Stage())
}

func TestJobEvent_JobName(t *testing.T) {
	assert.Equal(t, "sonarqube", JobEvent{BuildName: "sonarqube"}.JobName())
}

func TestJobEvent_RunnerDescription(t *testing.T) {
	assert.Equal(
		t,
		"gitlab-runner-shared-69fd46fcd-w7x5m",
		JobEvent{Runner: struct {
			Description string `json:"description"`
		}{Description: "gitlab-runner-shared-69fd46fcd-w7x5m"},
		}.RunnerDescription(),
	)
}
