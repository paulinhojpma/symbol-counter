# Symbol Counter

Symbol Counter é um package em Golang projetado para realizar a contagem de símbolos recebidos em fluxos de texto, onde cada símbolo deve estar separado por `\n`.

O projeto foi desenvolvido para suportar cenários de altíssima concorrência, sendo otimizado para processar milhões de requisições por segundo mantendo baixa latência e alta disponibilidade.

Cada conteúdo recebido é analisado e todos os símbolos únicos são contabilizados individualmente.

## Exemplo

Entrada:

```text
AU_ieu13
103956
ghgqb
10002
a012ne
AU_ieu13
```

Resultado:

```json
{
  "AU_ieu13": 2,
  "103956": 1,
  "ghgqb": 1,
  "10002": 1,
  "a012ne": 1
}
```

## Funcionalidades

- Processamento assíncrono
- Suporte para milhões de requisições por segundo
- Sistema de subscribers
- Agregação concorrente
- Graceful shutdown
- Redução de contenção por mutex

## API

### Analyse

```go
Analyse(context.Context, []byte)
```

Responsável por processar o conteúdo recebido e iniciar a contagem dos símbolos.

### GetCurrentCounts

```go
GetCurrentCounts() map[string]uint64
```

Retorna os valores atuais sem reiniciar os contadores.

### Subscribe

```go
Subscribe(
    writer io.WriteCloser,
    interval time.Duration,
) error
```

Adiciona um assinante para recebimento periódico dos relatórios.

### Shutdown

```go
Shutdown(ctx context.Context) error
```

Executa o encerramento correto do sistema.

# Arquitetura

```text
Requisições
      |
      v
Channel com Buffer
      |
      v
Worker
      |
      v
Goroutines de Contagem
      |
      v
Maps Concorrentes
      |
      v
Agregador
      |
      v
Mapa Principal
      |
      +--> Subscribers
```

## Estratégias de Performance

### 1. Channel com Buffer

As requisições são colocadas inicialmente em um channel com buffer que atua como fila e absorve picos de tráfego.

### 2. Processamento Assíncrono

Uma goroutine principal recebe os dados e cria goroutines dedicadas para executar a soma dos símbolos.

### 3. Distribuição de Locks

Para evitar gargalos em um único `map[string]uint64`, o sistema cria múltiplos mapas:

```go
[]struct{
    mutex sync.Mutex
    data map[string]uint64
}
```

Cada mapa possui seu próprio mutex e recebe parte da carga aleatoriamente.

Exemplo:

```text
Requisição A -> Map 1
Requisição B -> Map 4
Requisição C -> Map 2
```

## Processo de Agregação

Os dados dos mapas concorrentes são drenados:

- Automaticamente a cada segundo
- Quando `GetCurrentCounts()` é chamado

Fluxo:

```text
Maps Concorrentes
      |
      v
Drain
      |
      v
Mapa Principal
```

## Sistema de Subscribers

Uma goroutine dedicada envia relatórios periodicamente sem impactar a ingestão.

## Graceful Shutdown

Garantias:

- Finalização da fila
- Encerramento dos workers
- Flush dos subscribers
- Agregação final
- Liberação de recursos

## Resumo

| Componente | Responsabilidade |
|------------|------------------|
| Channel Buffer | Fila |
| Worker | Consumo |
| Goroutines | Contagem |
| Maps Concorrentes | Redução de lock |
| Agregador | Consolidação |
| Subscribers | Relatórios |
| Mapa Principal | Resultado |
