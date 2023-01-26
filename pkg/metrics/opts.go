package metrics

import "go.uber.org/zap"

type (
	Option func(webhook *Handler)
)

// WithZapLogger configures the webhook with a preconfigured zap instance.
func WithZapLogger(logger *zap.Logger) Option {
	return func(handler *Handler) {
		if logger != nil {
			handler.log = logger
		}
	}
}
