package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type Cotacao struct {
	Bid string `json:"bid"`
}

func main() {
	GetCotacaoEChamadaCriacaoArquivo()
}

func GetCotacaoEChamadaCriacaoArquivo() {
	c := http.Client{Timeout: 900 * time.Millisecond}

	res, err := c.Get("http://localhost:8080/GetCotacao")
	if err != nil {
		fmt.Println(http.StatusInternalServerError)
	}
	defer res.Body.Close()

	resp, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err.Error(), http.StatusInternalServerError)
		return
	}

	var cotacaoBRL Cotacao

	erro := json.Unmarshal(resp, &cotacaoBRL)
	if erro != nil {
		fmt.Println(http.StatusInternalServerError)
	}
	//fmt.Println(cotacaoBRL.Bid)
	criaArquivoTxt(cotacaoBRL)
}

func criaArquivoTxt(cot Cotacao) {
	// Criar um contexto com timeout de 3 segundos
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	//Criação de arquivo txt
	file, err := os.Create("cotacao.txt")
	if err != nil {
		fmt.Println("Erro ao criar arquivo!")
	}
	filename := "cotacao.txt"

	// Tenta gravar o arquivo
	if err := writeFile(ctx, filename, []byte(cot.Bid)); err != nil {
		if err == context.DeadlineExceeded {
			fmt.Println("Gravação cancelada: prazo excedido")
		} else {
			fmt.Println("Erro ao gravar o arquivo:", err)
		}
	} else {
		fmt.Println("Arquivo gravado com sucesso!")
	}
	//file.Write([]byte(cot.Bid))
	file.Close()
}

func writeFile(ctx context.Context, filename string, data []byte) error {
	// Cria um canal para sinalizar a conclusão da operação de gravação
	done := make(chan error)

	// Executa a operação de gravação em uma goroutine
	go func() {
		// Tenta escrever no arquivo
		err := ioutil.WriteFile(filename, data, 0644)
		done <- err
	}()

	// Espera pela operação de gravação ou pelo cancelamento do contexto
	select {
	case err := <-done:
		return err // Retorna o erro da gravação, se houver
	case <-ctx.Done():
		return ctx.Err() // Retorna o erro do contexto (cancelado ou timeout)
	}
}
