package app

import (
	"github.com/Yasuhiro-gh/jwt-app/internal/config"
	"github.com/Yasuhiro-gh/jwt-app/internal/db"
	"github.com/Yasuhiro-gh/jwt-app/internal/handlers"
	"github.com/Yasuhiro-gh/jwt-app/internal/usecase"
	"net/http"
)

func Run() {
	config.Run()
	pdb := db.NewPostgresDB()
	err := pdb.OpenConnection()
	if err != nil {
		panic(err)
	}
	defer func(pdb *db.PostgresDB) {
		err := pdb.CloseConnection()
		if err != nil {
			panic(err)
		}
	}(pdb)
	err = db.CreateTable(pdb)
	if err != nil {
		panic(err)
	}

	ts := usecase.NewTokenStorage(pdb)
	mux := handlers.Router(ts)

	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
