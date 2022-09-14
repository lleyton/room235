package app

import (
	"context"
	"database/sql"
	"errors"

	"github.com/gocopper/copper/csql"
)

var ErrRecordNotFound = sql.ErrNoRows

func NewQueries(querier csql.Querier) *Queries {
	return &Queries{
		querier: querier,
	}
}

type Queries struct {
	querier csql.Querier
}

/*
Here are some example queries that use Querier to unmarshal results into Go strcuts

func (q *Queries) ListPosts(ctx context.Context) ([]Post, error) {
	const query = "SELECT * FROM posts ORDER BY created_at DESC"

	var (
	    posts []Post
	    err = q.querier.Select(ctx, &posts, query)
    )

	return posts, err
}

func (q *Queries) GetPostByID(ctx context.Context, id string) (*Post, error) {
	const query = "SELECT * from posts where id=?"

	var (
	    post Post
	    err = q.querier.Get(ctx, &post, query, id)
    )

	return &post, err
}

func (q *Queries) SavePost(ctx context.Context, post *Post) error {
	const query = `
	INSERT INTO posts (id, title, url, poster)
	VALUES (?, ?, ?, ?)
	ON CONFLICT (id) DO UPDATE SET title=?, url=?`

	_, err := q.querier.Exec(ctx, query,
		post.ID,
		post.Title,
		post.URL,
		post.Poster,
		post.Title,
		post.URL,
	)

	return err
}
*/

type Domain struct {
	Name string
	IP   string
}

func (q *Queries) ListDomains(ctx context.Context) ([]Domain, error) {
	const query = "SELECT * FROM domains"

	var (
		domains []Domain
		err     = q.querier.Select(ctx, &domains, query)
	)

	return domains, err
}

func (q *Queries) GetDomainByIP(ctx context.Context, ip string) (*Domain, error) {
	const query = "SELECT * FROM domains WHERE ip=?"

	var (
		domain Domain
		err    = q.querier.Get(ctx, &domain, query, ip)
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	return &domain, err
}

func (q *Queries) GetDomainByName(ctx context.Context, name string) (*Domain, error) {
	const query = "SELECT * FROM domains WHERE name=?"

	var (
		domain Domain
		err    = q.querier.Get(ctx, &domain, query, name)
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	return &domain, err
}

func (q *Queries) SaveDomain(ctx context.Context, domain *Domain) error {
	const query = `
	INSERT INTO domains (name, ip)
	VALUES (?, ?)`

	_, err := q.querier.Exec(ctx, query,
		domain.Name,
		domain.IP,
	)

	return err
}

func (q *Queries) DeleteDomain(ctx context.Context, name string) error {
	const query = "DELETE FROM domains WHERE name=?"

	_, err := q.querier.Exec(ctx, query, name)

	return err
}
