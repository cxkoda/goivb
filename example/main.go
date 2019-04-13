package main

import (
	"fmt"
	"github.com/cxkoda/goivb"
)

func main() {
	stops := goivb.GetStops()
	goivb.Printall(stops)
	fmt.Println()

	//smi := goivb.GetSmartinfo(1187) // Hauptbahnhof
	smi := goivb.GetSmartinfo(61549) // Höttinger Auffahrt
	goivb.Printall(smi)
	fmt.Println()
	goivb.FPrint(smi, 8)
}
