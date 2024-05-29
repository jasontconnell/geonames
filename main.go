package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"log"
	"os"
	"strconv"
	"strings"
)

var fields map[string]int = map[string]int{
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
	Name      string  `json:"city"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Country   string  `json:"country"`
	State     string  `json:"state"`
	TimeZone  string  `json:"timezone"`
}

func main() {
	file := flag.String("f", "", "filepath")
	out := flag.String("out", "cities.json", "output filename (json)")
	flag.Parse()

	lines, err := read(*file)
	if err != nil {
		log.Fatal(err)
	}

	cities := []City{}
	for _, line := range lines {
		sp := strings.Split(line, "\t")
		lat, perr := strconv.ParseFloat(sp[fields["latitude"]], 64)
		if perr != nil {
			log.Println(perr)
		}
		lng, perr := strconv.ParseFloat(sp[fields["longitude"]], 64)
		if perr != nil {
			log.Println(perr)
		}
		city := City{
			Name:      sp[fields["name"]],
			Latitude:  lat,
			Longitude: lng,
			Country:   sp[fields["country code"]],
			State:     sp[fields["admin1 code"]],
			TimeZone:  sp[fields["timezone"]],
		}
		cities = append(cities, city)
	}

	o, err := os.OpenFile(*out, os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	defer o.Close()
	enc := json.NewEncoder(o)
	enc.Encode(cities)
}

func read(filename string) ([]string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	s := bufio.NewScanner(f)
	lines := []string{}
	for s.Scan() {
		line := s.Text()
		lines = append(lines, line)
	}
	return lines, nil
}
