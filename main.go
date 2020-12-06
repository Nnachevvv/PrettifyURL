package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/nnachevv/PretifyURL/server"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/encode", server.Encode)
	r.HandleFunc("/{short}", server.Short)

	http.Handle("/", r)
	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
