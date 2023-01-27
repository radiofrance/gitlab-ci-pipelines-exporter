package gitlab_test

import (
	"testing"

	"github.com/radiofrance/gitlab-ci-pipelines-exporter/pkg/gitlab"
	"github.com/stretchr/testify/assert"
	gogitlab "github.com/xanzy/go-gitlab"
)

func TestJobEvent_ProjectName(t *testing.T) {
	t.Parallel()

	tcases := map[string]struct {
		event    gitlab.JobEvent
		expected string
	}{
		"simple": {
			gitlab.JobEvent{JobEvent: gogitlab.JobEvent{ProjectName: "aaa / bbb"}},
			"aaa/bbb",
		},
		"formatted": {
			gitlab.JobEvent{JobEvent: gogitlab.JobEvent{ProjectName: "aaa/bbb"}},
			"aaa/bbb",
		},
		"multiple": {
			gitlab.JobEvent{JobEvent: gogitlab.JobEvent{ProjectName: "aaa / bbb / ccc / ddd"}},
			"aaa/bbb/ccc/ddd",
		},
		"random": {
			gitlab.JobEvent{JobEvent: gogitlab.JobEvent{ProjectName: "aaa / bbb/ccc / ddd"}},
			"aaa/bbb/ccc/ddd",
		},
	}

	for name, tcase := range tcases {
		tcase := tcase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tcase.expected, tcase.event.ProjectName())
		})
	}
}

func TestJobEvent_Ref(t *testing.T) {
	t.Parallel()

	tcases := map[string]struct {
		event    gitlab.JobEvent
		expected string
	}{
		"ref:branch": {
			gitlab.JobEvent{JobEvent: gogitlab.JobEvent{Ref: "master"}},
			"master",
		},
		"ref:tag": {
			gitlab.JobEvent{JobEvent: gogitlab.JobEvent{Ref: "v1.0.0", Tag: true}},
			"v1.0.0",
		},
		"ref:merge_request": {
			gitlab.JobEvent{JobEvent: gogitlab.JobEvent{Ref: "refs/merge-requests/1029/merge"}},
			"1029",
		},
	}

	for name, tcase := range tcases {
		tcase := tcase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tcase.expected, tcase.event.Ref())
		})
	}
}

func TestJobEvent_RefKind(t *testing.T) {
	t.Parallel()

	tcases := map[string]struct {
		event    gitlab.JobEvent
		expected gitlab.Kind
	}{
		"ref:branch": {
			gitlab.JobEvent{JobEvent: gogitlab.JobEvent{Ref: "master"}},
			gitlab.KindBranch,
		},
		"ref:tag": {
			gitlab.JobEvent{JobEvent: gogitlab.JobEvent{Ref: "v1.0.0", Tag: true}},
			gitlab.KindTag,
		},
		"ref:merge_request": {
			gitlab.JobEvent{JobEvent: gogitlab.JobEvent{Ref: "refs/merge-requests/1029/merge"}},
			gitlab.KindMergeRequest,
		},
	}

	for name, tcase := range tcases {
		tcase := tcase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tcase.expected, tcase.event.RefKind())
		})
	}
}
