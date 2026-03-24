package main

import (
	"html/template"
	"math"
	"net/http"
	"strconv"
)

type InputData struct {
	ShrSumPH     float64
	ShrSumPHKv   float64
	ShrSumPHKvTg float64
	ShrSumPH2    float64
	ShrKp        float64

	ShopSumPH     float64
	ShopSumPHKv   float64
	ShopSumPHKvTg float64
	ShopSumPH2    float64
	ShopKp        float64
}

type CalculationResult struct {
	ShrKv float64
	ShrNe float64
	ShrPp float64
	ShrQp float64
	ShrSp float64
	ShrIp float64

	ShopKv float64
	ShopNe float64
	ShopPp float64
	ShopQp float64
	ShopSp float64
	ShopIp float64
}

func main() {
	tmpl := template.Must(template.ParseFiles("index.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			tmpl.Execute(w, nil)
			return
		}

		parseFloat := func(key string) float64 {
			val, _ := strconv.ParseFloat(r.FormValue(key), 64)
			return val
		}

		in := InputData{
			ShrSumPH:      parseFloat("shr_sum_ph"),
			ShrSumPHKv:    parseFloat("shr_sum_ph_kv"),
			ShrSumPHKvTg:  parseFloat("shr_sum_ph_kv_tg"),
			ShrSumPH2:     parseFloat("shr_sum_ph2"),
			ShrKp:         parseFloat("shr_kp"),
			ShopSumPH:     parseFloat("shop_sum_ph"),
			ShopSumPHKv:   parseFloat("shop_sum_ph_kv"),
			ShopSumPHKvTg: parseFloat("shop_sum_ph_kv_tg"),
			ShopSumPH2:    parseFloat("shop_sum_ph2"),
			ShopKp:        parseFloat("shop_kp"),
		}

		var res CalculationResult

		res.ShrKv = in.ShrSumPHKv / in.ShrSumPH
		res.ShrNe = math.Ceil(math.Pow(in.ShrSumPH, 2) / in.ShrSumPH2)
		res.ShrPp = in.ShrKp * in.ShrSumPHKv
		res.ShrQp = 1.0 * in.ShrSumPHKvTg
		res.ShrSp = math.Sqrt(math.Pow(res.ShrPp, 2) + math.Pow(res.ShrQp, 2))
		res.ShrIp = res.ShrPp / 0.38

		res.ShopKv = in.ShopSumPHKv / in.ShopSumPH
		res.ShopNe = math.Floor(math.Pow(in.ShopSumPH, 2) / in.ShopSumPH2)
		res.ShopPp = in.ShopKp * in.ShopSumPHKv
		res.ShopQp = in.ShopKp * in.ShopSumPHKvTg
		res.ShopSp = math.Sqrt(math.Pow(res.ShopPp, 2) + math.Pow(res.ShopQp, 2))
		res.ShopIp = res.ShopPp / 0.38

		tmpl.Execute(w, res)
	})

	http.ListenAndServe(":8080", nil)
}
