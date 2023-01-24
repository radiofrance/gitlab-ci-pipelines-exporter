package webhook

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWithZapLogger(t *testing.T) {
	logger := zap.NewNop()
	webhook := NewWebhook("", nil, WithZapLogger(logger))
	assert.Equal(t, logger, webhook.log)
}
