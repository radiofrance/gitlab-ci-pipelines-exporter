package metrics //nolint:testpackage

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWithZapLogger(t *testing.T) {
	t.Parallel()

	logger := zap.NewNop()
	handler := NewHandler("/metrics", NewPrometheusCollectors(), WithZapLogger(logger))
	assert.Equal(t, logger, handler.log)
}
