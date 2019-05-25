package goivb

import (
	"github.com/levigross/grequests"
	"github.com/tidwall/gjson"
	"log"
	"strings"
	"fmt"
	"encoding/json"
)

type StopT struct {
	Uid  string `json:"uid"`
	Name string `json:"name"`
	Lat  string `json:"lat"`
	Lon  string `json:"lon"`
}

type Stops []struct {
	stops StopT `json:"stop"`
}

type stopResponse []struct {
	Stop struct {
		UID  string `json:"uid"`
		Name string `json:"name"`
		Lat  string `json:"lat"`
		Lon  string `json:"lon"`
	} `json:"stop"`
}


func GetStops (hostname string) Stops {
	if !Config.IsSet {
		log.Fatalln("Goivb Configuration not set")
	}

	resp, err := grequests.Get(hostname, nil)

	if err != nil {
		log.Fatalln("Unable to make request: ", err)
	}

	data := resp.String()
	data = strings.Replace(data, "\\", "", -1)

	if !gjson.Valid(data) {
		log.Fatalln("invalid json")
	}

	fmt.Println(data)

	//stops := Stops{}

	response := stopResponse{}
	json.Unmarshal([]byte(data), &response)

	fmt.Println(response)

	//return Stops(gjson.Parse(data).Get("#.stop"))
	return Stops{}
}
