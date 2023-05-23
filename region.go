package main

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
)

type Region struct {
	Region []Stations `xml:"stations"`
}

type Stations struct {
	Stations []Station `xml:"station"`
}

type Station struct {
	ID     string `xml:"id"`
	Name   string `xml:"name"`
	AreaID string `xml:"area_id"`
}

const (
	RegionFullAPI = "http://radiko.jp/v3/station/region/full.xml"
)

// GetRegion returns Region
func GetRegion() (Region, error) {
	region := Region{}

	response, err := http.Get(RegionFullAPI)
	if err != nil {
		return region, err
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return region, err
	}
	defer response.Body.Close()

	err = xml.Unmarshal([]byte(string(body)), &region)
	if err != nil {
		return region, err
	}

	return region, nil
}

/*
func main() {
	region, err := GetRegion()
	if err != nil {
		panic(err)
	}
	for _, stations := range region.Region {
		for _, station := range stations.Stations {
			fmt.Println(station.Name)
		}
	}
}
*/
