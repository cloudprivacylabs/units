package main

import (
	"flag"
	"fmt"
	"net/http"
)

type config struct {
	port int
}

type application struct {
	config config
}

func main() {
	run()
}

func run() {
	var cfg config
	flag.IntVar(&cfg.port, "port", 8080, "Server port to listen on")
	flag.Parse()
	app := &application{
		config: cfg,
	}
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.port),
		Handler: app.routes(),
	}
	fmt.Println("Starting server on port", cfg.port)
	err := srv.ListenAndServe()
	if err != nil {
		fmt.Println(err)

	}
}
