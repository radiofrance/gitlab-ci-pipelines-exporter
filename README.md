<!-- markdownlint-disable MD033 -->
<h1 align="center">
  <a href="https://gitlab.com">
    <img src="https://about.gitlab.com/images/press/logo/png/gitlab-logo-100.png" alt="Gitlab logo" height="50px">
  </a>
  &nbsp;&nbsp;+&nbsp;&nbsp;
  <a href="https://prometheus.io/">
    <img src="https://cncf-branding.netlify.app/img/projects/prometheus/horizontal/color/prometheus-horizontal-color.png" alt="Gitlab logo" height="50px">
  </a>
</h1>

<h4 align="center">gitlab-ci-pipelines-exporter - Monitor your <a href="https://docs.gitlab.com/ee/ci/pipelines/">
GitlabCI pipelines</a></h4>

<div align="center">
  <a href="https://github.com/radiofrance/gitlab-ci-pipelines-exporter/issues/new">Report a Bug</a> ·
  <a href="https://github.com/radiofrance/gitlab-ci-pipelines-exporter/issues/new">Request a Feature</a> ·
  <a href="https://github.com/radiofrance/gitlab-ci-pipelines-exporter/discussions">Ask a Question</a>
  <br/>
  <br/>

[![GoReportCard](https://goreportcard.com/badge/github.com/radiofrance/gitlab-ci-pipelines-exporter)](https://goreportcard.com/report/github.com/radiofrance/gitlab-ci-pipelines-exporter)
[![Codecov branch](https://img.shields.io/codecov/c/github/radiofrance/gitlab-ci-pipelines-exporter/main?label=code%20coverage)](https://app.codecov.io/gh/radiofrance/gitlab-ci-pipelines-exporter/tree/main)
[![GoDoc reference](https://img.shields.io/badge/godoc-reference-blue.svg)](https://pkg.go.dev/github.com/radiofrance/gitlab-ci-pipelines-exporter)
<br/>
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/radiofrance/gitlab-ci-pipelines-exporter?logo=go&logoColor=white&logoWidth=20)
![Gitlab compatible version](https://img.shields.io/badge/Gitlab_CI_compatibility-≥_15.4.2-success?logo=gitlab&logoColor=white&logoWidth=20)
[![License](https://img.shields.io/badge/license-CeCILL%202.1-blue?logo=git&logoColor=white&logoWidth=20)](LICENSE)

<a href="#about">About</a> ·
<a href="#install">How to Install?</a> ·
<a href="#exported-metrics">Metrics</a> ·
<a href="#support">Support</a> ·
<a href="#contributing">Contributing</a> ·
<a href="#security">Security</a>

</div>

---
<!-- markdownlint-enable MD033 -->

## About

`gitlab-ci-pipelines-exporter` allows you to monitor
your [GitLab CI pipelines](https://docs.gitlab.com/ee/ci/pipelines/)
with [Prometheus](https://prometheus.io/).
You can find more information on GitLab docs about how it takes part improving your pipeline efficiency.

This project is based on the more
featured [`gitlab-ci-pipelines-exporter`](https://github.com/mvisonneau/gitlab-ci-pipelines-exporter)
from [`mvisonneau`](https://github.com/mvisonneau), and metrics should be compatible with it.

### Why rewriting the exporter ?

The original project has many features and is very complete (it manages metrics
in [OpenMetrics](https://github.com/OpenObservability/OpenMetrics)
format, exports environment metrics, etc... we recommend to take a look before). However, it has a downside: it depends
on the Gitlab API to fetch data.  
Unfortunately, in our case, this results in either very aggressive rate limiting and not being able to get enough
metrics for our pipelines, or API saturation.

Using only the data provided by the webhook payload, we have enough information to export ~3/4 of the metrics without
relying on the Gitlab API.

## Install

### Go

```shell
go install github.com/radiofrance/gitlab-ci-pipelines-exporter/cmd/gitlab-ci-pipelines-exporter@latest
gitlab-ci-pipelines-exporter --gitlab.webhook-secret-token "GITLAB_SECRET_TOKEN"
```

### Docker

```shell
docker pull ghcr.io/radiofrance/gitlab-ci-pipelines-exporter
docker run --publish 8080 --publish 9252 ghcr.io/radiofrance/gitlab-ci-pipelines-exporter --gitlab.webhook-secret-token "GITLAB_SECRET_TOKEN"
```

### Helm

```shell
helm repo add radiofrance-gcpe https://radiofrance.github.io/gitlab-ci-pipelines-exporter
helm upgrade --install gitlab-ci-pipelines-exporter radiofrance-gcpe/gitlab-ci-pipelines-exporter \
  --namespace gitlab-ci-pipelines-exporter \
  --create-namespace \
  --wait \
  --set gitlab.webhook-secret-token="GITLAB_SECRET_TOKEN"
helm test gitlab-ci-pipelines-exporter --namespace gitlab-ci-pipelines-exporter
```

## Usage

```shell
NAME:
   gitlab-ci-pipelines-exporter - Export metrics about GitLab CI pipelines statuses

USAGE:
   gitlab-ci-pipelines-exporter [global options] command [command options] [arguments...]

VERSION:
   devel

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --web.listen-address value           address:port to listen on for incoming webhooks (default: ":8080") [$WEB_LISTEN_ADDRESS]
   --telemetry.listen-address value     address:port to listen on for telemetry (default: ":9252") [$TELEMETRY_LISTEN_ADDRESS]
   --telemetry.path value               Path under which to expose telemetry endpoint (default: "/metrics") [$TELEMETRY_PATH]
   --gitlab.webhook-secret-token token  token used to authenticate legitimate requests (overrides config file parameter) [$GITLAB_WEBHOOK_SECRET_TOKEN]
   --log.level value                    Log verbosity (default: "info") [$LOG_LEVEL]
   --help, -h                           show help (default: false)
   --version, -v                        print the version (default: false)
```

## Exported metrics

| Metric name                                      | Description                                                                  | Labels                                                                        |
|--------------------------------------------------|------------------------------------------------------------------------------|-------------------------------------------------------------------------------|
| `gitlab_ci_pipeline_duration_seconds`            | Duration in seconds of the most recent pipeline                              | [project], [ref], [kind]                                                      |
| `gitlab_ci_pipeline_id`                          | ID of the most recent pipeline                                               | [project], [ref], [kind]                                                      |
| `gitlab_ci_pipeline_job_duration_seconds`        | Duration in seconds of the most recent job                                   | [project], [ref], [runner_description], [kind], [stage], [job_name]           |
| `gitlab_ci_pipeline_job_id`                      | ID of the most recent job                                                    | [project], [ref], [runner_description], [kind], [stage], [job_name]           |
| `gitlab_ci_pipeline_job_queued_duration_seconds` | Duration in seconds the most recent job has been queued before starting      | [project], [ref], [runner_description], [kind], [stage], [job_name]           |
| `gitlab_ci_pipeline_job_run_count`               | Number of executions of a job                                                | [project], [ref], [runner_description], [kind], [stage], [job_name]           |
| `gitlab_ci_pipeline_job_status`                  | Status of the most recent job                                                | [project], [ref], [runner_description], [kind], [stage], [job_name], [status] |
| `gitlab_ci_pipeline_job_timestamp`               | Creation date timestamp of the the most recent job                           | [project], [ref], [runner_description], [kind], [stage], [job_name]           |
| `gitlab_ci_pipeline_queued_duration_seconds`     | Duration in seconds the most recent pipeline has been queued before starting | [project], [ref], [kind]                                                      |
| `gitlab_ci_pipeline_run_count`                   | Number of executions of a pipeline                                           | [project], [ref], [kind]                                                      |
| `gitlab_ci_pipeline_status`                      | Status of the most recent pipeline                                           | [project], [ref], [kind], [status]                                            |
| `gitlab_ci_pipeline_timestamp`                   | Timestamp of the last update of the most recent pipeline                     | [project], [ref], [kind]                                                      |

### Compatibility with [mvisonneau/gitlab-ci-pipelines-exporter](https://github.com/mvisonneau/gitlab-ci-pipelines-exporter)

Some metrics has been removed:

- `gcpe_currently_queued_tasks_count` _(Number of tasks in the queue)_: only related
  to [mvisonneau/gitlab-ci-pipelines-exporter](https://github.com/mvisonneau/gitlab-ci-pipelines-exporter)
- `gcpe_environments_count` _(Number of GitLab environments being exported)_: only related
  to [mvisonneau/gitlab-ci-pipelines-exporter](https://github.com/mvisonneau/gitlab-ci-pipelines-exporter)
- `gcpe_executed_tasks_count` _(Number of tasks executed)_: only related
  to [mvisonneau/gitlab-ci-pipelines-exporter](https://github.com/mvisonneau/gitlab-ci-pipelines-exporter)
- `gcpe_gitlab_api_requests_count` _(GitLab API requests count)_: only related
  to [mvisonneau/gitlab-ci-pipelines-exporter](https://github.com/mvisonneau/gitlab-ci-pipelines-exporter)
- `gcpe_gitlab_api_requests_remaining` _(GitLab API requests remaining in the API Limit)_: only related
  to [mvisonneau/gitlab-ci-pipelines-exporter](https://github.com/mvisonneau/gitlab-ci-pipelines-exporter)
- `gcpe_gitlab_api_requests_limit` _(GitLab API requests available in the API Limit)_: only related
  to [mvisonneau/gitlab-ci-pipelines-exporter](https://github.com/mvisonneau/gitlab-ci-pipelines-exporter)
- `gcpe_metrics_count` _(Number of GitLab pipelines metrics being exported)_: only related
  to [mvisonneau/gitlab-ci-pipelines-exporter](https://github.com/mvisonneau/gitlab-ci-pipelines-exporter)
- `gcpe_projects_count` _(Number of GitLab projects being exported)_: only related
  to [mvisonneau/gitlab-ci-pipelines-exporter](https://github.com/mvisonneau/gitlab-ci-pipelines-exporter)
- `gcpe_refs_count` _(Number of GitLab refs being exported)_: only related
  to [mvisonneau/gitlab-ci-pipelines-exporter](https://github.com/mvisonneau/gitlab-ci-pipelines-exporter)


- `gitlab_ci_environment_behind_commits_count` _(Number of commits the environment is behind given its last
  deployment)_: `environment`'s metrics are currently not managed, maybe for a future release
- `gitlab_ci_environment_behind_duration_seconds` _(Duration in seconds the environment is behind the most recent commit
  given its last deployment)_: `environment`'s metrics are currently not managed, maybe for a future release
- `gitlab_ci_environment_deployment_count` _(Number of deployments for an environment)_: `environment`'s metrics are
  currently not managed, maybe for a future release
- `gitlab_ci_environment_deployment_duration_seconds` _(Duration in seconds of the most recent deployment of the
  environment)_: `environment`'s metrics are currently not managed, maybe for a future release
- `gitlab_ci_environment_deployment_job_id` _(ID of the most recent deployment job for an environment)_: `environment`'s
  metrics are currently not managed, maybe for a future release
- `gitlab_ci_environment_deployment_status` _(Status of the most recent deployment of the environment)_: `environment`'s
  metrics are currently not managed, maybe for a future release
- `gitlab_ci_environment_deployment_timestamp` _(Creation date of the most recent deployment of the
  environment)_: `environment`'s metrics are currently not managed, maybe for a future release
- `gitlab_ci_environment_information` _(Information about the environment)_: `environment`'s metrics are currently not
  managed, maybe for a future release


- `gitlab_ci_pipeline_job_artifact_size_bytes` _(Artifact size in bytes (sum of all of them) of the most recent job)_:
  currently not managed, but `pipeline` events has this information into `builds[].artifacts_file.size`
- `gitlab_ci_pipeline_coverage` _(Coverage of the most recent pipeline)_: unfortunately, we didn't found any information
  about coverage inside `job` or `event` payload

We also remove two labels, that are not used
inside [Grafana dashboards](https://github.com/mvisonneau/gitlab-ci-pipelines-exporter#tldr):

- `topics` _(all metrics)_: we can't have this information inside webhook payloads
- `variables` _(all metrics)_: currently not handled (maybe for a future release)

## Support

Reach out to the maintainer at one of the following places:

- [GitHub Discussions](https://github.com/radiofrance/gitlab-ci-pipelines-exporter/discussions)
- Open an issue on [Github](https://github.com/radiofrance/gitlab-ci-pipelines-exporter/issues/new)

## Contributing

First off, thanks for taking the time to contribute! Contributions are what make the
open-source community such an amazing place to learn, inspire, and create. Any contributions
you make will benefit everybody else and are **greatly appreciated**.

Please read [our contribution guidelines](docs/CONTRIBUTING.md), and thank you for being involved!

## Security

`gitlab-ci-pipelines-exporter` follows good practices of security, but 100% security cannot be assured.
`gitlab-ci-pipelines-exporter` is provided **"as is"** without any **warranty**. Use at your own risk.

*For more information and to report security issues, please refer to our [security documentation](docs/SECURITY.md).*

## License

This project is licensed under the **CeCILL License 2.1**.

See [LICENSE](LICENSE) for more information.

## Acknowledgements

Thanks for these awesome resources and projects that were used during development:

- <https://github.com/mvisonneau/gitlab-ci-pipelines-exporter> - Original Gitlab CI pipelines exporter
- <https://github.com/uber-go/zap> - Blazing fast, structured, leveled logging
- <https://github.com/urfave/cli> - Simple and fast package for build the CLI
- <https://github.com/urfave/negroni> - Idiomatic approach to web middlewares
