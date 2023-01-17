package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/radiofrance/gitlab-ci-pipelines-exporter/pkg/metrics"
	"github.com/radiofrance/gitlab-ci-pipelines-exporter/pkg/webhook"

	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var version = "devel"

func main() {
	_, ctxCancel := context.WithCancel(context.Background())
	defer ctxCancel()

	app := cli.NewApp()
	app.Name = "gitlab-ci-pipelines-exporter"
	app.Usage = "Export metrics about GitLab CI pipelines statuses"
	app.Version = version
	app.EnableBashCompletion = true

	app.Flags = cli.FlagsByName{
		&cli.StringFlag{
			Name:    "web.listen-address",
			Usage:   "address:port to listen on for telemetry",
			Value:   ":8080",
			EnvVars: []string{"WEB_LISTEN_ADDRESS"},
		},
		&cli.StringFlag{
			Name:    "telemetry.listen-address",
			Usage:   "Path under which to expose metrics",
			Value:   ":9252",
			EnvVars: []string{"TELEMETRY_LISTEN_ADDRESS"},
		},
		&cli.StringFlag{
			Name:    "telemetry.path",
			Usage:   "Path under which to expose metrics",
			Value:   "/metrics",
			EnvVars: []string{"TELEMETRY_PATH"},
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

		// Graceful shutdowns
		onShutdown := make(chan os.Signal, 1)
		signal.Notify(onShutdown, syscall.SIGINT, syscall.SIGTERM, syscall.SIGABRT)

		webhookHandler := webhook.NewWebhook(
			ctx.String("gitlab.webhook-secret-token"),
			webhook.WithZapLogger(logger),
		)

		webhookSrv := &http.Server{
			Addr:    ctx.String("web.listen-address"),
			Handler: webhookHandler,
		}

		go func() {
			logger.Info(fmt.Sprintf("webhook server listening on %s", ctx.String("web.listen-address")))
			if err := webhookSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				logger.Fatal("")
			}
		}()

		metricsSrv := &http.Server{
			Addr: ctx.String("telemetry.listen-address"),
			Handler: metrics.NewHandler(
				ctx.String("telemetry.path"),
				webhookHandler.Collectors,
				metrics.WithZapLogger(logger),
			),
		}

		go func() {
			logger.Info(fmt.Sprintf("metrics server listening on %s", ctx.String("telemetry.listen-address")))
			if err := metricsSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				logger.Fatal("")
			}
		}()

		<-onShutdown
		logger.Info("received exit signal, gracefully shutting down...")
		ctxCancel()

		httpServerContext, forceHTTPServerShutdown := context.WithTimeout(context.Background(), 5*time.Second)
		defer forceHTTPServerShutdown()

		if err := webhookSrv.Shutdown(httpServerContext); err != nil {
			return err
		}
		if err := metricsSrv.Shutdown(httpServerContext); err != nil {
			return err
		}

		return nil
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
}
