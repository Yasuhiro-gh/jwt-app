package app

import (
	"github.com/Yasuhiro-gh/jwt-app/internal/handlers"
	"net/http"
)

func Run() {
	mux := handlers.Router()

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
