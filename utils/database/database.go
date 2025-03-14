package database

import (
	"billing_engine/config"
	"fmt"
	"log"

	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewDatabase(conf *config.Config) *sqlx.DB {
	db, err := sqlx.Connect("postgres", fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		conf.PostgresUser,
		conf.PostgresPassword,
		conf.PostgresHost,
		conf.PostgresPort,
		conf.PostgresDatabase))
	if err != nil {
		log.Fatalf("%v", err)
	}
	return db
}
