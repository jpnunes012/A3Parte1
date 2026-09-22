package main

import (
	"encoding/json"
	"fmt"
	"net"
)

type Request struct {
	Operacao string `json:"operacao"`
	Dados    []int  `json:"dados"`
}

type Response struct {
	Servidor  string `json:"servidor"`
	Resultado string `json:"resultado"`
}

func main() {
	listener, err := net.Listen("tcp", ":9001")
	if err != nil {
		fmt.Println("Erro ao iniciar servidor:", err)
		return
	}
	fmt.Println("Servidor escutando na porta 9001...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()

	var req Request
	if err := json.NewDecoder(conn).Decode(&req); err != nil {
		fmt.Println("Erro ao decodificar:", err)
		return
	}

	fmt.Println("Recebido:", req)

	resp := Response{
		Servidor:  "Servidor 9001",
		Resultado: fmt.Sprintf("Recebi %d números", len(req.Dados)),
	}

	json.NewEncoder(conn).Encode(resp)
}