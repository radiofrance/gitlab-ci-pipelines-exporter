package webhook

import "go.uber.org/zap"

type (
	Option func(webhook *Webhook)
)

// WithTimestamp allows to override the default the timestamp function.
func WithTimestamp(timestamp TimestampFunc) Option {
	return func(webhook *Webhook) {
		webhook.timestamp = timestamp
	}
}

// WithZapLogger configures the webhook with a preconfigured zap instance.
func WithZapLogger(logger *zap.Logger) Option {
	return func(webhook *Webhook) {
		if logger != nil {
			webhook.log = logger
		}
	}
}
