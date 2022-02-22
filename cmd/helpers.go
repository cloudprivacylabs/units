package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// func (app *application) readFromFlags() (string, string, error) {
// 	var hint string
// 	flag.StringVar(&hint, "hint", "", "Unit lookup hint, length, height, etc.")
// 	flag.Parse()
// 	if len(flag.Args()) != 1 {
// 		return "", "", errors.New("need one arg")
// 	}
// 	x, y, err := units.ParseUnits(flag.Args()[0], hint)
// 	if err != nil {
// 		return "", "", err
// 	}
// 	return x, y, nil
// }

// curl "localhost:8080/unit?value=5'"'4"'"&hint=length"
func (app *application) readFromParams(r *http.Request) (string, string, error) {
	query := r.URL.Query()
	value := query.Get("value")
	hint := query.Get("hint") // optional
	fmt.Println("user input: "+value, hint)
	// fmt.Println(units.Convert("gg", "kg", 100))
	return value, hint, nil
}

func (app *application) writeJSON(w http.ResponseWriter, status int, data interface{}) error {
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}
	js = append(js, '\n')

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)
	return nil
}
