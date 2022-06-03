package main

import "net/http"

/*
API: GET /unit?value=1.2 m
Parse it, split value and number, normalize using regex, call js UCUM lib (edited)
GET /unit?value=1.2&unit=m
GET /unit?value=1.2&unit=m&targetUnit=cm

API: GET /validate  (UCUM passthrough)
API: GET /convert   (UCUM passthrough)
*/

func (app *application) routes() http.Handler {
	router := http.NewServeMux()
	router.HandleFunc("/unit", app.normalizeMeasure)
	router.HandleFunc("/validate", app.passthrough)
	router.HandleFunc("/convert", app.passthrough)
	return app.enableCORS(router)
}
