package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"unicode"
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

func analisarTexto(texto string) Resposta {
	resposta := Resposta{
		Servidor: "Servidor 2",
	}

	resposta.Palavras = len(strings.Fields(texto))

	for _, caractere := range texto {
		if unicode.IsLetter(caractere) {
			resposta.Letras++
		}

		if unicode.IsDigit(caractere) {
			resposta.Numeros++
		}

		if strings.ContainsRune("aeiouAEIOUáéíóúÁÉÍÓÚãõÃÕâêôÂÊÔ", caractere) {
			resposta.Vogais++
		}
	}

	return resposta
}

func tratarConexao(conexao net.Conn) {
	defer conexao.Close()

	leitor := bufio.NewReader(conexao)

	mensagem, erro := leitor.ReadString('\n')
	if erro != nil {
		return
	}

	var requisicao Requisicao

	erro = json.Unmarshal([]byte(mensagem), &requisicao)
	if erro != nil {
		return
	}

	resposta := analisarTexto(requisicao.Texto)

	jsonResposta, erro := json.Marshal(resposta)
	if erro != nil {
		return
	}

	fmt.Fprintln(conexao, string(jsonResposta))
}

func main() {
	listener, erro := net.Listen("tcp", ":8002")
	if erro != nil {
		fmt.Println("Erro ao iniciar servidor:", erro)
		return
	}

	defer listener.Close()
	
	fmt.Println("Servidor 2 iniciado na porta 8002...")
	for {
		conexao, erro := listener.Accept()

		if erro != nil {
			continue
		}
		go tratarConexao(conexao)
	}
}