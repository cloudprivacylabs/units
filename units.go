package units

import (
	"regexp"
	"strconv"
	"strings"
)

type ErrAmbiguousUnit []string

func (e ErrAmbiguousUnit) Error() string {
	return "Ambiguous unit: " + strings.Join([]string(e), ", ")
}

type UnitRegex struct {
	Regex     *regexp.Regexp
	Converter func(matches []string) (value, unit string)
	UnitName  string
}

var hintedRegex = make(map[string][]UnitRegex, 0)

func init() {
	RegisterUnixRegex(`(?P<ft>[[:digit:]]+)\'(?P<in>[[:digit:]]+)\"`, ParseFeetInch, "[in_i]", "height", "length")
}

func RegisterUnixRegex(regex string, converter func(matches []string) (value, unit string),
	unitName string, hints ...string) {
	re := regexp.MustCompile(regex)
	ur := UnitRegex{Regex: re, Converter: converter, UnitName: unitName}
	for _, hint := range hints {
		hintedRegex[hint] = append(hintedRegex[hint], ur)
	}
	hintedRegex[""] = append(hintedRegex[""], ur)
}

func ParseUnits(in string, hint ...string) (value, unit string, err error) {
	for _, ht := range hint {
		for _, unitRegex := range hintedRegex[ht] {
			matches := unitRegex.Regex.FindAllStringSubmatch(in, -1)
			if len(matches) == 0 {
				continue
			}
			value, unit = unitRegex.Converter(matches[0])
			return value, unit, nil
		}
	}
	var ambUnits ErrAmbiguousUnit
	if len(hint) == 0 {
		for _, unitRegex := range hintedRegex[""] {
			matches := unitRegex.Regex.FindAllStringSubmatch(in, -1)
			if len(matches) == 0 {
				continue
			} else if len(matches) > 1 {
				switch {
				case unitRegex.Regex.String() == `(?P<ft>[[:digit:]]+)\'(?P<in>[[:digit:]]+)\"`:
					ambUnits = append(ambUnits, unitRegex.UnitName)
				}
			} else {
				value, unit = unitRegex.Converter(matches[0])
				return value, unit, nil
			}
		}
	}
	return value, "", ambUnits
}

func ParseFeetInch(matches []string) (value, unit string) {
	var sum int
	const FOOT = 12
	ft, _ := strconv.Atoi(matches[1])
	in, _ := strconv.Atoi(matches[2])
	sum += (ft * FOOT) + in
	return strconv.Itoa(sum), "[in_i]"
}
