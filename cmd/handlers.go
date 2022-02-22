package main

import (
	"fmt"
	"net/http"

	"github.com/cloudprivacylabs/units"
)

func (app *application) normalizeMeasure(w http.ResponseWriter, r *http.Request) {
	type input struct {
		Value string `json:"value"`
		Unit  string `json:"unit"`
	}
	v, hint, err := app.readFromParams(r)
	if err != nil {
		fmt.Println(err)
		return
	}
	var value, unit string
	if hint != "" {
		value, unit, err = units.ParseUnits(v, hint)
	} else {
		value, unit, err = units.ParseUnits(v)
	}
	if err != nil {
		fmt.Println(err)
		return
	}
	conversion := input{
		Value: value,
		Unit:  unit,
	}
	if err = app.writeJSON(w, http.StatusOK, conversion); err != nil {
		fmt.Println(err)
		return
	}
}
