package data

import (
	"encoding/json"
	"os"
	"strings"
)

var countriesFields map[string]int = map[string]int{
	"iso":     0,
	"iso3":    1,
	"country": 4,
}

type Country struct {
	Name  string `json:"name"`
	Abbr  string `json:"let2"`
	Abbr3 string `json:"let3"`
}

func ReadCountries(filename string) ([]Country, error) {
	lines, err := read(filename)
	if err != nil {
		return nil, err
	}
	countries := []Country{}
	for _, line := range lines {
		if line[0] == '#' {
			continue
		}
		sp := strings.Split(line, "\t")

		abbr := sp[countriesFields["iso"]]
		country := Country{
			Name:  sp[countriesFields["country"]],
			Abbr:  abbr,
			Abbr3: sp[countriesFields["iso3"]],
		}
		countries = append(countries, country)
	}
	return countries, nil
}

func ReadCountriesFiltered(filename, includeCountries string) ([]Country, error) {
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
			Abbr:  abbr,
			Abbr3: sp[countriesFields["iso3"]],
		}
		countries = append(countries, country)
	}
	return countries, nil
}

func ReadJsonCountries(filename string) ([]Country, error) {
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
