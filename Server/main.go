package main

import (
	"encoding/json"
	"io"
	"net/http"
)

type Cotacao struct {
	Usdbrl struct {
		Code       string `json:"code"`
		Codein     string `json:"codein"`
		Name       string `json:"name"`
		High       string `json:"high"`
		Low        string `json:"low"`
		VarBid     string `json:"varBid"`
		PctChange  string `json:"pctChange"`
		Bid        string `json:"bid"`
		Ask        string `json:"ask"`
		Timestamp  string `json:"timestamp"`
		CreateDate string `json:"create_date"`
	} `json:"USDBRL"`
}

func main() {

	mux := http.NewServeMux()

	mux.HandleFunc("/GetCotacao", GetCotacao)
	http.ListenAndServe(":8081", mux)

}

func GetCotacao(w http.ResponseWriter, r *http.Request) {
	c := http.Client{}

	if r.URL.Path != "/GetCotacao" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	req, err := c.Get("https://economia.awesomeapi.com.br/json/last/usd-brl")
	defer req.Body.Close()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res, err := io.ReadAll(req.Body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var cotacaoBRL Cotacao

	erro := json.Unmarshal(res, &cotacaoBRL)

	if erro != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	bidResponse := map[string]string{"bid": cotacaoBRL.Usdbrl.Bid}
	jsonResponse, err := json.Marshal(bidResponse)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(jsonResponse))

}
