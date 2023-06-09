package webhook_test

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/radiofrance/gitlab-ci-pipelines-exporter/pkg/metrics"
	"github.com/radiofrance/gitlab-ci-pipelines-exporter/pkg/webhook"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Webhook_ServeHTTP(t *testing.T) {
	t.Parallel()

	collectors := metrics.NewPrometheusCollectors()
	webhook := webhook.NewWebhook("secret_token", collectors)
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
		"Webhooks route responds 500 when not authenticated": {
			method:             http.MethodPost,
			uri:                "/hooks",
			expectedStatusCode: http.StatusInternalServerError,
			expectedBody:       `{"error":"invalid Gitlab webhook secret token"}`,
		},
		"Webhooks route responds 200 when authenticated": {
			method:             http.MethodPost,
			uri:                "/hooks",
			headers:            authenticatedHeaders,
			event:              `{"object_kind":"build"}`,
			expectedStatusCode: http.StatusOK,
		},
		"Webhooks route responds 405 when invalid method": {
			method:             http.MethodGet,
			uri:                "/hooks",
			headers:            authenticatedHeaders,
			expectedStatusCode: http.StatusMethodNotAllowed,
		},
		"Webhooks route responds 501 when object_kind unsupported": {
			method:             http.MethodPost,
			uri:                "/hooks",
			headers:            authenticatedHeaders,
			event:              `{"object_kind":"unknown"}`,
			expectedStatusCode: http.StatusNotImplemented,
		},
		"Webhooks route handles invalid event data": {
			method:             http.MethodPost,
			uri:                "/hooks",
			headers:            authenticatedHeaders,
			event:              `true`,
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `{"error":"json: cannot unmarshal bool into Go value of type struct { ObjectKind string \"json:\\\"object_kind\\\"\" }"}`, //nolint:lll
		},
		"Webhooks route handles invalid json payload": {
			method:             http.MethodPost,
			uri:                "/hooks",
			headers:            authenticatedHeaders,
			event:              `not a json`,
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `{"error":"invalid character 'o' in literal null (expecting 'u')"}`,
		},
		"Webhooks route handles invalid pipeline event": {
			method:             http.MethodPost,
			uri:                "/hooks",
			headers:            authenticatedHeaders,
			event:              `{"object_kind":"pipeline","object_attributes":["invalid"]}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `{"error":"failed to unmarshall json into *gitlab.PipelineEvent"}`,
		},
		"Webhooks route handles invalid job event": {
			method:             http.MethodPost,
			uri:                "/hooks",
			headers:            authenticatedHeaders,
			event:              `{"object_kind":"build","ref":["invalid"]}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `{"error":"failed to unmarshall json into *gitlab.JobEvent"}`,
		},
		"Webhooks route handles pipeline event": {
			method:             http.MethodPost,
			uri:                "/hooks",
			headers:            authenticatedHeaders,
			event:              `{"object_kind":"pipeline"}`,
			expectedStatusCode: http.StatusOK,
		},
		"Webhooks route handles job event": {
			method:             http.MethodPost,
			uri:                "/hooks",
			headers:            authenticatedHeaders,
			event:              `{"object_kind":"build"}`,
			expectedStatusCode: http.StatusOK,
		},
	}

	for name, tcase := range tcases {
		tcase := tcase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			requestURL, err := url.Parse(fmt.Sprintf("%s%s", server.URL, tcase.uri))
			require.NoError(t, err)

			request := httptest.NewRequest(tcase.method, server.URL, bytes.NewBuffer([]byte(tcase.event))) //nolint:mirror
			request.RequestURI = ""
			request.URL = requestURL
			request.Header = tcase.headers
			response, err := server.Client().Do(request)
			require.NoError(t, err)
			defer response.Body.Close()

			body, err := io.ReadAll(response.Body)
			require.NoError(t, err)

			assert.Equal(t, tcase.expectedStatusCode, response.StatusCode)
			assert.Equal(t, tcase.expectedBody, string(body))
		})
	}
}
