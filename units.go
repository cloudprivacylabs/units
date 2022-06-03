package units

import (
	"fmt"
	"regexp"
	"strings"
)

type ErrAmbiguousUnit []string

func (e ErrAmbiguousUnit) Error() string {
	return "Ambiguous unit: " + strings.Join([]string(e), ", ")
}

type unitRegex struct {
	regex     *regexp.Regexp
	parseFunc func(matches []string) (value, unit string, err error)
	unitName  string
}

var hintedRegex = make(map[string][]unitRegex, 0)

func RegisterUnixRegex(regex string, parseFunc func(matches []string) (value, unit string, err error),
	unitName string, hints ...string) {
	re := regexp.MustCompile(regex)
	ur := unitRegex{regex: re, parseFunc: parseFunc, unitName: unitName}
	for _, hint := range hints {
		hintedRegex[hint] = append(hintedRegex[hint], ur)
	}
	hintedRegex[""] = append(hintedRegex[""], ur)
}

var unitRegexp = regexp.MustCompile(`^([+\-]?(?:(?:0|[1-9]\d*)(?:\.\d*)?|\.\d+)(?:\d[eE][+\-]?\d+)?)(?:(?:\s+(\S+.*))|([^\seE\d]+.*))$`)

// parseMeasure parses a number and then a string for units
func parseMeasure(in string) (string, string, error) {
	values := make([]string, 0, 2)
	for _, v := range unitRegexp.FindAllStringSubmatch(in, -1) {
		for _, x := range v[1:] {
			x := strings.TrimSpace(x)
			if len(x) > 0 {
				values = append(values, x)
			}
		}
	}
	if len(values) != 2 {
		return "", "", fmt.Errorf("Input does not look like a value/unit pair: %s", in)
	}
	return values[0], values[1], nil
}

func ParseUnits(in string, hint ...string) (value, unit string, err error) {
	for _, ht := range hint {
		for _, unitRegex := range hintedRegex[ht] {
			matches := unitRegex.regex.FindAllStringSubmatch(in, -1)
			if len(matches) == 0 {
				continue
			}
			value, unit, err = unitRegex.parseFunc(matches[0])
			return value, unit, err
		}
	}
	var ambUnits []string
	if len(hint) == 0 {
		for _, unitRegex := range hintedRegex[""] {
			matches := unitRegex.regex.FindAllStringSubmatch(in, -1)
			if len(matches) == 0 {
				continue
			} else if len(matches) > 1 {
				ambUnits = append(ambUnits, unitRegex.unitName)
			} else {
				value, unit, err = unitRegex.parseFunc(matches[0])
				return value, unit, err
			}
		}
	}
	if len(ambUnits) > 1 {
		return "", "", ErrAmbiguousUnit(ambUnits)
	}
	value, unit, err = parseMeasure(in)
	if err == nil {
		return
	}
	return value, "", nil
}
