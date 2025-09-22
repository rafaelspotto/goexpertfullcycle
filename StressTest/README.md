# StressTest - Sistema CLI para Testes de Carga

Sistema CLI em Go para realizar testes de carga em serviços web com relatórios detalhados.

## Funcionalidades

- ✅ Testes de carga com controle de concorrência
- ✅ Relatórios detalhados com estatísticas
- ✅ Suporte a HTTPS
- ✅ Containerização com Docker
- ✅ Interface CLI intuitiva

## Parâmetros

- `--url`: URL do serviço a ser testado (obrigatório)
- `--requests`: Número total de requests (obrigatório)
- `--concurrency`: Número de chamadas simultâneas (obrigatório)

## Uso Local

### Compilação
```bash
go build -o stresstest .
```

### Execução
```bash
./stresstest --url=http://google.com --requests=1000 --concurrency=10
```

## Uso com Docker

### Build da imagem
```bash
docker build -t stresstest .
```

### Execução
```bash
docker run stresstest --url=http://google.com --requests=1000 --concurrency=10
```

## Exemplos de Uso

### Teste básico
```bash
./stresstest --url=https://httpbin.org/get --requests=50 --concurrency=5
```

### Teste de alta carga
```bash
./stresstest --url=http://google.com --requests=1000 --concurrency=20
```

## Relatório Gerado

O sistema gera um relatório completo contendo:

- ⏱️ Tempo total de execução
- 📊 Total de requests realizados
- ✅ Requests com status HTTP 200
- 📈 Taxa de sucesso
- 🚀 Requests por segundo
- 📋 Distribuição de códigos de status HTTP
- ⏱️ Estatísticas de duração (mínima, máxima, média)
- ❌ Contagem de erros

## Estrutura do Projeto

```
StressTest/
├── main.go          # Código principal da aplicação
├── go.mod           # Dependências do Go
├── go.sum           # Checksums das dependências
├── Dockerfile       # Configuração do Docker
├── build.sh         # Script de build
└── README.md        # Este arquivo
```

## Tecnologias Utilizadas

- **Go 1.21**: Linguagem principal
- **Goroutines**: Para concorrência
- **Channels**: Para comunicação entre goroutines
- **HTTP Client**: Para requisições HTTP
- **Docker**: Para containerização
- **Alpine Linux**: Imagem base otimizada

## Características Técnicas

- **Controle de Concorrência**: Usa semáforos para limitar goroutines simultâneas
- **Timeout**: 30 segundos por requisição
- **Coleta de Métricas**: Em tempo real durante a execução
- **Relatório Detalhado**: Estatísticas completas de performance
- **Tratamento de Erros**: Captura e relata erros de conexão
