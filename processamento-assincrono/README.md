# Processamento assíncrono

Exemplo simples de dois microsserviços processando a troca de informação de forma assíncrona.
Os microsserviços poderiam estar em repositórios separados, mas por uma questão de conveniência, estão neste mesmo monorepo.
Observe que este exemplo tem como foco o processamento assíncrono, não a melhor estruturação de um projeto para rodar em produção. Em outras palavras, existe código duplicado que poderia ser resolvido com uma lib comum, por exemplo.

Este exemplo faz parte da [série](https://dev.to/odilonjk/processamento-assincrono-parte-1-31pn) que apresento sobre processamento assíncrono no site [DEV.to](https://dev.to/).

### Estoque

Microsserviço que cria uma mensagem e envia para fila no RabbitMQ.
Também é o serviço responsável por criar as `exchanges` e `queues` no RabbitMQ.

### OrdemDeProducao

Microsserviço que consome mensagens de uma fila no RabbitMQ.
Neste serviço está concentrada a regra relacionada a retentativa das filas e tempo de delay entre o consumo das mensagens.

---

## Dependências

- Docker
- Go

---

## Como rodar o projeto?

É muito simples para rodar o exemplo fazendo uso do Makefile contido neste projeto.

O primeiro passo é criar a imagem do RabbitMQ com os plugins Management e Delayed Message ativos.

```shell
    make docker_build
```

Em seguida é necessário iniciar o RabbitMQ em seu Docker.

```shell
    make rabbit_start
```

Ao final dos seus testes, para desligar o RabbitMQ basta executar o comando abaixo:

```shell
    make rabbit_stop
```

### Estoque

Sempre que quiser adicionar uma nova mensagem na fila, basta rodar o exemplo que simula um microsserviço de Estoque.
Para isto você pode entrar no diretório `processamento-assincrono/estoque` e rodar `go run .` ou no diretório `processamento-assincrono` executar:

```shell
    make estoque_run
```

### OrdemDeProducao

Para iniciar o exemplo de microsserviço OrdemDeProducao, que irá consumir as mensagens no RabbitMQ, basta acessar o diretório `processamento-assincrono/ordem-de-producao` e executar `go run .`, ou no diretório `processamento-assincrono` executar o comando abaixo:

```shell
    make ordem_producao_start
```