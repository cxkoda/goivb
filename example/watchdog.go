package main

import (
	"github.com/cxkoda/goivb"
	"os"
)


func main() {
	goivb.Watchdog(os.Args[1])
}
