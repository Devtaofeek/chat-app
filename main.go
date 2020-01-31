package main

import (
	"chat_app/auth"
	"chat_app/chat"
	"flag"
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"text/template"
)

// templ is a single template
type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	t.templ.Execute(w, r)
}

func main() {
	var port = flag.String("port", ":8080", "The port where our application runs.")
	flag.Parse() // parse flags and extract appropriate information

	// new room setup
	r := chat.NewRoom()

	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.Handle("/chat", auth.Required(&templateHandler{filename: "chat.html"}))
	http.Handle("/room", r)

	// initialize the room
	go r.Run()

	// start the web server or log error
	log.Println("Starting server on", *port)
	if err := http.ListenAndServe(*port, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
