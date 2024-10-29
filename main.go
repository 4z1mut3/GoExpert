package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type ViaCEP struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Unidade     string `json:"unidade"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Estado      string `json:"estado"`
	Regiao      string `json:"regiao"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

func main() {
	http.HandleFunc("/BuscaCep", HandleBuscaCEP)
	http.ListenAndServe(":8080", nil)

}

func HandleBuscaCEP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/BuscaCep" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	CepParam := r.URL.Query().Get("cep")
	if CepParam == "" {
		w.WriteHeader(http.StatusBadRequest)
	}
	cep, err := BuscaCEP(CepParam)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(cep)
}

func BuscaCEP(cep string) (*ViaCEP, error) {
	resp, err := http.Get("https://viacep.com.br/ws/" + cep + "/json/")

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var c ViaCEP
	err = json.Unmarshal(body, &c)
	return &c, nil
}
