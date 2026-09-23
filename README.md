# Avaliação 3, atividade 1 - Prof. Saulo Popov
## Sistemas distribuídos e mobile

Aplicação distribuída desenvolvida em Go para demonstrar conceitos de comunicação entre processos, sockets, concorrência e distribuição de processamento.

O sistema recebe um texto no cliente, divide esse texto em partes e distribui o processamento entre dois servidores. Cada servidor analisa sua parte do texto e retorna os resultados ao cliente, que reúne as informações e apresenta o resultado final.

Cada execução gera um arquivo `.json` contendo os dados processados, permitindo manter um histórico das operações realizadas.

## Objetivo

Apresentação simplificada de uma aplicação distribuída utilizando uma arquitetura cliente-servidor: O cliente é responsável por receber o texto informado pelo usuário e dividir o conteúdo entre os servidores disponíveis, e os servidores recebem suas respectivas partes por meio de sockets TCP, realizam o processamento e retornam os resultados utilizando mensagens no formato JSON.

## Arquitetura

O projeto utiliza três processos principais:

- **Cliente:** recebe o texto do usuário, divide o conteúdo, envia as partes para os servidores e reúne os resultados.
- **Servidor 1:** recebe e processa a primeira parte do texto.
- **Servidor 2:** recebe e processa a segunda parte do texto.

Fluxo simplificado:

```text
                 +----------------+
                 |     CLIENTE    |
                 +----------------+
                         |
                  Divide o texto
                         |
                +--------+--------+
                |                 |
                v                 v
        +---------------+ +---------------+
        |   SERVIDOR 1  | |   SERVIDOR 2  |
        |   Porta 8001  | |   Porta 8002  |
        +---------------+ +---------------+
                |                 |
                |   Processamento |
                |                 |
                +--------+--------+
                         |
                         v
                 +-----------------+
                 |     CLIENTE     |
                 | Junta resultados|
                 +-----------------+
```

## Funcionalidades

Cada servidor realiza uma análise da parte do texto recebida.

Atualmente são contabilizados:

- quantidade de palavras;
- quantidade de letras;
- quantidade de vogais;
- quantidade de números.

Ao final, o cliente soma os resultados dos dois servidores e apresenta a análise completa do texto.

## Comunicação

A comunicação entre cliente e servidores é realizada utilizando **sockets TCP**.
Os dados são enviados utilizando o formato **JSON**.

Exemplo de requisição enviada pelo cliente:

```json
{
    "texto": "Programação distribuída permite dividir tarefas"
}
```

Exemplo de resposta de um servidor:

```json
{
    "servidor": "Servidor 1",
    "palavras": 5,
    "vogais": 17,
    "letras": 42,
    "numeros": 0
}
```

## Concorrência

O projeto utiliza **goroutines** para permitir processamento concorrente.
No cliente, duas goroutines são utilizadas para enviar as partes do texto simultaneamente para os dois servidores.

```go
go func() {
    resposta1, erro1 = enviarParaServidor("localhost:8001", parte1)
}()

go func() {
    resposta2, erro2 = enviarParaServidor("localhost:8002", parte2)
}()
```

Também é utilizado `sync.WaitGroup` para garantir que o cliente espere a resposta dos dois servidores antes de calcular o resultado final.
Os servidores também utilizam goroutines para tratar diferentes conexões de maneira concorrente.


## Logs

Após o processamento, o cliente gera automaticamente um arquivo `.json` contendo informações sobre a execução.

Exemplo:

```json
{
    "dataHora": "2026-09-22T20:45:31-03:00",
    "textoCompleto": "Programação distribuída permite utilizar 2 computadores.",
    "partes": [
        {
            "servidor": "Servidor 1",
            "texto": "Programação distribuída permite",
            "palavras": 3,
            "vogais": 13,
            "letras": 30,
            "numeros": 0
        },
        {
            "servidor": "Servidor 2",
            "texto": "utilizar 2 computadores.",
            "palavras": 3,
            "vogais": 8,
            "letras": 19,
            "numeros": 1
        }
    ],
    "total": {
        "palavras": 6,
        "vogais": 21,
        "letras": 49,
        "numeros": 1
    }
}
```

Os arquivos são armazenados na pasta:

```text
client/logs/
```

## Estrutura do projeto

```text
a3/
│
├── go.mod
│
├── client/
│   ├── logs/
│   └── main.go
│
├── server1/
│   └── s1.go
│
├── server2/
│   └── s2.go
│
└── README.md
```

## Tecnologias utilizadas

- Go
- TCP Sockets
- JSON
- Goroutines
- WaitGroup
- Programação concorrente
- Programação distribuída

## Como executar

Necessário possuir [Go](https://go.dev/dl/) em sua máquina. Abra o link, clique para baixar a versão compatível com seu sistema operacional e siga as instruções do instalador.

Após instalado, abra um terminal na pasta principal do projeto.

### 1. Iniciar o Servidor 1

```bash
go run ./server1
```

O servidor ficará disponível em:

```text
localhost:8001
```

### 2. Iniciar o Servidor 2

Abra outro terminal:

```bash
go run ./server2
```

O servidor ficará disponível em:

```text
localhost:8002
```

### 3. Iniciar o cliente

Em um terceiro terminal:

```bash
go run ./client
```

O programa solicitará um texto:

```text
Digite um texto:
```

Exemplo:

```text
Programação distribuída permite utilizar vários computadores para processar uma tarefa 123.
```

O texto será dividido entre os dois servidores, e dará o resultado:

```text
========== Resultado ==========

Palavras: 11 | Letras: 79 | Vogais: 36 | Números: 3

Servidor 1:
Programação distribuída permite utilizar vários

Servidor 2:
computadores para processar uma tarefa 123.

Log salvo em: logs/log_2026-09-22_20-47-32.json
```

---

## Conceitos aplicados

Durante o desenvolvimento do TextSplit foram aplicados conceitos relacionados a:

- processos cliente e servidor;
- comunicação entre processos;
- sockets TCP;
- serialização e desserialização JSON;
- programação concorrente;
- goroutines;
- sincronização utilizando `WaitGroup`;
- distribuição de processamento;
- agregação de resultados;
- armazenamento de logs.

## Funcionamento resumido

1. O usuário informa um texto.
2. O cliente divide o texto em duas partes.
3. A primeira parte é enviada ao Servidor 1.
4. A segunda parte é enviada ao Servidor 2.
5. Os servidores processam suas partes simultaneamente.
6. Cada servidor retorna um JSON com os resultados.
7. O cliente reúne as duas respostas.
8. O resultado completo é apresentado ao usuário.
9. Um arquivo JSON contendo o histórico da operação é criado.

## Autores

<div align="center">

  | Foto | Integrante |
  |-------|------------|
  | <img src="https://avatars.githubusercontent.com/GabHrq" width="120"> | Gabriel Henrique Passos |
  | <img src="https://avatars.githubusercontent.com/jpnunes012" width="120"> | João Pedro Nunes |
  | <img src="https://avatars.githubusercontent.com/Diblewy" width="120"> | Juan Mello Arjona |

</div>
