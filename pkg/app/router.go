package app

import (
	"database/sql"
	"net"
	"net/http"

	"github.com/gocopper/copper/chttp"
	"github.com/gocopper/copper/clogger"
	"github.com/miekg/dns"
)

type NewRouterParams struct {
	RW      *chttp.ReaderWriter
	Logger  clogger.Logger
	Queries *Queries
	Db      *sql.DB
}

func NewRouter(p NewRouterParams) *Router {
	resolver := Resolver{
		queries: p.Queries,
		db:      p.Db,
		client:  &dns.Client{},
	}

	go func() {
		if err := resolver.Serve(); err != nil {
			panic(err)
		}
	}()

	return &Router{
		queries: p.Queries,
		rw:      p.RW,
		logger:  p.Logger,
	}
}

type Router struct {
	rw      *chttp.ReaderWriter
	logger  clogger.Logger
	queries *Queries
}

func (ro *Router) Routes() []chttp.Route {
	return []chttp.Route{
		{
			Path:    "/remove-domain",
			Methods: []string{http.MethodPost},
			Handler: ro.HandleRemoveDomainForm,
		},

		{
			Path:    "/remove-domain",
			Methods: []string{http.MethodGet},
			Handler: ro.HandleRemoveDomainPage,
		},

		{
			Path:    "/register-domain",
			Methods: []string{http.MethodPost},
			Handler: ro.HandleRegisterDomainForm,
		},

		{
			Path:    "/register-domain",
			Methods: []string{http.MethodGet},
			Handler: ro.HandleRegisterDomainsPage,
		},

		{
			Path:    "/",
			Methods: []string{http.MethodGet},
			Handler: ro.HandleIndexPage,
		},
	}
}

func (ro *Router) HandleIndexPage(w http.ResponseWriter, r *http.Request) {
	host, _, _ := net.SplitHostPort(r.RemoteAddr)

	domains, err := ro.queries.ListDomains(r.Context())
	if err != nil {
		panic(err)
	}

	domain, err := ro.queries.GetDomainByIP(r.Context(), host)
	if err != nil {
		panic(err)
	}

	data := map[string]interface{}{
		"IP":    host,
		"Hosts": domains,
	}

	if domain != nil {
		data["Domain"] = domain.Name
	}

	// It's safe to ignore the error here because RemoteAddr is always a correct value.
	ro.rw.WriteHTML(w, r, chttp.WriteHTMLParams{
		PageTemplate: "index.html",
		Data:         data,
	})
}

func (ro *Router) HandleRegisterDomainsPage(w http.ResponseWriter, r *http.Request) {
	host, _, _ := net.SplitHostPort(r.RemoteAddr)

	ro.rw.WriteHTML(w, r, chttp.WriteHTMLParams{
		PageTemplate: "register-domain.html",
		Data: map[string]interface{}{
			"IP": host},
	})
}

func (ro *Router) HandleRegisterDomainForm(w http.ResponseWriter, r *http.Request) {
	domainName := r.FormValue("domain")
	host, _, _ := net.SplitHostPort(r.RemoteAddr)

	domain, err := ro.queries.GetDomainByName(r.Context(), domainName)
	if err != nil {
		panic(err)
	}

	if domain != nil {
		println(domain.IP)
		ro.logger.Info("domain already registered " + domainName)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	domain, err = ro.queries.GetDomainByIP(r.Context(), host)
	if err != nil {
		panic(err)
	}

	if domain != nil {
		if err := ro.queries.DeleteDomain(r.Context(), domain.Name); err != nil {
			panic(err)
		}
	}

	if err := ro.queries.SaveDomain(r.Context(), &Domain{
		domainName,
		host,
	}); err != nil {
		panic(err)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (ro *Router) HandleRemoveDomainPage(w http.ResponseWriter, r *http.Request) {
	host, _, _ := net.SplitHostPort(r.RemoteAddr)

	domain, err := ro.queries.GetDomainByIP(r.Context(), host)
	if err != nil {
		panic(err)
	}

	if domain == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	ro.rw.WriteHTML(w, r, chttp.WriteHTMLParams{
		PageTemplate: "remove-domain.html",
		Data: map[string]interface{}{
			"Name": domain.Name,
		},
	})
}

func (ro *Router) HandleRemoveDomainForm(w http.ResponseWriter, r *http.Request) {
	host, _, _ := net.SplitHostPort(r.RemoteAddr)

	domain, err := ro.queries.GetDomainByIP(r.Context(), host)
	if err != nil {
		panic(err)
	}

	if domain == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if err := ro.queries.DeleteDomain(r.Context(), domain.Name); err != nil {
		panic(err)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
