package goivb

import (
	"fmt"
	"github.com/levigross/grequests"
	"github.com/tidwall/gjson"
	"github.com/BurntSushi/toml"
	"log"
	"strconv"
	"strings"
	"time"
	"os"
	"os/exec"
)

// to change the flags on the default logger
var glog = log.New(os.Stderr, "", log.LstdFlags | log.Lshortfile)

type GoivbConfig struct {
	StopHost string `toml:"StopHost"`
	PassageHost string `toml:"PassageHost"`
}

type WatchdogConfig struct {
	StopUid int `toml:"StopUid"`
	Sleep float64 `toml:"Sleep"`
}

type TomlConfig struct {
	Goivb GoivbConfig `toml:"goivb"`
	Watchdogs map[string]WatchdogConfig `toml:"watchdog"`
	IsSet bool
}

var Config TomlConfig

func init ()  {

	if _, err := os.Stat("/etc/goivb.toml"); err == nil {	  
	if _, err := toml.DecodeFile("/etc/goivb.toml", &Config); err != nil {
		glog.Fatalln(err)
	} else {
		Config.IsSet = true
	}
	}

	if _, err := os.Stat("goivb.toml"); err == nil {
	if _, err := toml.DecodeFile("goivb.toml", &Config); err != nil {
		glog.Fatalln(err)
	} else {
		Config.IsSet = true
	}
	}

	if !Config.IsSet {
		glog.Fatalln("No Config file loaded")
	}
}


type GoivbStops gjson.Result
type GoivbSmi gjson.Result

func GetStops () gjson.Result {
	if !Config.IsSet {
		glog.Fatalln("Goivb Configuration not set")
	}

	resp, err := grequests.Get(Config.Goivb.StopHost, nil)

	if err != nil {
		glog.Fatalln("Unable to make request: ", err)
	}

	data := resp.String()
	data = strings.Replace(data, "\\", "", -1)

	if !gjson.Valid(data) {
		glog.Fatalln("invalid json")
	}

	return gjson.Parse(data).Get("#.stop")
}

func GetSmartinfo (stopUid int) gjson.Result {
	if !Config.IsSet {
		glog.Fatalln("Goivb Configuration not set")
	}

	resp, err := grequests.Post(Config.Goivb.PassageHost + "/?stopUid=" +  strconv.Itoa(stopUid), nil)

	if err != nil {
		glog.Fatalln("Unable to make request: ", err)
	}
	data := resp.String()
	data = strings.Replace(data, "\\", "", -1)

	if !gjson.Valid(data) {
		glog.Fatalln("invalid json")
	}

	if (len(data) <= 2) {
		glog.Fatalln("No data received. Wrong StopUid?")
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
		"Sadrach": "< ",
		"Rum Sanatorium": "< ",
		"J. Kerschb. Str.": "< ",
		"J.-Kerschbaumer-Straße": "< ",
		"Peerhofsiedlung": " >",
		"Technik West": " >",
		"Schützenstraße":" >",
		"Flughafen": " >",
		"Term. Marktplatz": "< ",
		"Terminal Marktplatz": "< ",
		"Kajetan-Sweth-Straße": "< ",
		"Technik": " >",
		"Ibk. Hauptbahnhof": "< "}}

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
		time = strings.Replace(time, " min", "'", -1)

		fmt.Printf("%v %v %v \n", value.Get("route"), direction, time)
		return true // keep iterating
	})
}


func Watchdog(watchdogId string) {
	cfg := Config.Watchdogs[watchdogId]
	clearOut, _ := exec.Command("clear").Output()
	for true {
		smi := GetSmartinfo(cfg.StopUid)
		os.Stdout.Write(clearOut)
		RpiPrint(smi, 6)
		time.Sleep(time.Duration(cfg.Sleep) * time.Second)
	}
}

