package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/xunleii/gitlab-ci-pipelines-exporter/pkg/webhook"

	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var version = "devel"

func main() {
	app := cli.NewApp()
	app.Name = "gitlab-ci-pipelines-exporter"
	app.Usage = "Export metrics about GitLab CI pipelines statuses"
	app.Version = version
	app.EnableBashCompletion = true

	app.Flags = cli.FlagsByName{
		&cli.StringFlag{
			Name:    "web.listen-address",
			Usage:   "address:port to listen on for telemetry",
			Value:   ":9252",
			EnvVars: []string{"WEB_LISTEN_ADDRESS"},
		},
		&cli.StringFlag{
			Name:    "web.telemetry-path",
			Usage:   "Path under which to expose metrics",
			Value:   "/metrics",
			EnvVars: []string{"WEB_TELEMETRY_PATH"},
		},
		&cli.StringFlag{
			Name:     "gitlab.webhook-secret-token",
			Usage:    "`token` used to authenticate legitimate requests (overrides config file parameter)",
			Required: true,
			EnvVars:  []string{"GITLAB_WEBHOOK_SECRET_TOKEN"},
		},
		&cli.StringFlag{
			Name:    "log.level",
			Usage:   "Log verbosity",
			Value:   "info",
			EnvVars: []string{"LOG_LEVEL"},
		},
		&cli.BoolFlag{
			Name:   "dev",
			Hidden: true,
		},
	}

	app.Action = func(ctx *cli.Context) error {
		lvl, err := zapcore.ParseLevel(ctx.String("log.level"))
		if err != nil {
			return fmt.Errorf("failed to parse log.level: %w", err)
		}

		var logger *zap.Logger
		if ctx.Bool("dev") {
			logger, err = zap.NewDevelopment(zap.IncreaseLevel(lvl))
		} else {
			logger, err = zap.Config{
				Level:       zap.NewAtomicLevelAt(lvl),
				Development: false,
				Sampling: &zap.SamplingConfig{
					Initial:    100,
					Thereafter: 100,
				},
				Encoding: "json",
				EncoderConfig: zapcore.EncoderConfig{
					TimeKey:        "ts",
					LevelKey:       "level",
					NameKey:        "logger",
					CallerKey:      "caller",
					FunctionKey:    zapcore.OmitKey,
					MessageKey:     "msg",
					StacktraceKey:  "stacktrace",
					LineEnding:     zapcore.DefaultLineEnding,
					EncodeLevel:    zapcore.LowercaseLevelEncoder,
					EncodeTime:     zapcore.RFC3339NanoTimeEncoder,
					EncodeDuration: zapcore.SecondsDurationEncoder,
					EncodeCaller:   zapcore.ShortCallerEncoder,
				},
				OutputPaths:      []string{"stderr"},
				ErrorOutputPaths: []string{"stderr"},
			}.Build()
		}
		if err != nil {
			return err
		}

		logger.Info(fmt.Sprintf("start listening on %s", ctx.String("web.listen-address")))
		return http.ListenAndServe(
			ctx.String("web.listen-address"),
			webhook.NewWebhook(
				ctx.String("web.telemetry-path"),
				ctx.String("gitlab.webhook-secret-token"),
				webhook.WithZapLogger(logger),
			),
		)
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
}
