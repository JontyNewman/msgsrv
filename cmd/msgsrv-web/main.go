package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"

	h "internal/http"
	r "internal/repo/http"
)

type Message struct {
	ID   uint
	Body string
}

func main() {

	addr := flag.String("addr", ":8080", "The TCP network address for this service to listen on.")
	baseAddress := flag.String("repo", "http://localhost:80/messages/", "The base address of the repository.")
	flag.Parse()

	repo := r.InitHttpMessageRepository(*baseAddress)
	tmpl := template.Must(template.ParseFiles("form.html"))

	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, nil)
	})

	fetch := h.NewFetchMessageHandler(repo.Fetch)
	http.Handle("GET /{id}", &fetch)

	http.HandleFunc("POST /", func(w http.ResponseWriter, r *http.Request) {

		message := r.FormValue("message")

		id, err := repo.Add(message)

		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tmpl.Execute(w, Message{ID: id, Body: message})
	})

	http.ListenAndServe(*addr, nil)
}
