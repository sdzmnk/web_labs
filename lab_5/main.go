package main

import (
	"html/template"
	"net/http"
	"strconv"
)

type Result struct {
	// Завдання 1
	OmegaSumVal float64
	TvosVal     float64
	KaosVal     float64
	KposVal     float64
	DoubleOmega float64
	SystemOmega float64

	// Завдання 2
	Wneda float64
	Wnedp float64
	Zper  float64
}

func handler(w http.ResponseWriter, r *http.Request) {
	result := Result{}

	r.ParseForm()
	formType := r.FormValue("form_type")

	if formType == "task1" {
		PL110KVOmegaVal, _ := strconv.ParseFloat(r.FormValue("PL110KVOmegaVal"), 64)
		T110KVOmegaVal, _ := strconv.ParseFloat(r.FormValue("T110KVOmegaVal"), 64)
		V110KVOmegaVal, _ := strconv.ParseFloat(r.FormValue("V110KVOmegaVal"), 64)
		V10KVOmegaVal, _ := strconv.ParseFloat(r.FormValue("V10KVOmegaVal"), 64)
		tiresOmegaVal, _ := strconv.ParseFloat(r.FormValue("tiresOmegaVal"), 64)

		PL110KVTviVal, _ := strconv.ParseFloat(r.FormValue("PL110KVTviVal"), 64)
		T110KVTviVal, _ := strconv.ParseFloat(r.FormValue("T110KVTviVal"), 64)
		V110KVTviVal, _ := strconv.ParseFloat(r.FormValue("V110KVTviVal"), 64)
		V10KVTviVal, _ := strconv.ParseFloat(r.FormValue("V10KVTviVal"), 64)
		tiresTviVal, _ := strconv.ParseFloat(r.FormValue("tiresTviVal"), 64)

		PlannedKMaxVal, _ := strconv.ParseFloat(r.FormValue("PlannedKMaxVal"), 64)

		OmegaSumVal := PL110KVOmegaVal*10 + T110KVOmegaVal + V110KVOmegaVal + V10KVOmegaVal + 6*tiresOmegaVal
		TvosVal := (PL110KVOmegaVal*10*PL110KVTviVal + T110KVOmegaVal*T110KVTviVal + V110KVOmegaVal*V110KVTviVal + V10KVOmegaVal*V10KVTviVal + tiresOmegaVal*6*tiresTviVal) / OmegaSumVal
		KaosVal := (OmegaSumVal * TvosVal) / 8760
		KposVal := 1.2 * (PlannedKMaxVal / 8760)
		DoubleOmega := 2 * OmegaSumVal * (KaosVal + KposVal)
		SystemOmega := DoubleOmega + 0.02

		result.OmegaSumVal = OmegaSumVal
		result.TvosVal = TvosVal
		result.KaosVal = KaosVal
		result.KposVal = KposVal
		result.DoubleOmega = DoubleOmega
		result.SystemOmega = SystemOmega
	}

	if formType == "task2" {
		ZperA, _ := strconv.ParseFloat(r.FormValue("ZperA"), 64)
		ZperP, _ := strconv.ParseFloat(r.FormValue("ZperP"), 64)
		Omega, _ := strconv.ParseFloat(r.FormValue("Omega"), 64)
		tv, _ := strconv.ParseFloat(r.FormValue("tv"), 64)
		Pm, _ := strconv.ParseFloat(r.FormValue("Pm"), 64)
		Tm, _ := strconv.ParseFloat(r.FormValue("Tm"), 64)
		kp, _ := strconv.ParseFloat(r.FormValue("kp"), 64)

		Wneda := Omega * tv * Pm * Tm
		Wnedp := kp * Pm * Tm
		Zper := ZperA*Wneda + ZperP*Wnedp

		result.Wneda = Wneda
		result.Wnedp = Wnedp
		result.Zper = Zper
	}

	funcMap := template.FuncMap{
		"mul": func(a, b float64) float64 { return a * b },
	}

	tmpl := template.Must(template.New("index.html").Funcs(funcMap).ParseFiles("index.html"))
	tmpl.Execute(w, result)
}

func main() {
	http.Handle("/style.css", http.FileServer(http.Dir(".")))
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
