package main

import (
	"regexp"
	"strconv"
)

/*
A Go function in the UCUM project that checks if the
input string matches a given set of regexes, and if so, processes
the input string with that regex and sends it to UCUM. Something
like this:
func ParseUnits(in string) (value string, unit string) {
}
There should be a struct like this:
type UnitRegexes struct {
  Regex *regexp.Regexp
  Converter func(matches []string) (value string,unit string)
}
Define a function
func RegisterUnit(regex string, converter func( ) )
to register units.
ParseUnits should check every regex, find a match, and call the converter func
The first obvious case is the one for 5'4". The converter func should return 64 [in_i]
The idea is to get the value from input, first pass it through this function,
and then pass the output of this function to UCUM

The idea is to have this little Go library as a filtering layer in front of UCUM.

The RegisterUnit function is there to register filters.

We'll have some defaults, like the feet-inch thing
but if we need to add more exceptional processing, this allows us to "Register"
those exceptions by specifying a regexp and a processor function.

So the Register function will simply add the registered unit processor to a slice of structs. That's it.

The ParseUnits function will scan that list, and if any regex matches, call the processing function.

The Go-JS thing: On the Go side, we'll call the ParseUnits function, and then pass the returned value to UCUM service.

This is the "registry" pattern, pretty common for things like: database drivers, the "terms" stuff we
have in the LSA library, text encoding, etc.

For instance: we have a NewTerm function that registers a term. It adds that term to a global map.
Then there is a function called GetTermMetadata() that looks up the term in that map and returns its metadata.
*/

type UnitRegex struct {
	Regex     *regexp.Regexp
	Converter func(matches []string) (value, unit string)
}

var HintedRegex map[string][]UnitRegex

// var acceptedRegex []UnitRegex

func init() {
	RegisterUnixRegex(`(?P<ft>[[:digit:]]+)\'(?P<in>[[:digit:]]+)\"`, ParseFeetInch, "[in_i]", "height", "length")
}

func RegisterUnixRegex(regex string, converter func(matches []string) (value, unit string),
	unitName string, hints ...string) {
}

// compile regex
// create instance of unitregex and append to slice
//

func ParseUnits(in string, hint ...string) (value, unit string) {
	for ht, hintRegex := range HintedRegex {
		for idx, rx := range hintRegex {
			// if there is no hint, potentially check every struct in the slice for a match
			if len(hint) == 0 || hint[0] == "" {
				matches := rx.Regex.FindAllStringSubmatch(in, -1)
				if len(matches) == 0 {
					continue
				}
				return rx.Converter(matches[0][1:])
			} else {
				// otherwise, if the provided hint matches a hint in the map
				if hint[idx] == ht {
					// list all potential matching regexes
					matchingRegex := make([]string, 0, len(hintRegex))
					matchingRegex = append(matchingRegex)
				}
			}
		}
	}
	// for _, regex := range acceptedRegex {
	// 	rx := regex.Regex
	// 	matches := rx.FindAllStringSubmatch(in, -1)
	// 	if len(matches) == 0 {
	// 		continue
	// 	}
	// 	return regex.Converter(matches[0][1:])
	// }
	// //fmt.Println(value, unit)
	// return value, ""
}

// func RegisterUnit(regex string, converter func()) {
// 	re := regexp.MustCompile(regex)
// 	u := UnitRegexes{Regex: re}
// 	AcceptedRegex = append(AcceptedRegex, u)
// }

func main() {

	ParseUnits(`5'4"`)
	//fmt.Println(AcceptedRegex)
}

// re := regexp.MustCompile(`(?P<ft>[[:digit:]]+)\'(?P<in>[[:digit:]]+)\"`)
// r2 := re.FindAllStringSubmatch(`10'5"`, -1)[0]

func ParseFeetInch(matches []string) (value, unit string) {
	var sum int
	const FOOT = 12
	ft, _ := strconv.Atoi(matches[1])
	in, _ := strconv.Atoi(matches[2])
	sum += (ft * FOOT) + in
	return strconv.Itoa(sum), "[in_i]"

	//value, unit = regex.Converter(matches)
}

// `(?P<ft>[[:digit:]]+)\'(?P<in>[[:digit:]]+)\"`:
