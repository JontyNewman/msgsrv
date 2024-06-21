package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

type AddMessageHandler struct {
	add func(string) uint
}

func NewAddMessageHandler(add func(string) uint) AddMessageHandler {
	return AddMessageHandler{add: add}
}

func (h *AddMessageHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	body, err := io.ReadAll(req.Body)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	id := h.add(string(body))

	json, err := json.Marshal(struct {
		Id uint `json:"id"`
	}{Id: id})

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")
	fmt.Fprintf(w, "%s", json)
}

type FetchMessageHandler struct {
	fetch func(uint) (string, bool)
}

func NewFetchMessageHandler(fetch func(uint) (string, bool)) FetchMessageHandler {
	return FetchMessageHandler{fetch: fetch}
}

func (h *FetchMessageHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	id, err := strconv.Atoi(req.PathValue("id"))

	if err != nil {
		fmt.Fprintln(os.Stdout, err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	message, found := h.fetch(uint(id))

	if !found {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	fmt.Fprintf(w, "%s", message)
}
