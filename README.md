# Avaliação 3, atividade 1 - Prof. Saulo Popov

## Sistemas distribuídos e mobile

Aplicação distribuída desenvolvida em **Go** para demonstrar conceitos de comunicação entre processos, sockets TCP, comunicação HTTP, concorrência, goroutines e distribuição de processamento.

O sistema recebe um texto informado pelo usuário através de uma interface web, divide o texto em duas partes e distribui o processamento entre dois servidores. Cada servidor analisa sua parte do texto e retorna os resultados ao cliente, que reúne as informações e apresenta o resultado final.

Após cada processamento, o sistema gera um arquivo `.json` contendo os dados da execução, permitindo manter um histórico das operações realizadas.

## Objetivo

O objetivo do projeto é apresentar uma aplicação distribuída utilizando uma arquitetura **cliente-servidor**, na qual o processamento de uma tarefa é distribuído entre diferentes processos.

O cliente é responsável por:

* receber o texto através da interface web;
* dividir o texto em duas partes;
* enviar cada parte para um servidor diferente;
* aguardar as respostas dos servidores;
* reunir os resultados;
* apresentar os dados ao usuário;
* gerar um arquivo JSON com o histórico da execução.

Os servidores recebem suas respectivas partes do texto através de **sockets TCP**, realizam a análise e retornam os resultados utilizando mensagens no formato **JSON**.

## Arquitetura

O projeto utiliza três processos principais:

* **Cliente:** disponibiliza a interface web, recebe o texto, divide o conteúdo, envia as partes para os servidores e reúne os resultados.
* **Servidor 1:** recebe e processa a primeira parte do texto através da porta `8001`.
* **Servidor 2:** recebe e processa a segunda parte do texto através da porta `8002`.

A aplicação também utiliza um servidor HTTP no cliente, executado na porta `8080`, responsável por disponibilizar a interface web e a API de processamento.

```text
                         +----------------------+
                         |       NAVEGADOR      |
                         |   Interface Web      |
                         +----------+-----------+
                                    |
                              HTTP :8080
                                    |
                                    v
                         +----------------------+
                         |       CLIENTE        |
                         |      Go / HTTP       |
                         +----------+-----------+
                                    |
                              Divide o texto
                           +--------+--------+
                           |                 |
                         TCP :8001         TCP :8002
                           |                 |
                           v                 v
                  +----------------+  +----------------+
                  |   SERVIDOR 1   |  |   SERVIDOR 2   |
                  |    Go / TCP    |  |    Go / TCP    |
                  +----------------+  +----------------+
                           |                 |
                           |   Processamento |
                           +--------+--------+
                                    |
                                    v
                         +----------------------+
                         |       CLIENTE        |
                         |  Junta os resultados |
                         +----------+-----------+
                                    |
                                    v
                         +----------------------+
                         |       NAVEGADOR      |
                         | Resultado da análise |
                         +----------------------+
                                    |
                                    v
                         +----------------------+
                         |      logs/*.json     |
                         | Histórico da execução|
                         +----------------------+
```

## Funcionalidades

Cada servidor realiza uma análise da parte do texto recebida.

Atualmente são contabilizados:

* quantidade de palavras;
* quantidade de letras;
* quantidade de vogais;
* quantidade de números.

### Contagem de palavras

As palavras são identificadas utilizando `strings.Fields`, considerando os espaços como separadores.

### Contagem de letras

As letras são identificadas utilizando `unicode.IsLetter`, permitindo reconhecer também caracteres acentuados.

### Contagem de números

Os números são identificados utilizando `unicode.IsDigit`.

### Contagem de vogais

O sistema reconhece vogais maiúsculas, minúsculas e algumas vogais acentuadas, incluindo:

```text
a e i o u
á é í ó ú
ã õ
â ê ô
```

Ao final, o cliente soma os resultados dos dois servidores e apresenta a análise completa do texto.

## Comunicação

A comunicação entre o cliente e os servidores é realizada utilizando **sockets TCP**.

Os dados enviados entre os processos são serializados utilizando o formato **JSON**.

O cliente se conecta aos servidores através dos seguintes endereços:

```text
localhost:8001
localhost:8002
```

### Exemplo de requisição

O cliente envia uma requisição semelhante a:

```json
{
    "texto": "Programação distribuída permite dividir tarefas"
}
```

### Exemplo de resposta

Um servidor retorna uma resposta semelhante a:

```json
{
    "servidor": "Servidor 1",
    "palavras": 5,
    "vogais": 17,
    "letras": 42,
    "numeros": 0
}
```

## Divisão do texto

O cliente divide o texto em duas partes utilizando a quantidade de palavras.

O processo é realizado pela função:

```go
func dividirTexto(texto string) (string, string)
```

Primeiramente, o texto é separado em palavras utilizando:

```go
strings.Fields(texto)
```

Depois, é calculado o ponto central:

```go
meio := len(palavras) / 2
```

A primeira metade é enviada ao **Servidor 1** e a segunda metade é enviada ao **Servidor 2**.

Por exemplo:

```text
Texto original:

Programação distribuída permite utilizar vários
computadores para processar uma tarefa 123.
```

Pode ser dividido em:

```text
Servidor 1:
Programação distribuída permite utilizar vários

Servidor 2:
computadores para processar uma tarefa 123.
```

## Concorrência

O projeto utiliza **goroutines** para permitir que as duas partes do texto sejam enviadas aos servidores de maneira concorrente.

No cliente são criadas duas goroutines:

```go
go func() {
    defer waitGroup.Done()
    resposta1, erro1 = enviarServidor("localhost:8001", parte1)
}()

go func() {
    defer waitGroup.Done()
    resposta2, erro2 = enviarServidor("localhost:8002", parte2)
}()
```

Dessa maneira, o cliente pode realizar as duas comunicações simultaneamente.

Para garantir que o processamento das duas goroutines seja concluído antes da apresentação do resultado, é utilizado `sync.WaitGroup`:

```go
waitGroup.Add(2)

...

waitGroup.Wait()
```

### Concorrência nos servidores

Os dois servidores também utilizam goroutines para atender diferentes conexões.

Quando uma conexão é aceita:

```go
conexao, erro := listener.Accept()
```

ela é encaminhada para uma goroutine:

```go
go tratarConexao(conexao)
```

Isso permite que o servidor continue aceitando novas conexões enquanto outras estão sendo processadas.

## Interface Web

O cliente disponibiliza uma interface web através de um servidor HTTP.

A aplicação utiliza a porta:

```text
localhost:8080
```

A interface permite:

1. inserir ou colar um texto;
2. clicar no botão **Dividir**;
3. enviar o texto para o cliente;
4. visualizar a divisão realizada;
5. visualizar os resultados de cada servidor;
6. visualizar os totais da análise.

A comunicação entre a interface e o cliente é realizada através de uma requisição HTTP `POST` para:

```text
/api/processar
```

Os dados enviados para a API utilizam JSON.

## Logs

Após o processamento, o cliente gera automaticamente um arquivo `.json` contendo informações sobre a execução.

Os arquivos são armazenados em:

```text
client/logs/
```

### Exemplo

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

O nome do arquivo é gerado automaticamente utilizando a data e hora da execução:

```text
logs/log_2026-09-22_20-47-32.json
```

## Estrutura do projeto

```text
A3Parte1/
│
├── go.mod
├── README.md
│
├── client/
│   ├── interface/
│   │   └── index.html
│   │
│   ├── logs/
│   │
│   └── main.go
│
├── server1/
│   └── s1.go
│
└── server2/
    └── s2.go
```

### Descrição dos arquivos

| Arquivo/Pasta                 | Função                                                                                    |
| ----------------------------- | ----------------------------------------------------------------------------------------- |
| `client/main.go`              | Implementa o cliente, servidor HTTP, divisão do texto, comunicação TCP e geração dos logs |
| `client/interface/index.html` | Interface web utilizada pelo usuário                                                      |
| `client/logs/`                | Armazena os arquivos JSON gerados após cada processamento                                 |
| `server1/s1.go`               | Implementa o Servidor 1                                                                   |
| `server2/s2.go`               | Implementa o Servidor 2                                                                   |
| `go.mod`                      | Arquivo de configuração do módulo Go                                                      |
| `README.md`                   | Documentação do projeto                                                                   |

## Tecnologias utilizadas

* **Go**
* **HTML**
* **JavaScript**
* **HTTP**
* **TCP Sockets**
* **JSON**
* **Goroutines**
* **sync.WaitGroup**
* **Programação concorrente**
* **Programação distribuída**

## Como executar

É necessário possuir o **Go** instalado na máquina.

O download pode ser realizado através do site oficial:

https://go.dev/dl/

Após instalar o Go, abra um terminal na pasta principal do projeto.

### 1. Iniciar o Servidor 1

Execute:

```bash
go run s1.go
```

O servidor ficará disponível na porta:

```text
localhost:8001
```

A mensagem apresentada será semelhante a:

```text
Servidor 1 iniciado na porta 8001...
```

### 2. Iniciar o Servidor 2

Abra outro terminal na pasta principal do projeto e execute:

```bash
go run s2.go
```

O servidor ficará disponível na porta:

```text
localhost:8002
```

A mensagem apresentada será semelhante a:

```text
Servidor 2 iniciado na porta 8002...
```

### 3. Iniciar o cliente

Abra um terceiro terminal na pasta principal do projeto e execute:

```bash
go run main.go
```

O cliente iniciará o servidor HTTP na porta `8080`.

Será apresentada uma mensagem assim:

```text
Aberto em http://localhost:8080
Os servidores de localhost:8001 e localhost:8002 precisam estar de pé.
```

### 4. Abrir a interface

Abra um navegador e acesse:

```text
http://localhost:8080
```

### 5. Processar um texto

Digite ou cole um texto na interface e clique no botão:

```text
Dividir
```

Por exemplo:

```text
Programação distribuída permite utilizar vários computadores para processar uma tarefa 123.
```

O cliente dividirá o texto e enviará cada parte para um servidor.

Após o processamento, serão apresentados os resultados individuais e o resultado total.

Exemplo:

```text
Servidor 1:
Programação distribuída permite utilizar vários

Palavras: 5 | Letras: 45 | Vogais: 21 | Números: 0

Servidor 2:
computadores para processar uma tarefa 123.

Palavras: 6 | Letras: 34 | Vogais: 15 | Números: 3

TOTAL

Palavras: 11
Letras: 79
Vogais: 36
Números: 3
```

Também será criado um arquivo de log dentro da pasta:

```text
client/logs/
```

## Funcionamento resumido

O funcionamento da aplicação pode ser dividido nas seguintes etapas:

1. O usuário acessa a interface web.
2. O usuário informa um texto.
3. A interface envia o texto para a API HTTP do cliente.
4. O cliente valida o texto recebido.
5. O cliente divide o texto em duas partes.
6. A primeira parte é enviada ao Servidor 1 através de TCP.
7. A segunda parte é enviada ao Servidor 2 através de TCP.
8. As duas comunicações são realizadas utilizando goroutines.
9. Cada servidor analisa sua respectiva parte do texto.
10. Os servidores retornam os resultados em JSON.
11. O cliente aguarda as duas respostas utilizando `WaitGroup`.
12. O cliente soma os resultados.
13. Os resultados são enviados para a interface web.
14. A interface apresenta os resultados ao usuário.
15. O cliente salva um arquivo JSON contendo o histórico da operação.

## Conceitos aplicados

Durante o desenvolvimento do projeto foram aplicados conceitos relacionados a:

* arquitetura cliente-servidor;
* processos distribuídos;
* comunicação entre processos;
* sockets TCP;
* servidor HTTP;
* API HTTP;
* serialização e desserialização JSON;
* programação concorrente;
* goroutines;
* sincronização utilizando `WaitGroup`;
* distribuição de processamento;
* divisão de tarefas;
* agregação de resultados;
* tratamento de conexões;
* armazenamento de logs;
* desenvolvimento de interface web.

## Autores

<div align="center">

| Foto                                                                     | Integrante              |
| ------------------------------------------------------------------------ | ----------------------- |
| <img src="https://avatars.githubusercontent.com/GabHrq" width="120">     | Gabriel Henrique Passos |
| <img src="https://avatars.githubusercontent.com/jpnunes012" width="120"> | João Pedro Nunes        |
| <img src="https://avatars.githubusercontent.com/Diblewy" width="120">    | Juan Mello Arjona       |

</div>
