package units

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

type UnitHandler struct {
	UCUMURL string
}

func Normalize(input string) (string, string, error) { // 1.2m
	value, unit, err := ParseUnits(input)
	if err != nil {
		return "", "", err
	}
	if value != "" && unit != "" {
		return value, unit, nil
	}
	re := regexp.MustCompile(`[-]?\d[\d,]*[\.]?[\d{2}]*`)
	// rel := regexp.MustCompile(`[a-zA-Z]+`)
	// letter_split := re.Split(input, -1) // may contain spaces
	//letters := strings.TrimSpace(strings.Join(rel.Split(input, -1), ""))
	letter_split := strings.TrimSpace(strings.Join(re.Split(input, -1), ""))
	num_split := strings.Join(re.FindAllString(input, -1), " ")
	return num_split, letter_split, nil
}

// unit=meter&unit=m&value=value
func (hnd UnitHandler) NormalizeUnit(fromUnit string) ([]string, error) {
	query := url.Values{
		"unit": []string{fromUnit},
	}
	// call UCUM validate
	resp, err := http.Get(hnd.UCUMURL + "/validate?" + query.Encode())
	if resp.StatusCode != 200 {
		panic("internal server error")
	}
	if err != nil {
		return []string{err.Error()}, err
	}
	var valid struct {
		Status   string   `json:"status"`
		UcumCode []string `json:"ucumCode"`
	}
	err = json.NewDecoder(resp.Body).Decode(&valid)
	if err != nil {
		return []string{err.Error()}, err
	}
	if valid.Status == "invalid" && len(valid.UcumCode) == 0 {
		return []string{errors.New("invalid input").Error()}, err
	}
	return valid.UcumCode, nil
}

func (hnd UnitHandler) Convert(fromUnit, targetUnit, value string) (string, error) {
	query := url.Values{
		"unit":   []string{fromUnit},
		"value":  []string{value},
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
