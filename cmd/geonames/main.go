package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"

	"github.com/jasontconnell/geonames/data"
)

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
		var modified []data.City
		if *mod != "" {
			mfcp, err := data.ReadJsonCities(*mod)
			if err != nil {
				log.Fatal(err)
			}
			modified = append(modified, mfcp...)
		}

		cities, err := data.ReadCities(*file, *includeCountries, modified)
		if err != nil {
			log.Fatal(err)
		}
		mncities, err := data.ReadJsonCities(*mn)
		if err != nil {
			log.Println("couldn't read manual cities, continuing.", *mn, err)
		}
		cities = append(cities, mncities...)
		jsonobj = cities
	} else if *format == "countries" {
		countries, err := data.ReadCountries(*file, *includeCountries)
		if err != nil {
			log.Fatal(err)
		}
		mncountries, err := data.ReadJsonCountries(*mn)
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
