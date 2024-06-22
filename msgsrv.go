package msgsrv

import (
	"net/http"

	h "internal/http"
	"internal/repo"
)

func Run(addr string, repo repo.MessageRepository) error {

	add := h.NewAddMessageHandler(repo.Add)
	fetch := h.NewFetchMessageHandler(repo.Fetch)

	http.Handle("POST /messages/{$}", &add)
	http.Handle("GET /messages/{id}", &fetch)

	return http.ListenAndServe(addr, nil)
}
