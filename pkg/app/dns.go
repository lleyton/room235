package app

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/gocopper/copper/csql"
	"github.com/miekg/dns"
)

type Resolver struct {
	queries *Queries
	db      *sql.DB
	client  *dns.Client
}

func (resolver *Resolver) resolveUpstream(w dns.ResponseWriter, r *dns.Msg) {
	println("Resolving upstream")
	ctx := context.TODO()

	r, _, err := resolver.client.ExchangeContext(ctx, r, "1.1.1.1:53")
	if err != nil {
		panic(err)
	}

	w.WriteMsg(r)
}

func (resolver *Resolver) parseQuery(ctx context.Context, m *dns.Msg) {
	// Yes, I know this is stupid
	for _, q := range m.Question {
		switch q.Qtype {
		case dns.TypeA:
			log.Printf("Query for %s\n", q.Name)
			domain, err := resolver.queries.GetDomainByName(ctx, strings.TrimSuffix(q.Name, ".lab."))
			if err != nil {
				continue
			}

			if domain == nil {
				continue
			}

			ip := domain.IP

			if ip != "" {
				rr, err := dns.NewRR(fmt.Sprintf("%s A %s", q.Name, ip))
				if err == nil {
					m.Answer = append(m.Answer, rr)
				}
			}
		}
	}
}

func (resolver *Resolver) handleDnsRequest(w dns.ResponseWriter, r *dns.Msg) {
	ctx, _, err := csql.CtxWithTx(context.TODO(), resolver.db, "sqlite3")
	tx, _ := csql.TxFromCtx(ctx)
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	defer tx.Commit()

	m := new(dns.Msg)
	m.SetReply(r)
	m.Compress = false

	switch r.Opcode {
	case dns.OpcodeQuery:
		resolver.parseQuery(ctx, m)
	}

	w.WriteMsg(m)
}

func (resolver *Resolver) Serve() error {
	dns.HandleFunc("lan.", resolver.handleDnsRequest)
	dns.HandleFunc(".", resolver.resolveUpstream)

	port := 5354
	server := &dns.Server{Addr: ":" + strconv.Itoa(port), Net: "udp"}
	log.Printf("Starting at %d\n", port)
	err := server.ListenAndServe()
	defer server.Shutdown()
	if err != nil {
		log.Fatalf("Failed to start server: %s\n ", err.Error())
	}

	return nil
}
