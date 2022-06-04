package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	pluralize "github.com/gertd/go-pluralize"

	"github.com/cloudprivacylabs/units"
)

type ucumValidationResponse struct {
	Status string `json:"status"`
	Code   string `json:"ucumCode"`
	Unit   struct {
		Code string `json:"code"`
		Name string `json:"name"`
	} `json:"unit"`
	Msg []string `json:"msg"`
}

func (app *application) validateUnit(unit string) (ucumValidationResponse, error) {
	q := url.Values{
		"unit": []string{unit},
	}
	// call UCUM validate
	resp, err := http.Get(app.ucumURL + "/validate?" + q.Encode())
	if err != nil {
		return ucumValidationResponse{}, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return ucumValidationResponse{}, fmt.Errorf("Cannot validate unit %s: %s", unit, string(body))
	}
	var valid ucumValidationResponse
	err = json.Unmarshal(body, &valid)
	if err != nil {
		return valid, err
	}
	return valid, nil
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
	ret := map[string]interface{}{
		"value": value,
		"unit":  unit,
		"valid": false,
	}

	if len(unit) > 0 {
		u, err := app.validateUnit(unit)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if u.Status == "valid" {
			ret["unit"] = u.Code
			ret["valid"] = true
		} else if u.Status == "invalid" && len(u.Code) > 0 {
			ret["unit"] = u.Code
			ret["valid"] = false
			ret["msg"] = strings.Join(u.Msg, "/")
		} else {
			plu := pluralize.NewClient()
			var pluralizedRsp ucumValidationResponse
			if plu.IsPlural(unit) {
				pluralizedRsp, err = app.validateUnit(plu.Singular(unit))
			} else {
				pluralizedRsp, err = app.validateUnit(plu.Plural(unit))
			}
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if pluralizedRsp.Status == "valid" {
				ret["unit"] = pluralizedRsp.Code
				ret["valid"] = false
				ret["msg"] = "Pluralization difference"
			} else if pluralizedRsp.Status == "invalid" && len(pluralizedRsp.Code) > 0 {
				ret["unit"] = pluralizedRsp.Code
				ret["valid"] = false
				ret["msg"] = strings.Join(u.Msg, "/")
			} else {
				ret["msg"] = strings.Join(u.Msg, "/")
			}
		}
	}

	app.writeJSON(w, http.StatusOK, ret)
}

func (app *application) passthrough(w http.ResponseWriter, r *http.Request) {
	rsp, err := http.Get(app.ucumURL + r.URL.Path + "?" + r.URL.Query().Encode())
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
