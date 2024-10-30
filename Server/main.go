package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3" // Importação do driver
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
	CriarTabelaCotacao()
	mux := http.NewServeMux()

	mux.HandleFunc("/GetCotacao", GetCotacao)
	http.ListenAndServe(":8080", mux)

}

func GetCotacao(w http.ResponseWriter, r *http.Request) {
	c := http.Client{Timeout: 800 * time.Millisecond}

	if r.URL.Path != "/GetCotacao" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	req, err := c.Get("https://economia.awesomeapi.com.br/json/last/usd-brl")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer req.Body.Close()

	res, err := io.ReadAll(req.Body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var Cot Cotacao

	erro := json.Unmarshal(res, &Cot)

	if erro != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	bidResponse := map[string]string{"bid": Cot.Usdbrl.Bid}
	jsonResponse, err := json.Marshal(bidResponse)

	InsertCotacao(Cot)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(jsonResponse))

}

func InsertCotacao(cot Cotacao) {

	// Criar um contexto com um prazo
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel() // Garante que o cancelamento será chamado

	db, err := sql.Open("sqlite3", "./cotacao.db")
	if err != nil {
		fmt.Println("Erro ao abrir o banco de dados:", err)
		return
	}
	// Inserir a moeda na tabela
	_, err = db.ExecContext(ctx, `
       INSERT INTO Cotacao (
           code, codein, name, high, low, varBid, pctChange, bid, ask, timestamp, create_date
       ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
   `, cot.Usdbrl.Code, cot.Usdbrl.Codein, cot.Usdbrl.Name, cot.Usdbrl.High, cot.Usdbrl.Low,
		cot.Usdbrl.VarBid, cot.Usdbrl.PctChange, cot.Usdbrl.Bid, cot.Usdbrl.Ask,
		cot.Usdbrl.Timestamp, cot.Usdbrl.CreateDate)

	if err != nil {
		fmt.Println("Erro ao inserir cotacao em  banco de dados:", err)
		return
	}
}

func CriarTabelaCotacao() {

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
