package webhook

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/radiofrance/gitlab-ci-pipelines-exporter/pkg/metrics"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Webhook_ServeHTTP(t *testing.T) {
	collectors := metrics.NewPrometheusCollectors()
	webhook := NewWebhook("secret_token", collectors)
	server := httptest.NewServer(webhook)

	authenticatedHeaders := http.Header{
		"X-Gitlab-Token": {"secret_token"},
	}

	tcases := map[string]struct {
		method  string
		uri     string
		headers map[string][]string
		event   string

		expectedStatusCode int
		expectedBody       string
	}{
		"Healthcheck route responds 200": {
			method:             http.MethodGet,
			uri:                "/healthz",
			expectedStatusCode: http.StatusOK,
		},
		"Pipelines route responds 500 when not authenticated": {
			method:             http.MethodPost,
			uri:                "/pipeline",
			expectedStatusCode: http.StatusInternalServerError,
			expectedBody:       `{"error":"invalid Gitlab webhook secret token"}`,
		},
		"Pipelines route responds 200 when authenticated": {
			method:             http.MethodPost,
			uri:                "/pipeline",
			headers:            authenticatedHeaders,
			event:              `{"foo":"bar"}`,
			expectedStatusCode: http.StatusOK,
		},
		"Pipeline route responds 405 when invalid method": {
			method:             http.MethodGet,
			uri:                "/pipeline",
			headers:            authenticatedHeaders,
			expectedStatusCode: http.StatusMethodNotAllowed,
		},
		"Pipeline route handles invalid event data": {
			method:             http.MethodPost,
			uri:                "/pipeline",
			headers:            authenticatedHeaders,
			event:              `true`,
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `{"error":"json: cannot unmarshal bool into Go value of type gitlab_events.PipelineEvent"}`,
		},
		"Pipeline route handles invalid json payload": {
			method:             http.MethodPost,
			uri:                "/pipeline",
			headers:            authenticatedHeaders,
			event:              `not a json`,
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `{"error":"invalid character 'o' in literal null (expecting 'u')"}`,
		},
		"Job route responds 500 when not authenticated": {
			method:             http.MethodPost,
			uri:                "/job",
			expectedStatusCode: http.StatusInternalServerError,
			expectedBody:       `{"error":"invalid Gitlab webhook secret token"}`,
		},
		"Job route responds 200 when authenticated": {
			method:             http.MethodPost,
			uri:                "/job",
			headers:            authenticatedHeaders,
			event:              `{"foo":"bar"}`,
			expectedStatusCode: http.StatusOK,
		},
		"Job route responds 405 when invalid method": {
			method:             http.MethodGet,
			uri:                "/job",
			headers:            authenticatedHeaders,
			expectedStatusCode: http.StatusMethodNotAllowed,
		},
		"Job route handles invalid event data": {
			method:             http.MethodPost,
			uri:                "/job",
			headers:            authenticatedHeaders,
			event:              `true`,
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `{"error":"json: cannot unmarshal bool into Go value of type gitlab_events.JobEvent"}`,
		},
		"Job route handles invalid json payload": {
			method:             http.MethodPost,
			uri:                "/job",
			headers:            authenticatedHeaders,
			event:              `not a json`,
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `{"error":"invalid character 'o' in literal null (expecting 'u')"}`,
		},
	}

	for name, tcase := range tcases {
		t.Run(name, func(t *testing.T) {
			requestURL, err := url.Parse(fmt.Sprintf("%s%s", server.URL, tcase.uri))
			require.NoError(t, err)

			request := httptest.NewRequest(tcase.method, server.URL, bytes.NewBuffer([]byte(tcase.event)))
			request.RequestURI = ""
			request.URL = requestURL
			request.Header = tcase.headers
			response, err := server.Client().Do(request)
			require.NoError(t, err)

			body, err := io.ReadAll(response.Body)
			require.NoError(t, err)

			assert.Equal(t, tcase.expectedStatusCode, response.StatusCode)
			assert.Equal(t, tcase.expectedBody, string(body))
		})
	}
}
