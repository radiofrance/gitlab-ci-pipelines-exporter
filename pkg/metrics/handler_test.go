package metrics_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/radiofrance/gitlab-ci-pipelines-exporter/pkg/metrics"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeCollectors struct{}

func (c fakeCollectors) MustRegister() {
	fakeMetric := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "fake_total",
			Help: "A fake counter metric",
		},
		[]string{"label"},
	)
	fakeMetric.With(prometheus.Labels{"label": "value"}).Inc()

	prometheus.MustRegister(fakeMetric)
}

func Test_Handler_ServeHTTP(t *testing.T) {
	collectors := &fakeCollectors{}
	handler := metrics.NewHandler("/metrics", collectors)
	server := httptest.NewServer(handler)

	response, err := server.Client().Get(fmt.Sprintf("%s/metrics", server.URL))
	require.NoError(t, err)

	body, err := io.ReadAll(response.Body)
	require.NoError(t, err)

	expectedBody := `# HELP fake_total A fake counter metric
# TYPE fake_total counter
fake_total{label="value"} 1
`
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Contains(t, string(body), expectedBody)
}
