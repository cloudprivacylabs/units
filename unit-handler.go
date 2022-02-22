package units

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type UnitHandler struct {
	UCUMURL string
}

func (hnd UnitHandler) Normalize(input string, hint ...string) (string, string, error) {
	value, unit, err := ParseUnits(input, hint...)
	if err != nil {
		return "", "", err
	}
	return value, unit, nil
}

// unit=meter&unit=m&value=value
func (hnd UnitHandler) NormalizeUnit(fromUnit string) (string, error) {
	query := url.Values{
		"unit": []string{fromUnit},
	}
	// call UCUM validate
	resp, err := http.Get(hnd.UCUMURL + "/validate?" + query.Encode())
	if err != nil {
		return err.Error(), err
	}
	var valid struct {
		Status   string `json:"status"`
		UcumCode string `json:"ucumCode"`
	}
	err = json.NewDecoder(resp.Body).Decode(&valid)
	if err != nil {
		return err.Error(), err
	}
	if valid.Status == "invalid" && valid.UcumCode == "" {
		return errors.New("invalid input").Error(), err
	}
	return valid.UcumCode, nil
}

func (hnd UnitHandler) Convert(fromUnit, targetUnit string, value int) (string, error) {
	query := url.Values{
		"unit":   []string{fromUnit},
		"value":  []string{strconv.Itoa(value)},
		"output": []string{targetUnit},
	}
	// call UCUM convert
	resp, err := http.Get(hnd.UCUMURL + "/convert?" + query.Encode())
	if err != nil {
		return err.Error(), nil
	}
	var conversion struct {
		Status  string   `json:"status"`
		ToVal   float64  `json:"toVal"`
		Message []string `json:"msg"`
	}
	err = json.NewDecoder(resp.Body).Decode(&conversion)
	if err != nil {
		return err.Error(), err
	}
	if conversion.Status == "failed" && conversion.ToVal == 0 {
		return strings.Join(conversion.Message, ","), errors.New("invalid UCUM codes")
	}
	js, err := json.MarshalIndent(conversion, "", "\t")
	if err != nil {
		return err.Error(), err
	}
	return string(js), nil
}
