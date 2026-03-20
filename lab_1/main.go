package main

import (
	"encoding/json"
	"net/http"
)

type FuelInput struct {
	H float64 `json:"h"`
	C float64 `json:"c"`
	S float64 `json:"s"`
	N float64 `json:"n"`
	O float64 `json:"o"`
	W float64 `json:"w"`
	A float64 `json:"a"`
}

type FuelResult struct {
	KPC             float64            `json:"kpc"`
	KPG             float64            `json:"kpg"`
	QPH             float64            `json:"qph"`
	QCH             float64            `json:"qch"`
	QGH             float64            `json:"qgh"`
	DryComposition  map[string]float64 `json:"dry"`
	BurnComposition map[string]float64 `json:"burn"`
}

func calculateFuel(w http.ResponseWriter, r *http.Request) {
	var input FuelInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	KPC := 100 / (100 - input.W)
	KPG := 100 / (100 - input.W - input.A)

	QPH := 339*input.C + 1030*input.H - 108.8*(input.O-input.S) - 25*input.W
	QCH := (QPH/1000 + 0.025*input.W) * KPC
	QGH := (QPH/1000 + 0.025*input.W) * KPG

	dry := map[string]float64{
		"H": input.H * KPC, "C": input.C * KPC, "S": input.S * KPC,
		"N": input.N * KPC, "O": input.O * KPC, "A": input.A * KPC,
	}

	burn := map[string]float64{
		"H": input.H * KPG, "C": input.C * KPG, "S": input.S * KPG,
		"N": input.N * KPG, "O": input.O * KPG, "A": 0.0,
	}

	result := FuelResult{
		KPC: KPC, KPG: KPG, QPH: QPH, QCH: QCH, QGH: QGH,
		DryComposition: dry, BurnComposition: burn,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

type MazutInput struct {
	C float64 `json:"c"`
	H float64 `json:"h"`
	O float64 `json:"o"`
	S float64 `json:"s"`
	W float64 `json:"w"`
	A float64 `json:"a"`
	Q float64 `json:"q"`
	V float64 `json:"v"`
}

type MazutResult struct {
	C float64 `json:"c"`
	H float64 `json:"h"`
	O float64 `json:"o"`
	S float64 `json:"s"`
	A float64 `json:"a"`
	V float64 `json:"v"`
	Q float64 `json:"q"`
}

func calculateMazut(w http.ResponseWriter, r *http.Request) {
	var input MazutInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result := MazutResult{
		C: input.C * (100 - input.W - input.A) / 100,
		H: input.H * (100 - input.W - input.A) / 100,
		O: input.O * (100 - input.W/10 - input.A/10) / 100,
		S: input.S * (100 - input.W - input.A) / 100,
		A: input.A * (100 - input.W) / 100,
		V: input.V * (100 - input.W) / 100,
		Q: input.Q*(100-input.W-input.A)/100 - 0.025*input.W,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	http.HandleFunc("/calculate1", calculateFuel)
	http.HandleFunc("/calculate2", calculateMazut)

	http.ListenAndServe(":8080", nil)
}
