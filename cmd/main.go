package main

import (
	"flag"
	"fmt"

	"github.com/cloudprivacylabs/units"
)

func main() {
	var hint string
	flag.StringVar(&hint, "hint", "", "Unit lookup hint, length, height, etc.")
	flag.Parse()
	if len(flag.Args()) != 1 {
		fmt.Println("Need one arg")
		return
	}
	x, y, err := units.ParseUnits(flag.Args()[0], hint)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(x, y)
}
