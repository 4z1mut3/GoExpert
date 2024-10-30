package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

type Cotacao struct {
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
}

func main() {
	cotacao := Cotacao{
		Code:       "USD",
		Codein:     "BRL",
		Name:       "Dólar",
		High:       "5.30",
		Low:        "5.00",
		VarBid:     "0.10",
		PctChange:  "1.89",
		Bid:        "5.25",
		Ask:        "5.35",
		Timestamp:  "1638312000",
		CreateDate: "2023-10-29",
	}

	criarTabelaCotacao()
	insertCotacao(cotacao)

}

func GetCotacao(w http.ResponseWriter, r *http.Request) {
	c := http.Client{}

	if r.URL.Path != "/GetCotacao" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	req, err := c.Get("https://localhost:8080/GetCotacao")
	defer req.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	var cotacaoBRL Cotacao

	erro := json.Unmarshal(r, &cotacaoBRL)
	if erro != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func insertCotacao(cot Cotacao) {
	db, err := sql.Open("sqlite3", "./cotacao.db")
	if err != nil {
		fmt.Println("Erro ao abrir o banco de dados:", err)
		return
	}
	// Inserir a moeda na tabela
	_, err = db.Exec(`
       INSERT INTO Cotacao (
           code, codein, name, high, low, varBid, pctChange, bid, ask, timestamp, create_date
       ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
   `, cot.Code, cot.Codein, cot.Name, cot.High, cot.Low,
		cot.VarBid, cot.PctChange, cot.Bid, cot.Ask,
		cot.Timestamp, cot.CreateDate)

	if err != nil {
		fmt.Println("Erro ao inserir cotacao em  banco de dados:", err)
		return
	}
}

func criarTabelaCotacao() {

	db, err := sql.Open("sqlite3", "./cotacao.db")
	if err != nil {
		fmt.Println("Erro ao abrir o banco de dados:", err)
		return
	}
	defer db.Close()

	// Criar tabela
	sqlStmt := `
    CREATE TABLE IF NOT EXISTS Cotacao (
        code TEXT,
        codein TEXT,
        name TEXT,
        high TEXT,
        low TEXT,
        varBid TEXT,
        pctChange TEXT,
        bid TEXT,
        ask TEXT,
        timestamp TEXT,
        create_date TEXT
    );
    `
	_, err = db.Exec(sqlStmt)
	if err != nil {
		fmt.Println("Erro ao criar a tabela:", err)
		return
	}

	fmt.Println("Tabela criada com sucesso!")
}
