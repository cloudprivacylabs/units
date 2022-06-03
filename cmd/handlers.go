package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	pluralize "github.com/gertd/go-pluralize"

	"github.com/cloudprivacylabs/units"
)

func (app *application) validateUnit(unit string) (string, error) {
	q := url.Values{
		"unit": []string{unit},
	}
	fmt.Println("Validate", q)
	// call UCUM validate
	resp, err := http.Get(app.ucumURL + "/validate?" + q.Encode())
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("Cannot validate unit %s: %s", unit, string(body))
	}
	var valid struct {
		Status   string `json:"status"`
		UcumCode string `json:"ucumCode"`
	}
	err = json.Unmarshal(body, &valid)
	if err != nil {
		return "", err
	}
	if len(valid.UcumCode) > 0 {
		return valid.UcumCode, nil
	}
	return "", nil
}

func (app *application) normalizeMeasure(w http.ResponseWriter, r *http.Request) {
	var err error
	query := r.URL.Query()
	v := query.Get("value")
	hint := query.Get("hint") // optional
	var value, unit string
	if hint != "" {
		value, unit, err = units.ParseUnits(v, hint)
	} else {
		value, unit, err = units.ParseUnits(v)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(unit) > 0 {
		plu := pluralize.NewClient()
		u, err := app.validateUnit(unit)
		if err != nil || len(u) == 0 {
			fmt.Println("Checking plural/singular", unit)
			if plu.IsPlural(unit) {
				u, err = app.validateUnit(plu.Singular(unit))
			} else {
				u, err = app.validateUnit(plu.Plural(unit))
			}
			if err == nil && len(u) > 0 {
				unit = u
			}
		} else {
			unit = u
		}
	}

	app.writeJSON(w, http.StatusOK, map[string]interface{}{
		"value": value, "unit": unit})
}

func (app *application) passthrough(w http.ResponseWriter, r *http.Request) {
	rsp, err := http.Get(app.ucumURL + r.URL.Query().Encode())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rsp.Body.Close()
	if (rsp.StatusCode % 100) != 2 {
		data, _ := ioutil.ReadAll(rsp.Body)
		http.Error(w, string(data), rsp.StatusCode)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	io.Copy(w, rsp.Body)
}
