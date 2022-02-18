package main

import "net/http"

/*
API: GET /unit?value=1.2 m
Parse it, split value and number, normalize using regex, call js UCUM lib (edited)
GET /unit?value=1.2&unit=m
GET /unit?value=1.2&unit=m&targetUnit=cm
*/

func (app *application) routes() http.Handler {
	router := http.NewServeMux()
	router.HandleFunc("/unit", app.normalizeUnit)
	return app.enableCORS(router)
}
