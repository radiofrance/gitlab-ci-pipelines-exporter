package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	gcpehttp "github.com/radiofrance/gitlab-ci-pipelines-exporter/pkg/http"
	"github.com/urfave/negroni"
	"go.uber.org/zap"
)

type (
	Collectors interface {
		MustRegister()
	}

	Handler struct {
		collectors Collectors
		mux        *http.ServeMux
		log        *zap.Logger
	}
)

// NewHandler creates a new instance of a Handler object.
func NewHandler(pattern string, collectors Collectors, opts ...Option) *Handler {
	handler := &Handler{
		collectors: collectors,
		mux:        http.NewServeMux(),
		log:        zap.Must(zap.NewProduction()),
	}

	collectors.MustRegister()

	for _, opt := range opts {
		opt(handler)
	}

	root := negroni.New(gcpehttp.NewRecoverMiddleware(handler.log), gcpehttp.NewZapMiddleware(handler.log))
	handler.mux.Handle(pattern, root.With(negroni.Wrap(promhttp.Handler())))

	return handler
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}
