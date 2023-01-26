package gitlab_events

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xanzy/go-gitlab"
)

func TestJobEvent_ProjectName(t *testing.T) {
	tcases := map[string]struct {
		event    JobEvent
		expected string
	}{
		"simple":    {JobEvent{JobEvent: gitlab.JobEvent{ProjectName: "aaa / bbb"}}, "aaa/bbb"},
		"formatted": {JobEvent{JobEvent: gitlab.JobEvent{ProjectName: "aaa/bbb"}}, "aaa/bbb"},
		"multiple":  {JobEvent{JobEvent: gitlab.JobEvent{ProjectName: "aaa / bbb / ccc / ddd"}}, "aaa/bbb/ccc/ddd"},
		"random":    {JobEvent{JobEvent: gitlab.JobEvent{ProjectName: "aaa / bbb/ccc / ddd"}}, "aaa/bbb/ccc/ddd"},
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
		"ref:branch":        {JobEvent{JobEvent: gitlab.JobEvent{Ref: "master"}}, "master"},
		"ref:tag":           {JobEvent{JobEvent: gitlab.JobEvent{Ref: "v1.0.0", Tag: true}}, "v1.0.0"},
		"ref:merge_request": {JobEvent{JobEvent: gitlab.JobEvent{Ref: "refs/merge-requests/1029/merge"}}, "1029"},
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
		expected Kind
	}{
		"ref:branch":        {JobEvent{JobEvent: gitlab.JobEvent{Ref: "master"}}, KindBranch},
		"ref:tag":           {JobEvent{JobEvent: gitlab.JobEvent{Ref: "v1.0.0", Tag: true}}, KindTag},
		"ref:merge_request": {JobEvent{JobEvent: gitlab.JobEvent{Ref: "refs/merge-requests/1029/merge"}}, KindMergeRequest},
	}

	for name, tcase := range tcases {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tcase.expected, tcase.event.RefKind())
		})
	}
}
