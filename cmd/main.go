package main

import (
	"flag"
	"fmt"
	"net/http"
)

type config struct {
}

type application struct {
	listen  string
	ucumURL string
}

func main() {
	run()
}

func run() {
	var app application
	flag.StringVar(&app.listen, "listen", ":8081", "Server address to listen on")
	flag.StringVar(&app.ucumURL, "ucum", "http://localhost:8080", "UCUM server address")
	flag.Parse()
	srv := &http.Server{
		Addr:    app.listen,
		Handler: app.routes(),
	}
	fmt.Println("Starting server on ", app.listen)
	srv.ListenAndServe()
}
