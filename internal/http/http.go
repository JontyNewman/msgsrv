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
	add func(string) (uint, error)
}

func NewAddMessageHandler(add func(string) (uint, error)) AddMessageHandler {
	return AddMessageHandler{add: add}
}

func (h *AddMessageHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	body, err := io.ReadAll(req.Body)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	id, err := h.add(string(body))

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

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
	fetch func(uint) (string, bool, error)
}

func NewFetchMessageHandler(fetch func(uint) (string, bool, error)) FetchMessageHandler {
	return FetchMessageHandler{fetch: fetch}
}

func (h *FetchMessageHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	id, err := strconv.Atoi(req.PathValue("id"))

	if err != nil {
		fmt.Fprintln(os.Stdout, err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	message, found, err := h.fetch(uint(id))

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	if !found {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	fmt.Fprintf(w, "%s", message)
}
