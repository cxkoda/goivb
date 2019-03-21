package goivb

import (
	"fmt"
	"github.com/levigross/grequests"
	"github.com/tidwall/gjson"
	"log"
	"strconv"
	"strings"
)

func GetData (stopUid int) gjson.Result {
	resp, err := grequests.Post("http://webservices.ivb.at/smiapi/1.0/Passage/?stopUid=" +  strconv.Itoa(stopUid), nil)

	if err != nil {
		log.Fatalln("Unable to make request: ", err)
	}
	data := resp.String()
	data = strings.Replace(data, "\\", "", -1)

	if !gjson.Valid(data) {
		log.Fatalln("invalid json")
	}

	if (len(data) <= 2) {
		log.Fatalln("No data received. Wrong StopUid?")
	}

	return gjson.Parse(data)
}

func Printall(parser gjson.Result) {
	parser.ForEach(func(key, value gjson.Result) bool {
		fmt.Println(value.String())
		return true // keep iterating
	})
}

func FPrint(parser gjson.Result, num int) {
	rowSep := "|---------------------------------------------|\n"
	fmt.Printf(rowSep)
	fmt.Printf("| Haltestelle:%30s  |\n", parser.Get("#.stopidname").Get("0"))
	fmt.Printf(rowSep)

	smi := parser.Get("#.smartinfo")
	nCurr := 0
	smi.ForEach(func(key, value gjson.Result) bool {
		nCurr++
		if nCurr > num {
			return false
		}
		fmt.Printf("| %-5v| %-25v| %-10v|\n", value.Get("route"), value.Get("direction"), value.Get("time"))
		return true // keep iterating
	})

	fmt.Printf(rowSep)
}
