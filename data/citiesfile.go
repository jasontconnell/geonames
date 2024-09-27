package data

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
	"strings"
)

var citiesFields map[string]int = map[string]int{
	"geonameid":         0,
	"name":              1,
	"asciiname":         2,
	"alternatenames":    3,
	"latitude":          4,
	"longitude":         5,
	"feature class":     6,
	"feature code":      7,
	"country code":      8,
	"cc2":               9,
	"admin1 code":       10,
	"admin2 code":       11,
	"admin3 code":       12,
	"admin4 code":       13,
	"population":        14,
	"elevation":         15,
	"dem":               16,
	"timezone":          17,
	"modification date": 18,
}

type City struct {
	Name           string   `json:"city"`
	AlternateNames []string `json:"alternateNames"`
	Latitude       float64  `json:"latitude"`
	Longitude      float64  `json:"longitude"`
	Country        string   `json:"country"`
	State          string   `json:"state"`
	TimeZone       string   `json:"timezone"`
}

func ReadCities(filename, includeCountries string, modified []City) ([]City, error) {
	lines, err := read(filename)
	if err != nil {
		return nil, err
	}

	inc, err := read(includeCountries)
	if err != nil {
		return nil, err
	}

	modmap := make(map[string]City)
	for _, m := range modified {
		modmap[m.Name+" "+m.State+" "+m.Country] = m
	}

	incLookup := make(map[string]bool)
	for _, line := range inc {
		incLookup[line] = true
	}

	cities := []City{}
	for _, line := range lines {
		sp := strings.Split(line, "\t")

		country := sp[citiesFields["country code"]]
		if _, ok := incLookup[country]; !ok {
			continue
		}

		lat, laterr := strconv.ParseFloat(sp[citiesFields["latitude"]], 64)
		lng, lngerr := strconv.ParseFloat(sp[citiesFields["longitude"]], 64)

		if laterr != nil || lngerr != nil {
			log.Println(laterr, lngerr)
			continue
		}

		splitAltNames := strings.Split(sp[citiesFields["alternatenames"]], ",")
		altNames := []string{}
		for _, s := range splitAltNames {
			if isASCII(s) {
				altNames = append(altNames, s)
			}
		}

		state := sp[citiesFields["admin1 code"]]
		if country != "US" {
			state = ""
		}

		city := City{
			Name:           sp[citiesFields["asciiname"]],
			AlternateNames: altNames,
			Latitude:       lat,
			Longitude:      lng,
			Country:        sp[citiesFields["country code"]],
			State:          state,
			TimeZone:       sp[citiesFields["timezone"]],
		}

		if mc, ok := modmap[city.Name+" "+city.State+" "+city.Country]; ok {
			if mc.AlternateNames != nil && len(mc.AlternateNames) > 0 {
				city.AlternateNames = append(city.AlternateNames, mc.AlternateNames...)
			}
		}

		cities = append(cities, city)
	}
	return cities, nil
}

func ReadJsonCities(filename string) ([]City, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	dec := json.NewDecoder(f)
	var manual []City
	err = dec.Decode(&manual)

	return manual, err
}
