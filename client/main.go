package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
    "time"
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

func salvarLog(
	textoCompleto string,
	parte1 string,
	parte2 string,
	resposta1 Resposta,
	resposta2 Resposta,
) error {

	log := Log{
		DataHora:      time.Now().Format(time.RFC3339),
		TextoCompleto: textoCompleto,

		Partes: []ParteLog{
			{
				Servidor: resposta1.Servidor,
				Texto:    parte1,
				Palavras: resposta1.Palavras,
				Vogais:   resposta1.Vogais,
				Letras:   resposta1.Letras,
				Numeros:  resposta1.Numeros,
			},
			{
				Servidor: resposta2.Servidor,
				Texto:    parte2,
				Palavras: resposta2.Palavras,
				Vogais:   resposta2.Vogais,
				Letras:   resposta2.Letras,
				Numeros:  resposta2.Numeros,
			},
		},

		Total: TotalLog{
			Palavras: resposta1.Palavras + resposta2.Palavras,
			Vogais:   resposta1.Vogais + resposta2.Vogais,
			Letras:   resposta1.Letras + resposta2.Letras,
			Numeros:  resposta1.Numeros + resposta2.Numeros,
		},
	}

	dadosJSON, erro := json.MarshalIndent(log, "", "    ")

	if erro != nil {
		return erro
	}

	erro = os.MkdirAll("logs", 0755)

	if erro != nil {
		return erro
	}

	nomeArquivo := fmt.Sprintf(
		"logs/log_%s.json",
		time.Now().Format("2006-01-02_15-04-05"),
	)

	erro = os.WriteFile(nomeArquivo, dadosJSON, 0644)

	if erro != nil {
		return erro
	}

	fmt.Println("Log salvo em:", nomeArquivo)

	return nil
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

    fmt.Printf(
        "\n========== Resultado ==========\n"+
            "Palavras: %d | Letras: %d | Vogais: %d | Números: %d\n\n"+
            "Servidor 1: %s\n"+
            "Servidor 2: %s\n",
        totalPalavras,
        totalLetras,
        totalVogais,
        totalNumeros,
        parte1,
        parte2,
    )

    erro = salvarLog(
	    texto,
	    parte1,
	    parte2,
	    resposta1,
	    resposta2,
    )

    if erro != nil {
	    fmt.Println("Erro ao salvar log:", erro)
    }
}