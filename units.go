package units

import (
	"fmt"
	"regexp"
	"strconv"
)

type UnitRegex struct {
	Regex     *regexp.Regexp
	Converter func(matches []string) (value, unit string)
}

var HintedRegex = make(map[string][]UnitRegex, 0)

func init() {
	RegisterUnixRegex(`(?P<ft>[[:digit:]]+)\'(?P<in>[[:digit:]]+)\"`, ParseFeetInch, "[in_i]", "height", "length")
}

func RegisterUnixRegex(regex string, converter func(matches []string) (value, unit string),
	unitName string, hints ...string) {
	re := regexp.MustCompile(regex)
	ur := UnitRegex{Regex: re}
	for _, hint := range hints {
		HintedRegex[hint] = append(HintedRegex[hint], ur)
		HintedRegex[""] = append(HintedRegex[""], ur)
	}
}

func ParseUnits(in string, hint ...string) (value, unit string) {
	for _, ht := range hint {
		for _, unitRegex := range HintedRegex[ht] {
			matches := unitRegex.Regex.FindAllStringSubmatch(in, -1)
			if len(matches) == 0 {
				continue
			}
			switch {
			case ht == "length":
				unitRegex.Converter = ParseFeetInch
			case ht == "height":
				unitRegex.Converter = ParseFeetInch
			}
			return unitRegex.Converter(matches[0][:])
		}
	}
	if len(hint) == 0 {
		for _, unitRegex := range HintedRegex[""] {
			matches := unitRegex.Regex.FindAllStringSubmatch(in, -1)
			if len(matches) == 0 {
				continue
			}
			fmt.Println("unit:", matches[0][0], "-> AMBIGUOUS")
		}
	}
	return value, ""
}

func ParseFeetInch(matches []string) (value, unit string) {
	var sum int
	const FOOT = 12
	ft, _ := strconv.Atoi(matches[1])
	in, _ := strconv.Atoi(matches[2])
	sum += (ft * FOOT) + in
	return strconv.Itoa(sum), "[in_i]"
}
