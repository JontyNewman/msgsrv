package main

import (
	"log"
	"net/http"

	h "github.com/jontynewman/msgsrv/internal/http"
	"github.com/jontynewman/msgsrv/internal/repo"
)

func main() {

	repo := repo.RuntimeMessageRepository{}

	http.Handle("POST /messages/{$}", &h.AddMessageHandler{Repo: &repo})
	http.Handle("GET /messages/{id}", &h.FetchMessageHandler{Repo: &repo})

	log.Fatal(http.ListenAndServe("", nil))
}
