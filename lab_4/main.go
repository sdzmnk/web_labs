package main

import (
	"html/template"
	"math"
	"net/http"
	"strconv"
)

type Result struct {
	ImVal    float64
	ImPaVal  float64
	SekVal   float64
	SSMinVal float64

	SumXVal float64
	Ip0Val  float64

	XtVal      float64
	ZshVal     float64
	XshVal     float64
	ZshMinVal  float64
	XshMinVal  float64
	Ish3Val    float64
	Ish2Val    float64
	Ish3MinVal float64
	Ish2MinVal float64
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl := template.Must(template.ParseFiles("index.html"))
		tmpl.Execute(w, nil)
		return
	}

	SmVal, _ := strconv.ParseFloat(r.FormValue("SmVal"), 64)
	IkVal, _ := strconv.ParseFloat(r.FormValue("IkVal"), 64)
	tPhiVal, _ := strconv.ParseFloat(r.FormValue("tPhiVal"), 64)

	ImVal := (SmVal / 2) / (math.Sqrt(3) * 10)
	ImPaVal := 2 * ImVal
	SekVal := ImVal / 1.4
	SSMinVal := (IkVal * math.Sqrt(tPhiVal)) / 92

	KZPower, _ := strconv.ParseFloat(r.FormValue("KZPower"), 64)

	sumXVal := (10.5*10.5)/KZPower + (10.5/100)*((10.5*10.5)/6.3)
	Ip0Val := 10.5 / (math.Sqrt(3) * sumXVal)

	Ukmax, _ := strconv.ParseFloat(r.FormValue("Ukmax"), 64)
	Uvn, _ := strconv.ParseFloat(r.FormValue("Uvn"), 64)
	Snom, _ := strconv.ParseFloat(r.FormValue("Snom"), 64)
	Rsh, _ := strconv.ParseFloat(r.FormValue("Rsh"), 64)
	Xsh, _ := strconv.ParseFloat(r.FormValue("Xsh"), 64)
	RshMin, _ := strconv.ParseFloat(r.FormValue("RshMin"), 64)
	XshMin, _ := strconv.ParseFloat(r.FormValue("XshMin"), 64)

	Xt := (Ukmax / 100) * (Uvn * Uvn / Snom)

	Zsh := math.Sqrt(Rsh*Rsh + math.Pow(Xt+Xsh, 2))
	ZshMin := math.Sqrt(RshMin*RshMin + math.Pow(Xt+XshMin, 2))

	Ish3 := (Uvn * 1000) / (math.Sqrt(3) * Zsh)
	Ish2 := Ish3 * (math.Sqrt(3) / 2)

	Ish3Min := (Uvn * 1000) / (math.Sqrt(3) * ZshMin)
	Ish2Min := Ish3Min * (math.Sqrt(3) / 2)

	result := Result{
		ImVal: ImVal, ImPaVal: ImPaVal, SekVal: SekVal, SSMinVal: SSMinVal,
		SumXVal: sumXVal, Ip0Val: Ip0Val,
		XtVal: Xt, ZshVal: Zsh, XshVal: Xt + Xsh,
		ZshMinVal: ZshMin, XshMinVal: Xt + XshMin,
		Ish3Val: Ish3, Ish2Val: Ish2,
		Ish3MinVal: Ish3Min, Ish2MinVal: Ish2Min,
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
