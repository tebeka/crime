package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"gonum.org/v1/gonum/stat"
)

func load(csvFile string, countryCol, valueCol int) (map[string]float64, error) {
	file, err := os.Open(csvFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	r := csv.NewReader(file)
	r.LazyQuotes = true
	byCountry := make(map[string]float64) // country -> value
	lnum := 0
	for {
		fields, err := r.Read()
		if err == io.EOF {
			break
		}

		lnum++
		if err != nil {
			return nil, fmt.Errorf("%d: %s", lnum, err)
		}

		if lnum == 1 { // Skip header
			continue
		}

		// 1,083 -> 1083
		s := strings.Replace(fields[valueCol], ",", "", -1)
		value, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return nil, fmt.Errorf("%d: %s", lnum, err)
		}
		country := fields[countryCol]
		i := strings.Index(country, "\u202f*")
		if i != -1 {
			country = country[:i]
		}

		byCountry[country] = value
	}

	return byCountry, nil
}

func main() {
	// "ranking","country","crimeIndex","pop2022"
	crime, err := load("crime.csv", 1, 2)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println(crime)

	// Country or dependency,Number of officers,Rate.mw-parser-output .nobold{font-weight:normal}(per 100k people),Year
	police, err := load("police.csv", 0, 2)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println(police)

	var (
		cVals []float64
		pVals []float64
	)

	for country, cVal := range crime {
		pVal, ok := police[country]
		if !ok {
			continue
		}
		cVals = append(cVals, cVal)
		pVals = append(pVals, pVal)
	}

	corr := stat.Correlation(cVals, pVals, nil)
	fmt.Printf("police ration to crime correlation: %.2f\n", corr)
}
