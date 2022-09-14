package app

import (
	"github.com/gocopper/copper/chttp"
	"github.com/gocopper/copper/clogger"
	"net/http"
)

type NewRouterParams struct {
	RW     *chttp.ReaderWriter
	Logger clogger.Logger
}

func NewRouter(p NewRouterParams) *Router {
	return &Router{
		rw:     p.RW,
		logger: p.Logger,
	}
}

type Router struct {
	rw     *chttp.ReaderWriter
	logger clogger.Logger
}

func (ro *Router) Routes() []chttp.Route {
	return []chttp.Route{
		{
			Path:    "/",
			Methods: []string{http.MethodGet},
			Handler: ro.HandleIndexPage,
		},
	}
}

func (ro *Router) HandleIndexPage(w http.ResponseWriter, r *http.Request) {

	ro.rw.WriteHTML(w, r, chttp.WriteHTMLParams{
		PageTemplate: "index.html",
	})
}
