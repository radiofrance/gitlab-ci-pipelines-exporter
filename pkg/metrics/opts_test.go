package metrics

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWithZapLogger(t *testing.T) {
	logger := zap.NewNop()
	handler := NewHandler("/metrics", NewPrometheusCollectors(), WithZapLogger(logger))
	assert.Equal(t, logger, handler.log)
}
