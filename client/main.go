package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
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

func enviarParaServidor(endereco string, texto string) (Resposta, error) {
	var resposta Resposta

	conexao, erro := net.Dial("tcp", endereco)
	if erro != nil {
		return resposta, erro
	}

	defer conexao.Close()

	requisicao := Requisicao{
		Texto: texto,
	}

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

func main() {
	leitor := bufio.NewReader(os.Stdin)
	fmt.Println("Digite um texto:")

	texto, erro := leitor.ReadString('\n')

	if erro != nil {
		fmt.Println("Erro ao ler texto:", erro)
		return
	}

	texto = strings.TrimSpace(texto)

	if texto == "" {
		fmt.Println("O texto não pode estar vazio.")
		return
	}

	parte1, parte2 := dividirTexto(texto)

	var resposta1 Resposta
	var resposta2 Resposta

	var erro1 error
	var erro2 error

	var waitGroup sync.WaitGroup

	waitGroup.Add(2)

	go func() {
		defer waitGroup.Done()

		resposta1, erro1 = enviarParaServidor(
			"localhost:8001",
			parte1,
		)
	}()

	go func() {
		defer waitGroup.Done()

		resposta2, erro2 = enviarParaServidor(
			"localhost:8002",
			parte2,
		)
	}()

	waitGroup.Wait()

	if erro1 != nil {
		fmt.Println("Erro no Servidor 1:", erro1)
		return
	}

	if erro2 != nil {
		fmt.Println("Erro no Servidor 2:", erro2)
		return
	}

	totalPalavras := resposta1.Palavras + resposta2.Palavras
	totalLetras := resposta1.Letras + resposta2.Letras
	totalVogais := resposta1.Vogais + resposta2.Vogais
	totalNumeros := resposta1.Numeros + resposta2.Numeros

	fmt.Println("Palavras:", totalPalavras)
	fmt.Println("Letras:", totalLetras)
	fmt.Println("Vogais:", totalVogais)
	fmt.Println("Números:", totalNumeros)

	fmt.Println()
	fmt.Println("Servidor 1 processou:")
	fmt.Println(parte1)

	fmt.Println()
	fmt.Println("Servidor 2 processou:")
	fmt.Println(parte2)
}