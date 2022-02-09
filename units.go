package units

import (
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
	return value, "", nil
}
