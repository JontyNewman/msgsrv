package http

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/jontynewman/msgsrv/internal/repo"
)

type AddMessageHandler struct {
	Repo repo.MessageRepository
}

func (h *AddMessageHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	body, err := io.ReadAll(req.Body)

	if err != nil {
		fmt.Fprintln(log.Writer(), err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	id := h.Repo.Add(string(body))

	json, err := json.Marshal(struct {
		Id uint `json:"id"`
	}{Id: id})

	if err != nil {
		fmt.Fprintln(log.Writer(), err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")
	fmt.Fprintf(w, "%s", json)
}

type FetchMessageHandler struct {
	Repo repo.MessageRepository
}

func (h *FetchMessageHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	id, err := strconv.Atoi(req.PathValue("id"))

	if err != nil {
		fmt.Fprintln(log.Writer(), err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	message, found := h.Repo.Fetch(uint(id))

	if !found {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	fmt.Fprintf(w, "%s", message)
}
