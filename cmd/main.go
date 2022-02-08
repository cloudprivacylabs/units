package main

import (
	"fmt"

	"github.com/cloudprivacylabs/units"
)

func main() {
	x, y := units.ParseUnits(`5'4"`)
	x, y = units.ParseUnits(`5'4"`, "length")
	x, y = units.ParseUnits(`5'4"`, "height")
	fmt.Println(x, y)
	fmt.Println(units.HintedRegex)
}
