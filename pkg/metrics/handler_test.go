package metrics_test

import (
	"context"
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
	t.Parallel()

	collectors := &fakeCollectors{}
	handler := metrics.NewHandler("/metrics", collectors)
	server := httptest.NewServer(handler)

	url := fmt.Sprintf("%s/metrics", server.URL)
	request, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	require.NoError(t, err)
	response, err := server.Client().Do(request)
	require.NoError(t, err)
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	require.NoError(t, err)

	expectedBody := `# HELP fake_total A fake counter metric
# TYPE fake_total counter
fake_total{label="value"} 1
`
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Contains(t, string(body), expectedBody)
}
