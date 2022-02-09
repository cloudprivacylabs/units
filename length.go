package units

import (
	"strconv"
)

func init() {
	RegisterUnixRegex(`(?P<ft>[[:digit:]]+)\'(?P<in>[[:digit:]]+)\"`, LengthFeetInch, "[in_i]", "height", "length")
}

func LengthFeetInch(matches []string) (value, unit string, err error) {
	var sum int
	const FOOT = 12
	ft, err := strconv.Atoi(matches[1])
	if err != nil {
		return "", "", err
	}
	in, err := strconv.Atoi(matches[2])
	if err != nil {
		return "", "", err
	}
	sum += (ft * FOOT) + in
	return strconv.Itoa(sum), "[in_i]", nil
}
