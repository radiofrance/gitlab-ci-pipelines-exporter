package webhook

import "go.uber.org/zap"

type (
	Option func(webhook *Webhook)
)

// WithZapLogger configures the webhook with a preconfigured zap instance.
func WithZapLogger(logger *zap.Logger) Option {
	return func(webhook *Webhook) {
		if logger != nil {
			webhook.log = logger
		}
	}
}
