package main

import (
	"encoding/json"
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

func (app *application) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	enc := json.NewEncoder(w)
	enc.Encode(data)
}
