package webhook_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/radiofrance/gitlab-ci-pipelines-exporter/pkg/webhook"
	"github.com/stretchr/testify/assert"
)

func TestNewGitlabSecretTokenMiddleware(t *testing.T) {
	t.Parallel()

	next := http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) { writer.WriteHeader(http.StatusOK) })
	middleware := webhook.NewGitlabSecretTokenMiddleware("token")

	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/", nil)
	recorder := httptest.NewRecorder()
	middleware.ServeHTTP(recorder, req, next)
	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	assert.JSONEq(t, `{"error":"invalid Gitlab webhook secret token"}`, recorder.Body.String())

	req.Header.Add("X-Gitlab-Token", "token")

	recorder = httptest.NewRecorder()
	middleware.ServeHTTP(recorder, req, next)
	assert.Equal(t, http.StatusOK, recorder.Code)
}
