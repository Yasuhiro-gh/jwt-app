package config

import (
	"flag"
	"os"
)

var Options struct {
	DatabaseDSN string
}

func Run() {
	flag.StringVar(&Options.DatabaseDSN, "d", "postgres://postgres:1234@localhost:5432/postgres?sslmode=disable", "database-dsn")

	flag.Parse()

	if databaseDSN := os.Getenv("DATABASE_DSN"); databaseDSN != "" {
		Options.DatabaseDSN = databaseDSN
	}
}
