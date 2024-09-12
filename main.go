package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode"
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

var countriesFields map[string]int = map[string]int{
	"iso":     0,
	"iso3":    1,
	"country": 4,
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

type Country struct {
	Name  string `json:"name"`
	Abbr  string `json:"let2"`
	Abbr3 string `json:"let3"`
}

func main() {
	file := flag.String("f", "", "filepath")
	out := flag.String("out", "cities.json", "output filename (json)")
	format := flag.String("format", "cities", "cities or countries")
	includeCountries := flag.String("inc", "includeCountries.txt", "list of countries to include")
	mn := flag.String("manual", "manual.json", "manually added cities")
	mod := flag.String("mod", "modifycities.json", "modify cities")
	flag.Parse()

	var jsonobj interface{}

	if *format == "cities" {
		var modified []City
		if *mod != "" {
			mfcp, err := readManualCities(*mod)
			if err != nil {
				log.Fatal(err)
			}
			modified = append(modified, mfcp...)
		}

		cities, err := readCities(*file, *includeCountries, modified)
		if err != nil {
			log.Fatal(err)
		}
		mncities, err := readManualCities(*mn)
		if err != nil {
			log.Println("couldn't read manual cities, continuing.", *mn, err)
		}
		cities = append(cities, mncities...)
		jsonobj = cities
	} else if *format == "countries" {
		countries, err := readCountries(*file, *includeCountries)
		if err != nil {
			log.Fatal(err)
		}
		mncountries, err := readManualCountries(*mn)
		if err != nil {
			log.Println("couldn't read manual countries, continuing.", *mn, err)
		}
		countries = append(countries, mncountries...)
		jsonobj = countries
	}

	o, err := os.OpenFile(*out, os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	defer o.Close()
	enc := json.NewEncoder(o)
	err = enc.Encode(jsonobj)
	if err != nil {
		log.Fatal(err)
	}
}

func readCities(filename, includeCountries string, modified []City) ([]City, error) {
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

func readCountries(filename, includeCountries string) ([]Country, error) {
	lines, err := read(filename)
	if err != nil {
		return nil, err
	}

	inc, err := read(includeCountries)
	if err != nil {
		return nil, err
	}

	incLookup := make(map[string]bool)
	for _, line := range inc {
		incLookup[line] = true
	}

	countries := []Country{}
	for _, line := range lines {
		if line[0] == '#' {
			continue
		}
		sp := strings.Split(line, "\t")

		abbr := sp[countriesFields["iso"]]
		if _, ok := incLookup[abbr]; !ok {
			continue
		}

		country := Country{
			Name:  sp[countriesFields["country"]],
			Abbr:  sp[countriesFields["iso"]],
			Abbr3: sp[countriesFields["iso3"]],
		}
		countries = append(countries, country)
	}
	return countries, nil
}

func read(filename string) ([]string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	lines := []string{}
	for s.Scan() {
		line := s.Text()
		lines = append(lines, line)
	}
	return lines, nil
}

func readManualCities(filename string) ([]City, error) {
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

func readManualCountries(filename string) ([]Country, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	dec := json.NewDecoder(f)
	var manual []Country
	err = dec.Decode(&manual)

	return manual, err
}

func isASCII(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] > unicode.MaxASCII {
			return false
		}
	}
	return true
}
