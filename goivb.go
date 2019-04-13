package goivb

import (
	"fmt"
	"github.com/levigross/grequests"
	"github.com/tidwall/gjson"
	"log"
	"strconv"
	"strings"
	"time"
	"os"
	"os/exec"
)

func GetStops () gjson.Result {
	resp, err := grequests.Get("http://webservices.ivb.at/smiapi/1.0/Stops", nil)

	if err != nil {
		log.Fatalln("Unable to make request: ", err)
	}

	data := resp.String()
	data = strings.Replace(data, "\\", "", -1)

	if !gjson.Valid(data) {
		log.Fatalln("invalid json")
	}

	return gjson.Parse(data).Get("#.stop")
}

func GetSmartinfo (stopUid int) gjson.Result {
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


var directionMaps = map[string]map[string]string {
	//←↑→↓
	"Höttinger Auffahrt": map[string]string{
		"Sadrach": "<-",
		"Rum Sanatorium": "<-",
		"J. Kerschb. Str.": "<-",
		"J.-Kerschbaumer-Straße": "<-",
		"Peerhofsiedlung": "->",
		"Technik West": "->",
		"Schützenstraße":"->",
		"Flughafen": "->",
		"Term. Marktplatz":"<-",
		"Terminal Marktplatz": "<-"}}

func RpiPrint(parser gjson.Result, num int) {
	directionMap := directionMaps[parser.Get("#.stopidname").Get("0").String()]
	smi := parser.Get("#.smartinfo")
	nCurr := 0
	smi.ForEach(func(key, value gjson.Result) bool {
		nCurr++
		if nCurr > num {
			return false
		}

		direction, ok := directionMap[value.Get("direction").String()]
		if !ok {
			direction = value.Get("direction").String()
		}
		time := value.Get("time").String()
		time = strings.Replace(time, " min", "\"", -1)

		fmt.Printf("%2v %v %-3v \n", value.Get("route"), direction, time)
		return true // keep iterating
	})
}


func Watchdog(stopUid int, sleep float64) {
	clearOut, _ := exec.Command("clear").Output()
	for true {
		smi := GetSmartinfo(stopUid)	
		os.Stdout.Write(clearOut)	
		RpiPrint(smi, 5)
		time.Sleep(time.Duration(sleep) * time.Second)
	}
}

