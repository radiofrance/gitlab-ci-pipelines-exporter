package webhook_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/radiofrance/gitlab-ci-pipelines-exporter/pkg/webhook"
	"github.com/stretchr/testify/assert"
)

func TestNewGitlabSecretTokenMiddleware(t *testing.T) {
	next := http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) { writer.WriteHeader(http.StatusOK) })
	mw := webhook.NewGitlabSecretTokenMiddleware("token")

	req, _ := http.NewRequest(http.MethodPost, "https://:::0", nil)
	w := httptest.NewRecorder()
	mw.ServeHTTP(w, req, next)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, `{"error":"invalid Gitlab webhook secret token"}`, w.Body.String())

	req.Header.Add("X-Gitlab-Token", "token")
	w = httptest.NewRecorder()
	mw.ServeHTTP(w, req, next)
	assert.Equal(t, http.StatusOK, w.Code)
}
