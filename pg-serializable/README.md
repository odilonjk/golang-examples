#PG Serializable Isolation Level

Exemplo de uma das formas de garantir a consistência dos dados.
Neste exemplo é utilizado o lock na tabela utilizando o nível de isolamento `serializable`.

Caso você tenha interesse em se aprofundar no funcionamento do isolamento das transações no PostgreSQL, indico a leitura do artigo [Deeply understand Isolation levels and Read phenomena in MySQL & PostgreSQL](https://dev.to/techschoolguru/understand-isolation-levels-read-phenomena-in-mysql-postgres-c2e#isolation-levels-in-postgres).

A própria [documentação](http://pgdocptbr.sourceforge.net/pg80/transaction-iso.html) do PostgreSQL também é bastante rica em explicações de como funcionam os nível de isolamento.

### Problema
A API recebe chamadas para realizar reservas. Porém a regra de negócio não permite reservas com datas sobrepostas.
Pode haver N instâncias da aplicação rodando, portanto pode acontecer de 2 ou mais instâncias tentarem persistir a reserva no banco ao mesmo tempo.

### Solução
A aplicação quando abre a transação para o banco, utiliza o nível de isolamento `serializable`. Desta forma, é como se todas transações executassem de forma sequencial. Caso tentem executar ao mesmo tempo como o problema descrito acima, apenas uma transação é executada com sucesso e as demais tomam erro realizando assim o rollback da transação.

### Vantagens

- Garante que não haverá registros com datas de sobreponto, respeitando a regra de negócio.

### Desvantagens

- Pode acabar transformando as transações no banco em um gargalo, prejudicando a performance.

### Como rodar o exemplo?

> **_NOTE:_**  O volume do banco está mapeado para `$HOME/workspace/volumes/postgres`. Garanta que este caminho existe ou altere da forma que desejar no arquivo _docker-compose.yml_ 

- Caso seja a primeira vez que você irá rodar a aplicação, é necessário realizar o build do binário com o comando `make build` e em seguida gerar a imagem no Docker com o comando `make build_docker`.

- Para iniciar as instâncias da aplicação e o PostgreSQL, execute `make start`.

- Para parar as instâncias da aplicação e o PostgreSQL, executa `make stop`.

- Caso queira realizar uma requisição manualmente, basta executar o seguinte cURL: 
    ```
    curl -XPOST -d 'start_date=2021-01-05&end_date=2021-01-07' 'http://localhost:8080/bookings'
    ```

### Como testar?

> **_NOTE:_** A configuração do nível de isolamento está no arquivo _service.go_. Você pode alterar o nível para validar a diferença na prática. Por exemplo, tente alterar para _read committed_

- Com as instâncias e o PostgreSQL de pé, execute `make concurrent_calls` para que sejam realizadas as chamadas.

- Você pode acessar o banco e verificar se foi gravado apenas um registro (garantindo a consistência dos dados) ou se existem múltiplos registros (onde haveria inconsistência de dados).

### PostgreSQL

**Database:** booking

**User:** booking_app

**Pass:** pg