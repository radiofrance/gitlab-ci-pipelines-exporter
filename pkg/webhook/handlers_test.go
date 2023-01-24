package webhook

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/radiofrance/gitlab-ci-pipelines-exporter/pkg/metrics"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

type PipelineHandlerTestSuite struct {
	suite.Suite
	webhook *Webhook
}

func (suite *PipelineHandlerTestSuite) SetupSuite() {
	timestamp = func() float64 { return 0 }
}

func (suite *PipelineHandlerTestSuite) SetupTest() {
	suite.webhook = &Webhook{
		collectors: metrics.NewPrometheusCollectors(),
		log:        zap.NewNop(),
	}
}

func (suite *PipelineHandlerTestSuite) TestSingleEvent() {
	genericTestWebhookHandler(
		suite.T(), suite.webhook.handlePipelineEvent,

		[]string{
			`{"object_attributes":{"id":4223,"status":"created"},"merge_request":{"iid":9832},"project":{"path_with_namespace": "radiofrance/gitlab-ci-pipelines-exporter"}}`,
		},

		map[prometheus.Collector]string{
			suite.webhook.collectors.ID: `
# HELP gitlab_ci_pipeline_id ID of the most recent pipeline
# TYPE gitlab_ci_pipeline_id gauge
gitlab_ci_pipeline_id{kind="merge_request",project="radiofrance/gitlab-ci-pipelines-exporter",ref="9832"} 4223
`,
			suite.webhook.collectors.Timestamp: `
# HELP gitlab_ci_pipeline_timestamp Timestamp of the last update of the most recent pipeline
# TYPE gitlab_ci_pipeline_timestamp gauge
gitlab_ci_pipeline_timestamp{kind="merge_request",project="radiofrance/gitlab-ci-pipelines-exporter",ref="9832"} 0
`,
			suite.webhook.collectors.QueuedDurationSeconds: ``,
			suite.webhook.collectors.DurationSeconds:       ``,
			suite.webhook.collectors.RunCount:              ``,
			suite.webhook.collectors.Status: `
# HELP gitlab_ci_pipeline_status Status of the most recent pipeline
# TYPE gitlab_ci_pipeline_status gauge
gitlab_ci_pipeline_status{kind="merge_request",project="radiofrance/gitlab-ci-pipelines-exporter",ref="9832",status="canceled"} 0
gitlab_ci_pipeline_status{kind="merge_request",project="radiofrance/gitlab-ci-pipelines-exporter",ref="9832",status="created"} 1
gitlab_ci_pipeline_status{kind="merge_request",project="radiofrance/gitlab-ci-pipelines-exporter",ref="9832",status="failed"} 0
gitlab_ci_pipeline_status{kind="merge_request",project="radiofrance/gitlab-ci-pipelines-exporter",ref="9832",status="manual"} 0
gitlab_ci_pipeline_status{kind="merge_request",project="radiofrance/gitlab-ci-pipelines-exporter",ref="9832",status="pending"} 0
gitlab_ci_pipeline_status{kind="merge_request",project="radiofrance/gitlab-ci-pipelines-exporter",ref="9832",status="preparing"} 0
gitlab_ci_pipeline_status{kind="merge_request",project="radiofrance/gitlab-ci-pipelines-exporter",ref="9832",status="running"} 0
gitlab_ci_pipeline_status{kind="merge_request",project="radiofrance/gitlab-ci-pipelines-exporter",ref="9832",status="scheduled"} 0
gitlab_ci_pipeline_status{kind="merge_request",project="radiofrance/gitlab-ci-pipelines-exporter",ref="9832",status="skipped"} 0
gitlab_ci_pipeline_status{kind="merge_request",project="radiofrance/gitlab-ci-pipelines-exporter",ref="9832",status="success"} 0
gitlab_ci_pipeline_status{kind="merge_request",project="radiofrance/gitlab-ci-pipelines-exporter",ref="9832",status="waiting_for_resource"} 0
`,
		},
	)
}

func (suite *PipelineHandlerTestSuite) TestMultipleEvents() {
	genericTestWebhookHandler(
		suite.T(), suite.webhook.handlePipelineEvent,

		[]string{
			`{"object_attributes":{"id":4223,"ref":"master","status":"created"},"project":{"path_with_namespace": "radiofrance/gitlab-ci-pipelines-exporter"}}`,
			`{"object_attributes":{"id":4224,"ref":"master","status":"pending"},"project":{"path_with_namespace": "radiofrance/gitlab-ci-pipelines-exporter"}}`,
			`{"object_attributes":{"id":4225,"ref":"master","status":"running","queued_duration":1.223745},"project":{"path_with_namespace": "radiofrance/gitlab-ci-pipelines-exporter"}}`,
			`{"object_attributes":{"id":4226,"ref":"master","status":"failed","duration":0.088020544,"queued_duration":1.223745},"project":{"path_with_namespace": "radiofrance/gitlab-ci-pipelines-exporter"}}`,
			`{"object_attributes":{"id":4227,"ref":"master","status":"running","queued_duration":2.343087599},"project":{"path_with_namespace": "radiofrance/gitlab-ci-pipelines-exporter"}}`,
			`{"object_attributes":{"id":4228,"ref":"master","status":"success","duration":91.140114,"queued_duration":2.343087599},"project":{"path_with_namespace": "radiofrance/gitlab-ci-pipelines-exporter"}}`,
			`{"object_attributes":{"id":4229,"ref":"v1.0.0","tag":true,"status":"created"},"project":{"path_with_namespace": "radiofrance/gitlab-ci-pipelines-exporter"}}`,
		},

		map[prometheus.Collector]string{
			suite.webhook.collectors.ID: `
# HELP gitlab_ci_pipeline_id ID of the most recent pipeline
# TYPE gitlab_ci_pipeline_id gauge
gitlab_ci_pipeline_id{kind="branch",project="radiofrance/gitlab-ci-pipelines-exporter",ref="master"} 4228
gitlab_ci_pipeline_id{kind="tag",project="radiofrance/gitlab-ci-pipelines-exporter",ref="v1.0.0"} 4229
`,
			suite.webhook.collectors.Timestamp: `
# HELP gitlab_ci_pipeline_timestamp Timestamp of the last update of the most recent pipeline
# TYPE gitlab_ci_pipeline_timestamp gauge
gitlab_ci_pipeline_timestamp{kind="branch",project="radiofrance/gitlab-ci-pipelines-exporter",ref="master"} 0
gitlab_ci_pipeline_timestamp{kind="tag",project="radiofrance/gitlab-ci-pipelines-exporter",ref="v1.0.0"} 0
`,
			suite.webhook.collectors.QueuedDurationSeconds: `
# HELP gitlab_ci_pipeline_queued_duration_seconds Duration in seconds the most recent pipeline has been queued before starting
# TYPE gitlab_ci_pipeline_queued_duration_seconds gauge
gitlab_ci_pipeline_queued_duration_seconds{kind="branch",project="radiofrance/gitlab-ci-pipelines-exporter",ref="master"} 2.343087599
`,
			suite.webhook.collectors.DurationSeconds: `
# HELP gitlab_ci_pipeline_duration_seconds Duration in seconds of the most recent pipeline
# TYPE gitlab_ci_pipeline_duration_seconds gauge
gitlab_ci_pipeline_duration_seconds{kind="branch",project="radiofrance/gitlab-ci-pipelines-exporter",ref="master"} 91.140114
`,
			suite.webhook.collectors.RunCount: `
# HELP gitlab_ci_pipeline_run_count Number of executions of a pipeline
# TYPE gitlab_ci_pipeline_run_count counter
gitlab_ci_pipeline_run_count{kind="branch",project="radiofrance/gitlab-ci-pipelines-exporter",ref="master"} 2
`,
			suite.webhook.collectors.Status: `
# HELP gitlab_ci_pipeline_status Status of the most recent pipeline
# TYPE gitlab_ci_pipeline_status gauge
gitlab_ci_pipeline_status{kind="branch",project="radiofrance/gitlab-ci-pipelines-exporter",ref="master",status="canceled"} 0
gitlab_ci_pipeline_status{kind="branch",project="radiofrance/gitlab-ci-pipelines-exporter",ref="master",status="created"} 0
gitlab_ci_pipeline_status{kind="branch",project="radiofrance/gitlab-ci-pipelines-exporter",ref="master",status="failed"} 0
gitlab_ci_pipeline_status{kind="branch",project="radiofrance/gitlab-ci-pipelines-exporter",ref="master",status="manual"} 0
gitlab_ci_pipeline_status{kind="branch",project="radiofrance/gitlab-ci-pipelines-exporter",ref="master",status="pending"} 0
gitlab_ci_pipeline_status{kind="branch",project="radiofrance/gitlab-ci-pipelines-exporter",ref="master",status="preparing"} 0
gitlab_ci_pipeline_status{kind="branch",project="radiofrance/gitlab-ci-pipelines-exporter",ref="master",status="running"} 0
gitlab_ci_pipeline_status{kind="branch",project="radiofrance/gitlab-ci-pipelines-exporter",ref="master",status="scheduled"} 0
gitlab_ci_pipeline_status{kind="branch",project="radiofrance/gitlab-ci-pipelines-exporter",ref="master",status="skipped"} 0
gitlab_ci_pipeline_status{kind="branch",project="radiofrance/gitlab-ci-pipelines-exporter",ref="master",status="success"} 1
gitlab_ci_pipeline_status{kind="branch",project="radiofrance/gitlab-ci-pipelines-exporter",ref="master",status="waiting_for_resource"} 0
gitlab_ci_pipeline_status{kind="tag",project="radiofrance/gitlab-ci-pipelines-exporter",ref="v1.0.0",status="canceled"} 0
gitlab_ci_pipeline_status{kind="tag",project="radiofrance/gitlab-ci-pipelines-exporter",ref="v1.0.0",status="created"} 1
gitlab_ci_pipeline_status{kind="tag",project="radiofrance/gitlab-ci-pipelines-exporter",ref="v1.0.0",status="failed"} 0
gitlab_ci_pipeline_status{kind="tag",project="radiofrance/gitlab-ci-pipelines-exporter",ref="v1.0.0",status="manual"} 0
gitlab_ci_pipeline_status{kind="tag",project="radiofrance/gitlab-ci-pipelines-exporter",ref="v1.0.0",status="pending"} 0
gitlab_ci_pipeline_status{kind="tag",project="radiofrance/gitlab-ci-pipelines-exporter",ref="v1.0.0",status="preparing"} 0
gitlab_ci_pipeline_status{kind="tag",project="radiofrance/gitlab-ci-pipelines-exporter",ref="v1.0.0",status="running"} 0
gitlab_ci_pipeline_status{kind="tag",project="radiofrance/gitlab-ci-pipelines-exporter",ref="v1.0.0",status="scheduled"} 0
gitlab_ci_pipeline_status{kind="tag",project="radiofrance/gitlab-ci-pipelines-exporter",ref="v1.0.0",status="skipped"} 0
gitlab_ci_pipeline_status{kind="tag",project="radiofrance/gitlab-ci-pipelines-exporter",ref="v1.0.0",status="success"} 0
gitlab_ci_pipeline_status{kind="tag",project="radiofrance/gitlab-ci-pipelines-exporter",ref="v1.0.0",status="waiting_for_resource"} 0
`,
		},
	)
}

func TestPipelineHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(PipelineHandlerTestSuite))
}

type JobHandlerTestSuite struct {
	suite.Suite
	webhook *Webhook
}

func (suite *JobHandlerTestSuite) SetupSuite() {
	timestamp = func() float64 { return 0 }
}

func (suite *JobHandlerTestSuite) SetupTest() {
	suite.webhook = &Webhook{
		collectors: metrics.NewPrometheusCollectors(),
		log:        zap.NewNop(),
	}
}

func (suite *JobHandlerTestSuite) TestSingleEvent() {
	genericTestWebhookHandler(
		suite.T(), suite.webhook.handleJobEvent,
		[]string{
			`{"ref":"refs/merge-requests/9832/merge","build_id":0,"build_name":"golang-ci-lint","build_stage":"quality","build_status":"created","project_name":"radiofrance / gitlab-ci-pipelines-exporter"}`,
		},
		map[prometheus.Collector]string{
			suite.webhook.collectors.JobID: `
# HELP gitlab_ci_pipeline_job_id ID of the most recent job
# TYPE gitlab_ci_pipeline_job_id gauge
gitlab_ci_pipeline_job_id{job_name="golang-ci-lint",kind="merge_request",project="radiofrance/gitlab-ci-pipelines-exporter",ref="9832",stage="quality"} 0
`,
			suite.webhook.collectors.JobTimestamp: `
# HELP gitlab_ci_pipeline_job_timestamp Creation date timestamp of the the most recent job
# TYPE gitlab_ci_pipeline_job_timestamp gauge
gitlab_ci_pipeline_job_timestamp{job_name="golang-ci-lint",kind="merge_request",project="radiofrance/gitlab-ci-pipelines-exporter",ref="9832",stage="quality"} 0
`,
			suite.webhook.collectors.JobDurationSeconds:       ``,
			suite.webhook.collectors.JobQueuedDurationSeconds: ``,
			suite.webhook.collectors.JobRunCount:              ``,
			suite.webhook.collectors.JobStatus: `
# HELP gitlab_ci_pipeline_job_status Status of the most recent job
# TYPE gitlab_ci_pipeline_job_status gauge
gitlab_ci_pipeline_job_status{job_name="golang-ci-lint",kind="merge_request",project="radiofrance/gitlab-ci-pipelines-exporter",ref="9832",stage="quality",status="canceled"} 0
gitlab_ci_pipeline_job_status{job_name="golang-ci-lint",kind="merge_request",project="radiofrance/gitlab-ci-pipelines-exporter",ref="9832",stage="quality",status="created"} 1
gitlab_ci_pipeline_job_status{job_name="golang-ci-lint",kind="merge_request",project="radiofrance/gitlab-ci-pipelines-exporter",ref="9832",stage="quality",status="failed"} 0
gitlab_ci_pipeline_job_status{job_name="golang-ci-lint",kind="merge_request",project="radiofrance/gitlab-ci-pipelines-exporter",ref="9832",stage="quality",status="manual"} 0
gitlab_ci_pipeline_job_status{job_name="golang-ci-lint",kind="merge_request",project="radiofrance/gitlab-ci-pipelines-exporter",ref="9832",stage="quality",status="pending"} 0
gitlab_ci_pipeline_job_status{job_name="golang-ci-lint",kind="merge_request",project="radiofrance/gitlab-ci-pipelines-exporter",ref="9832",stage="quality",status="preparing"} 0
gitlab_ci_pipeline_job_status{job_name="golang-ci-lint",kind="merge_request",project="radiofrance/gitlab-ci-pipelines-exporter",ref="9832",stage="quality",status="running"} 0
gitlab_ci_pipeline_job_status{job_name="golang-ci-lint",kind="merge_request",project="radiofrance/gitlab-ci-pipelines-exporter",ref="9832",stage="quality",status="scheduled"} 0
gitlab_ci_pipeline_job_status{job_name="golang-ci-lint",kind="merge_request",project="radiofrance/gitlab-ci-pipelines-exporter",ref="9832",stage="quality",status="skipped"} 0
gitlab_ci_pipeline_job_status{job_name="golang-ci-lint",kind="merge_request",project="radiofrance/gitlab-ci-pipelines-exporter",ref="9832",stage="quality",status="success"} 0
gitlab_ci_pipeline_job_status{job_name="golang-ci-lint",kind="merge_request",project="radiofrance/gitlab-ci-pipelines-exporter",ref="9832",stage="quality",status="waiting_for_resource"} 0
`,
		},
	)
}

func (suite *JobHandlerTestSuite) TestMultipleEvents() {
	genericTestWebhookHandler(
		suite.T(), suite.webhook.handleJobEvent,
		[]string{
			`{"ref":"master","build_id":4229,"build_name":"golang-ci-lint","build_stage":"quality","build_status":"created","project_name":"radiofrance / gitlab-ci-pipelines-exporter"}`,
			`{"ref":"master","build_id":4230,"build_name":"golang-ci-lint","build_stage":"quality","build_status":"pending","project_name":"radiofrance / gitlab-ci-pipelines-exporter"}`,
			`{"ref":"master","build_id":4231,"build_name":"golang-ci-lint","build_stage":"quality","build_status":"running","build_queued_duration":0.169516333,"project_name":"radiofrance / gitlab-ci-pipelines-exporter","runner":{"description":"gitlab-runner-standard-6f49f9-4b4qn"}}`,
			`{"ref":"master","build_id":4232,"build_name":"golang-ci-lint","build_stage":"quality","build_status":"failed","build_duration":1.223745182,"build_queued_duration":0.169516333,"project_name":"radiofrance / gitlab-ci-pipelines-exporter","runner":{"description":"gitlab-runner-standard-6f49f9-4b4qn"}}`,
			`{"ref":"master","build_id":4233,"build_name":"golang-ci-lint","build_stage":"quality","build_status":"running","build_queued_duration":0.074549967,"project_name":"radiofrance / gitlab-ci-pipelines-exporter","runner":{"description":"gitlab-runner-standard-6f49f9-4b4qn"}}`,
			`{"ref":"master","build_id":4234,"build_name":"golang-ci-lint","build_stage":"quality","build_status":"success","build_duration":196.868193,"build_queued_duration":0.074549967,"project_name":"radiofrance / gitlab-ci-pipelines-exporter","runner":{"description":"gitlab-runner-standard-6f49f9-4b4qn"}}`,
			`{"ref":"master","build_id":4235,"build_name":"build-image","build_stage":"build","build_status":"created","project_name":"radiofrance / gitlab-ci-pipelines-exporter"}`,
		},
		map[prometheus.Collector]string{
			suite.webhook.collectors.JobID: `
# HELP gitlab_ci_pipeline_job_id ID of the most recent job
# TYPE gitlab_ci_pipeline_job_id gauge
gitlab_ci_pipeline_job_id{job_name="build-image",kind="branch",project="radiofrance/gitlab-ci-pipelines-exporter",ref="master",stage="build"} 4235
gitlab_ci_pipeline_job_id{job_name="golang-ci-lint",kind="branch",project="radiofrance/gitlab-ci-pipelines-exporter",ref="master",stage="quality"} 4234
`,
			suite.webhook.collectors.JobTimestamp: `
# HELP gitlab_ci_pipeline_job_timestamp Creation date timestamp of the the most recent job
# TYPE gitlab_ci_pipeline_job_timestamp gauge
gitlab_ci_pipeline_job_timestamp{job_name="build-image",kind="branch",project="radiofrance/gitlab-ci-pipelines-exporter",ref="master",stage="build"} 0
gitlab_ci_pipeline_job_timestamp{job_name="golang-ci-lint",kind="branch",project="radiofrance/gitlab-ci-pipelines-exporter",ref="master",stage="quality"} 0
`,
			suite.webhook.collectors.JobDurationSeconds: `
# HELP gitlab_ci_pipeline_job_duration_seconds Duration in seconds of the most recent job
# TYPE gitlab_ci_pipeline_job_duration_seconds gauge
gitlab_ci_pipeline_job_duration_seconds{job_name="golang-ci-lint",kind="branch",project="radiofrance/gitlab-ci-pipelines-exporter",ref="master",stage="quality"} 196.868193
`,
			suite.webhook.collectors.JobQueuedDurationSeconds: `
# HELP gitlab_ci_pipeline_job_queued_duration_seconds Duration in seconds the most recent job has been queued before starting
# TYPE gitlab_ci_pipeline_job_queued_duration_seconds gauge
gitlab_ci_pipeline_job_queued_duration_seconds{job_name="golang-ci-lint",kind="branch",project="radiofrance/gitlab-ci-pipelines-exporter",ref="master",stage="quality"} 0.074549967
`,
			suite.webhook.collectors.JobRunCount: `
# HELP gitlab_ci_pipeline_job_run_count Number of executions of a job
# TYPE gitlab_ci_pipeline_job_run_count counter
gitlab_ci_pipeline_job_run_count{job_name="golang-ci-lint",kind="branch",project="radiofrance/gitlab-ci-pipelines-exporter",ref="master",stage="quality"} 2
`,
			suite.webhook.collectors.JobStatus: `
# HELP gitlab_ci_pipeline_job_status Status of the most recent job
# TYPE gitlab_ci_pipeline_job_status gauge
gitlab_ci_pipeline_job_status{job_name="build-image",kind="branch",project="radiofrance/gitlab-ci-pipelines-exporter",ref="master",stage="build",status="canceled"} 0
gitlab_ci_pipeline_job_status{job_name="build-image",kind="branch",project="radiofrance/gitlab-ci-pipelines-exporter",ref="master",stage="build",status="created"} 1
gitlab_ci_pipeline_job_status{job_name="build-image",kind="branch",project="radiofrance/gitlab-ci-pipelines-exporter",ref="master",stage="build",status="failed"} 0
gitlab_ci_pipeline_job_status{job_name="build-image",kind="branch",project="radiofrance/gitlab-ci-pipelines-exporter",ref="master",stage="build",status="manual"} 0
gitlab_ci_pipeline_job_status{job_name="build-image",kind="branch",project="radiofrance/gitlab-ci-pipelines-exporter",ref="master",stage="build",status="pending"} 0
gitlab_ci_pipeline_job_status{job_name="build-image",kind="branch",project="radiofrance/gitlab-ci-pipelines-exporter",ref="master",stage="build",status="preparing"} 0
gitlab_ci_pipeline_job_status{job_name="build-image",kind="branch",project="radiofrance/gitlab-ci-pipelines-exporter",ref="master",stage="build",status="running"} 0
gitlab_ci_pipeline_job_status{job_name="build-image",kind="branch",project="radiofrance/gitlab-ci-pipelines-exporter",ref="master",stage="build",status="scheduled"} 0
gitlab_ci_pipeline_job_status{job_name="build-image",kind="branch",project="radiofrance/gitlab-ci-pipelines-exporter",ref="master",stage="build",status="skipped"} 0
gitlab_ci_pipeline_job_status{job_name="build-image",kind="branch",project="radiofrance/gitlab-ci-pipelines-exporter",ref="master",stage="build",status="success"} 0
gitlab_ci_pipeline_job_status{job_name="build-image",kind="branch",project="radiofrance/gitlab-ci-pipelines-exporter",ref="master",stage="build",status="waiting_for_resource"} 0
gitlab_ci_pipeline_job_status{job_name="golang-ci-lint",kind="branch",project="radiofrance/gitlab-ci-pipelines-exporter",ref="master",stage="quality",status="canceled"} 0
gitlab_ci_pipeline_job_status{job_name="golang-ci-lint",kind="branch",project="radiofrance/gitlab-ci-pipelines-exporter",ref="master",stage="quality",status="created"} 0
gitlab_ci_pipeline_job_status{job_name="golang-ci-lint",kind="branch",project="radiofrance/gitlab-ci-pipelines-exporter",ref="master",stage="quality",status="failed"} 0
gitlab_ci_pipeline_job_status{job_name="golang-ci-lint",kind="branch",project="radiofrance/gitlab-ci-pipelines-exporter",ref="master",stage="quality",status="manual"} 0
gitlab_ci_pipeline_job_status{job_name="golang-ci-lint",kind="branch",project="radiofrance/gitlab-ci-pipelines-exporter",ref="master",stage="quality",status="pending"} 0
gitlab_ci_pipeline_job_status{job_name="golang-ci-lint",kind="branch",project="radiofrance/gitlab-ci-pipelines-exporter",ref="master",stage="quality",status="preparing"} 0
gitlab_ci_pipeline_job_status{job_name="golang-ci-lint",kind="branch",project="radiofrance/gitlab-ci-pipelines-exporter",ref="master",stage="quality",status="running"} 0
gitlab_ci_pipeline_job_status{job_name="golang-ci-lint",kind="branch",project="radiofrance/gitlab-ci-pipelines-exporter",ref="master",stage="quality",status="scheduled"} 0
gitlab_ci_pipeline_job_status{job_name="golang-ci-lint",kind="branch",project="radiofrance/gitlab-ci-pipelines-exporter",ref="master",stage="quality",status="skipped"} 0
gitlab_ci_pipeline_job_status{job_name="golang-ci-lint",kind="branch",project="radiofrance/gitlab-ci-pipelines-exporter",ref="master",stage="quality",status="success"} 1
gitlab_ci_pipeline_job_status{job_name="golang-ci-lint",kind="branch",project="radiofrance/gitlab-ci-pipelines-exporter",ref="master",stage="quality",status="waiting_for_resource"} 0
`,
		},
	)
}

func TestJobHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(JobHandlerTestSuite))
}

func genericTestWebhookHandler[T any](t *testing.T, handler func(T) error, events []string, expected map[prometheus.Collector]string) {
	for _, str := range events {
		var event T

		err := json.Unmarshal([]byte(str), &event)
		require.NoError(t, err)

		err = handler(event)
		assert.NoError(t, err)
	}

	for collector, expect := range expected {
		err := testutil.CollectAndCompare(collector, bytes.NewBuffer([]byte(expect)))
		assert.NoError(t, err)
	}
}
