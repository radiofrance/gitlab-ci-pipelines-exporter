package webhook

import (
	"net/http"

	gcpehttp "github.com/radiofrance/gitlab-ci-pipelines-exporter/pkg/http"
	"github.com/urfave/negroni"
)

// NewGitlabSecretTokenMiddleware rejects all requests using the wrong Gitlab secret token.
func NewGitlabSecretTokenMiddleware(token string) negroni.HandlerFunc {
	return func(writer http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
		xtoken := req.Header.Get("X-Gitlab-Token")
		if xtoken != token {
			// NOTE: this returns 500 in order to notify Gitlab to disable this webhook. See
			//			 https://docs.gitlab.com/ee/user/project/integrations/webhooks.html#configure-your-webhook-receiver-endpoint
			//			 for more details.
			writer.WriteHeader(http.StatusInternalServerError)
			gcpehttp.WriteError(writer, "invalid Gitlab webhook secret token")
			return
		}

		next(writer, req)
	}
}
