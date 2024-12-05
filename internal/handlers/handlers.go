package handlers

import "net/http"

func Router() *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("/api/get_pair", GetPair())
	mux.Handle("/api/refresh", Refresh())

	return mux
}

func GetPair() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("Get pair endpoint"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}

func Refresh() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("Refresh endpoint"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}
