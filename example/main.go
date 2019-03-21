package main

import (
	"fmt"
	"github.com/cxkoda/goivb"
)

func main() {

	smi := goivb.GetData(1187) // Hauptbahnhof
	//smi := goivb.GetData(61549) // Höttinger Auffahrt
	goivb.Printall(smi)
	fmt.Println()
	goivb.FPrint(smi, 8)
}
