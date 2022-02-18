package main

import (
	"fmt"
	"net/http"

	"github.com/cloudprivacylabs/units"
)

func (app *application) normalizeUnit(w http.ResponseWriter, r *http.Request) {
	type input struct {
		Value string `json:"unit"`
		Hint  string `json:"hint"`
	}
	v, h, err := app.readFromParams(r)
	if err != nil {
		fmt.Println(err)
		return
	}
	value, hint, err := units.ParseUnits(v, h)
	if err != nil {
		fmt.Println(err)
		return
	}
	conversion := input{
		Value: value,
		Hint:  hint,
	}
	if err = app.writeJSON(w, http.StatusOK, conversion); err != nil {
		fmt.Println(err)
		return
	}
}
