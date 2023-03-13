package repositories

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Config struct {
	DatabaseHost     string `envconfig:"database_host"`
	DatabasePort     int    `envconfig:"database_port"`
	DatabaseName     string `envconfig:"database_name"`
	DatabaseUsername string `envconfig:"database_username"`
	DatabasePassword string `envconfig:"database_password"`
}

func (c Config) DBString() string {
	return fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%d sslmode=disable",
		c.DatabaseUsername,
		c.DatabasePassword,
		c.DatabaseName,
		c.DatabaseHost,
		c.DatabasePort,
	)
}

func InitDatabase(cfg Config) *sqlx.DB {
	return sqlx.MustConnect("postgres", cfg.DBString())
}
