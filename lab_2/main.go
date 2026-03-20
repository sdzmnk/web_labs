package main

import (
	"encoding/json"
	"math"
	"net/http"
)

type FuelInput struct {
	Coal       float64 `json:"coal"`
	FuelOil    float64 `json:"fuelOil"`
	NaturalGas float64 `json:"naturalGas"`
}

type EmissionResult struct {
	CoalEI    float64 `json:"coalEI"`
	CoalTotal float64 `json:"coalTotal"`
	OilEI     float64 `json:"oilEI"`
	OilTotal  float64 `json:"oilTotal"`
	GasEI     float64 `json:"gasEI"`
	GasTotal  float64 `json:"gasTotal"`
}

func calculateEmissions(w http.ResponseWriter, r *http.Request) {
	var input FuelInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	const (
		qCoal    = 20.47
		aVunCoal = 0.8
		aCoal    = 25.2
		gVunCoal = 1.5
		nzu      = 0.985

		qOil    = 40.40
		wOil    = 2.0
		aOil    = 0.15
		aVunOil = 1.0
		gVunOil = 0.0
	)

	coalEI := (math.Pow(10, 6) / qCoal) * aVunCoal * (aCoal / (100 - gVunCoal)) * (1 - nzu)
	coalTotal := math.Pow(10, -6) * coalEI * qCoal * input.Coal

	qOilWorking := qOil*(100-wOil-aOil)/100 - 0.025*wOil
	oilEI := (math.Pow(10, 6) / qOilWorking) * aVunOil * (aOil / (100 - gVunOil)) * (1 - nzu)
	oilTotal := math.Pow(10, -6) * oilEI * qOilWorking * input.FuelOil

	res := EmissionResult{
		CoalEI:    math.Round(coalEI*100) / 100,
		CoalTotal: math.Round(coalTotal*100) / 100,
		OilEI:     math.Round(oilEI*100) / 100,
		OilTotal:  math.Round(oilTotal*100) / 100,
		GasEI:     0.0,
		GasTotal:  0.0,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func main() {
	http.Handle("/", http.FileServer(http.Dir(".")))
	http.HandleFunc("/calculate", calculateEmissions)
	http.ListenAndServe(":8080", nil)
}
