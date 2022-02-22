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
	// var hint string
	// var value string
	// flag.StringVar(&hint, "hint", "", "Unit lookup hint, length, height, etc.")
	// flag.StringVar(&value, "value", "", "value")
	// flag.Parse()
	// if len(flag.Args()) > 2 {
	// 	errors.New("too many args")
	// }
	// fmt.Println(value, hint)
	// x, y, _ := units.ParseUnits(value, hint)
	// fmt.Println(x, y)
	run()
}

func run() {
	var cfg config
	flag.IntVar(&cfg.port, "port", 8081, "Server port to listen on")
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
