package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type Requisicao struct {
	Texto string `json:"texto"`
}

type Resposta struct {
	Servidor string `json:"servidor"`
	Palavras int    `json:"palavras"`
	Vogais   int    `json:"vogais"`
	Letras   int    `json:"letras"`
	Numeros  int    `json:"numeros"`
}

type ParteLog struct {
	Servidor string `json:"servidor"`
	Texto    string `json:"texto"`
	Palavras int    `json:"palavras"`
	Vogais   int    `json:"vogais"`
	Letras   int    `json:"letras"`
	Numeros  int    `json:"numeros"`
}

type TotalLog struct {
	Palavras int `json:"palavras"`
	Vogais   int `json:"vogais"`
	Letras   int `json:"letras"`
	Numeros  int `json:"numeros"`
}

type Log struct {
	DataHora      string     `json:"dataHora"`
	TextoCompleto string     `json:"textoCompleto"`
	Partes        []ParteLog `json:"partes"`
	Total         TotalLog   `json:"total"`
}

type ResultadoAPI struct {
	Parte1     string   `json:"parte1"`
	Parte2     string   `json:"parte2"`
	Resposta1  Resposta `json:"resposta1"`
	Resposta2  Resposta `json:"resposta2"`
	Total      TotalLog `json:"total"`
	LogArquivo string   `json:"logArquivo"`
}

func enviarServidor(endereco string, texto string) (Resposta, error) {
	var resposta Resposta

	conexao, erro := net.DialTimeout("tcp", endereco, 5*time.Second)
	if erro != nil {
		return resposta, erro
	}
	defer conexao.Close()

	requisicao := Requisicao{Texto: texto}

	jsonRequisicao, erro := json.Marshal(requisicao)
	if erro != nil {
		return resposta, erro
	}

	fmt.Fprintln(conexao, string(jsonRequisicao))

	leitor := bufio.NewReader(conexao)
	mensagem, erro := leitor.ReadString('\n')
	if erro != nil {
		return resposta, erro
	}

	erro = json.Unmarshal([]byte(mensagem), &resposta)
	return resposta, erro
}

func dividirTexto(texto string) (string, string) {
	palavras := strings.Fields(texto)
	meio := len(palavras) / 2

	parte1 := strings.Join(palavras[:meio], " ")
	parte2 := strings.Join(palavras[meio:], " ")

	return parte1, parte2
}

func salvarLog(textoCompleto, parte1, parte2 string, resposta1, resposta2 Resposta) (string, error) {
	registro := Log{
		DataHora:      time.Now().Format(time.RFC3339),
		TextoCompleto: textoCompleto,
		Partes: []ParteLog{
			{Servidor: resposta1.Servidor, Texto: parte1, Palavras: resposta1.Palavras, Vogais: resposta1.Vogais, Letras: resposta1.Letras, Numeros: resposta1.Numeros},
			{Servidor: resposta2.Servidor, Texto: parte2, Palavras: resposta2.Palavras, Vogais: resposta2.Vogais, Letras: resposta2.Letras, Numeros: resposta2.Numeros},
		},
		Total: TotalLog{
			Palavras: resposta1.Palavras + resposta2.Palavras,
			Vogais:   resposta1.Vogais + resposta2.Vogais,
			Letras:   resposta1.Letras + resposta2.Letras,
			Numeros:  resposta1.Numeros + resposta2.Numeros,
		},
	}

	dadosJSON, erro := json.MarshalIndent(registro, "", "    ")
	if erro != nil {
		return "", erro
	}

	if erro := os.MkdirAll("logs", 0755); erro != nil {
		return "", erro
	}

	nomeArquivo := fmt.Sprintf("logs/log_%s.json", time.Now().Format("2006-01-02_15-04-05"))
	if erro := os.WriteFile(nomeArquivo, dadosJSON, 0644); erro != nil {
		return "", erro
	}

	return nomeArquivo, nil
}

func processar(texto string) (ResultadoAPI, error) {
	var resultado ResultadoAPI

	texto = strings.TrimSpace(texto)
	if texto == "" {
		return resultado, fmt.Errorf("O texto não pode estar vazio")
	}

	parte1, parte2 := dividirTexto(texto)

	var resposta1, resposta2 Resposta
	var erro1, erro2 error
	var waitGroup sync.WaitGroup

	waitGroup.Add(2)
	go func() {
		defer waitGroup.Done()
		resposta1, erro1 = enviarServidor("localhost:8001", parte1)
	}()
	go func() {
		defer waitGroup.Done()
		resposta2, erro2 = enviarServidor("localhost:8002", parte2)
	}()
	waitGroup.Wait()

	if erro1 != nil {
		return resultado, fmt.Errorf("erro no Servidor 1: %w", erro1)
	}
	if erro2 != nil {
		return resultado, fmt.Errorf("erro no Servidor 2: %w", erro2)
	}

	nomeArquivo, erro := salvarLog(texto, parte1, parte2, resposta1, resposta2)
	if erro != nil {
		log.Println("Não foi possível salvar o log:", erro)
	}

	resultado = ResultadoAPI{
		Parte1:    parte1,
		Parte2:    parte2,
		Resposta1: resposta1,
		Resposta2: resposta2,
		Total: TotalLog{
			Palavras: resposta1.Palavras + resposta2.Palavras,
			Vogais:   resposta1.Vogais + resposta2.Vogais,
			Letras:   resposta1.Letras + resposta2.Letras,
			Numeros:  resposta1.Numeros + resposta2.Numeros,
		},
		LogArquivo: nomeArquivo,
	}

	return resultado, nil
}

func handlerProcessar(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	var requisicao Requisicao
	if erro := json.NewDecoder(r.Body).Decode(&requisicao); erro != nil {
		http.Error(w, "JSON inválido: "+erro.Error(), http.StatusBadRequest)
		return
	}

	resultado, erro := processar(requisicao.Texto)
	if erro != nil {
		http.Error(w, erro.Error(), http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resultado)
}

func main() {
	porta := os.Getenv("WEB_PORT")
	if porta == "" {
		porta = "8080"
	}

	http.HandleFunc("/api/processar", handlerProcessar)
	http.Handle("/", http.FileServer(http.Dir("./client/interface")))

	endereco := "localhost:" + porta
	fmt.Printf("Aberto em http://%s\n", endereco)
	fmt.Println("Os servidores de localhost:8001 e localhost:8002 precisam estar de pé.")

	if erro := http.ListenAndServe(endereco, nil); erro != nil {
		log.Fatal(erro)
	}
}
