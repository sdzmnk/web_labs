package main

import (
	"html/template"
	"math"
	"net/http"
	"strconv"
)

type Result struct {
	GroupUtilization     float64
	EffectiveNumEP       float64
	ActiveLoad           float64
	ReactiveLoad         float64
	TotalPower           float64
	GroupCurrent         float64
	GroupUtilizationShop float64
	EffectiveNumEPShop   float64
	ActiveLoadShop       float64
	ReactiveLoadShop     float64
	TotalPowerShop       float64
	GroupCurrentShop     float64
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl := template.Must(template.ParseFiles("index.html"))
		tmpl.Execute(w, nil)
		return
	}

	powerEP1, _ := strconv.ParseFloat(r.FormValue("powerEP1"), 64)
	utilEP6, _ := strconv.ParseFloat(r.FormValue("utilEP6"), 64)
	reactiveEP4, _ := strconv.ParseFloat(r.FormValue("reactiveEP4"), 64)

	numEP := []float64{4, 2, 4, 1, 1, 1, 1, 2}
	powerEP := []float64{powerEP1, 14, 42, 36, 20, 40, 32, 20}
	utilization := []float64{0.15, 0.12, 0.15, 0.3, 0.5, utilEP6, 0.2, 0.65}
	reactive := []float64{1.33, 1, 1.33, reactiveEP4, 0.75, 1, 1, 0.75}

	var totalPowerSum, sumPowerSum float64
	totalPower := make([]float64, len(numEP))

	for i := range numEP {
		totalPower[i] = numEP[i] * powerEP[i]
		sumPowerSum += numEP[i] * math.Pow(powerEP[i], 2)
		totalPowerSum += totalPower[i]
	}

	powerCoefficient := 1.25
	loadVoltage := 0.38

	var weightedSum float64
	var reactiveSum float64

	for i := range totalPower {
		weightedSum += totalPower[i] * utilization[i]
		reactiveSum += totalPower[i] * utilization[i] * reactive[i]
	}

	groupUtilization := weightedSum / totalPowerSum
	effectiveNumEP := math.Pow(totalPowerSum, 2) / sumPowerSum
	activeLoad := powerCoefficient * weightedSum
	reactiveLoad := reactiveSum

	totalPowerResult := math.Sqrt(math.Pow(activeLoad, 2) + math.Pow(reactiveLoad, 2))
	groupCurrent := activeLoad / loadVoltage

	powerShop := 2330.0
	powerUtil := 752.0
	powerReactive := 657.0
	powerTotal := 96388.0
	powerCoeffShop := 0.7

	activeLoadShop := powerCoeffShop * powerUtil
	reactiveLoadShop := powerCoeffShop * powerReactive
	groupUtilShop := powerUtil / powerShop
	effectiveShop := math.Pow(powerShop, 2) / powerTotal
	totalShop := math.Sqrt(math.Pow(activeLoadShop, 2) + math.Pow(reactiveLoadShop, 2))
	groupCurrentShop := activeLoadShop / loadVoltage

	result := Result{
		GroupUtilization:     groupUtilization,
		EffectiveNumEP:       effectiveNumEP,
		ActiveLoad:           activeLoad,
		ReactiveLoad:         reactiveLoad,
		TotalPower:           totalPowerResult,
		GroupCurrent:         groupCurrent,
		GroupUtilizationShop: groupUtilShop,
		EffectiveNumEPShop:   effectiveShop,
		ActiveLoadShop:       activeLoadShop,
		ReactiveLoadShop:     reactiveLoadShop,
		TotalPowerShop:       totalShop,
		GroupCurrentShop:     groupCurrentShop,
	}

	tmpl := template.Must(template.ParseFiles("index.html"))
	tmpl.Execute(w, result)
}

func main() {
	http.HandleFunc("/style.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "style.css")
	})

	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
