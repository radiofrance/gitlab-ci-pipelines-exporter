package webhook //nolint: testpackage

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWithTimeStamp(t *testing.T) {
	t.Parallel()

	timestamp := func() int64 { return 12345 }
	webhook := NewWebhook("", nil, WithTimestamp(timestamp))
	assert.Equal(t, int64(12345), webhook.timestamp())
}

func TestWithZapLogger(t *testing.T) {
	t.Parallel()

	logger := zap.NewNop()
	webhook := NewWebhook("", nil, WithZapLogger(logger))
	assert.Equal(t, logger, webhook.log)
}
