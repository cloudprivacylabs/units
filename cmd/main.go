package main

import (
	"flag"
	"fmt"

	"github.com/cloudprivacylabs/units"
)

func main() {
	var cfg config
	flag.StringVar(&cfg.unit, "unit", "", "define unit")
	flag.StringVar(&cfg.hint, "hint", "", "define hint")
	flag.Parse()
	x, y, z := units.ParseUnits(cfg.unit)
	fmt.Println(x, y, z)
}

type config struct {
	unit string
	hint string
}
