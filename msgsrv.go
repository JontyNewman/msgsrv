package main

import (
	"log"
	"net/http"

	h "github.com/jontynewman/msgsrv/internal/http"
	"github.com/jontynewman/msgsrv/internal/repo"
)

func main() {

	repo := new(repo.RuntimeMessageRepository)

	add := h.NewAddMessageHandler(repo.Add)
	fetch := h.NewFetchMessageHandler(repo.Fetch)

	http.Handle("POST /messages/{$}", &add)
	http.Handle("GET /messages/{id}", &fetch)

	log.Fatal(http.ListenAndServe("", nil))
}
