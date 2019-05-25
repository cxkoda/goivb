package main

import (
	"github.com/cxkoda/goivb"
)

func main() {
	goivb.GetStops("http://webservices.ivb.at/smiapi/1.0/Stops")
}