package context

import (
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/lstrgiang/cryptowatch-server/internal/infra/cache"
)

type (
	Context interface {
		GetDB() *sqlx.DB
		GetCache() cache.Cache
		GetAuthSecretKey() string
		GetHTTPClient() *http.Client
		GetDomain() string
		GetPath() string
	}
	context struct {
		db            *sqlx.DB
		cache         cache.Cache
		authSecretKey string
		domain        string
		path          string
	}
)

func NewContext(db *sqlx.DB, cache cache.Cache, authSecretKey string, domain string) Context {
	return context{
		db:            db,
		authSecretKey: authSecretKey,
		cache:         cache,
		domain:        domain,
		path:          "/", // should get from config also
	}
}

func (c context) GetDB() *sqlx.DB {
	return c.db
}

func (c context) GetCache() cache.Cache {
	return c.cache
}

func (c context) GetAuthSecretKey() string {
	return c.authSecretKey
}

func (c context) GetHTTPClient() *http.Client {
	return &http.Client{}
}

func (c context) GetDomain() string {
	return c.domain
}

func (c context) GetPath() string {
	return c.path
}
