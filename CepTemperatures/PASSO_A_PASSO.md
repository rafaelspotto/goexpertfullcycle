# Passo a Passo - Sistema de Consulta de Temperatura por CEP

## Visão Geral
Desenvolver um sistema em Go que recebe um CEP, identifica a cidade e retorna o clima atual em Celsius, Fahrenheit e Kelvin.

## Pré-requisitos
- Go 1.19+ instalado
- Conta no Google Cloud Platform
- Conta na WeatherAPI (https://www.weatherapi.com/)
- Editor de código (VS Code, GoLand, etc.)

## Passo 1: Configuração do Projeto
1. Criar estrutura de diretórios:
   ```
   CepTemperature/
   ├── cmd/
   │   └── server/
   │       └── main.go
   ├── internal/
   │   ├── handlers/
   │   ├── services/
   │   └── models/
   ├── pkg/
   │   └── weather/
   ├── go.mod
   ├── go.sum
   └── Dockerfile
   ```

2. Inicializar o módulo Go:
   ```bash
   go mod init cep-temperature
   ```

## Passo 2: Instalação de Dependências
```bash
go get github.com/gin-gonic/gin
go get github.com/joho/godotenv
go get github.com/stretchr/testify
```

## Passo 3: Criação dos Modelos de Dados
1. Criar `internal/models/weather.go`:
   - Struct para resposta da API de clima
   - Struct para resposta da API ViaCEP
   - Struct para resposta final do sistema

2. Criar `internal/models/response.go`:
   - Structs para diferentes tipos de resposta (sucesso, erro)

## Passo 4: Implementação dos Serviços
1. Criar `internal/services/cep_service.go`:
   - Função para validar formato do CEP (8 dígitos)
   - Função para consultar ViaCEP API
   - Tratamento de erros específicos

2. Criar `internal/services/weather_service.go`:
   - Função para consultar WeatherAPI
   - Função para converter temperaturas:
     - Celsius para Fahrenheit: F = C * 1.8 + 32
     - Celsius para Kelvin: K = C + 273

## Passo 5: Implementação dos Handlers
1. Criar `internal/handlers/weather_handler.go`:
   - Handler para endpoint GET /weather/:cep
   - Validação do CEP
   - Chamada dos serviços
   - Retorno das respostas padronizadas:
     - 200: `{"temp_C": 28.5, "temp_F": 83.3, "temp_K": 301.5}`
     - 422: `{"error": "invalid zipcode"}`
     - 404: `{"error": "can not find zipcode"}`

## Passo 6: Configuração do Servidor HTTP
1. Criar `cmd/server/main.go`:
   - Configuração do Gin router
   - Middleware de CORS
   - Middleware de logging
   - Registro das rotas
   - Configuração da porta (variável de ambiente)

## Passo 7: Configuração de Variáveis de Ambiente
1. Criar `.env`:
   ```
   PORT=8080
   WEATHER_API_KEY=sua_chave_da_weatherapi
   ```

2. Criar `.env.example`:
   ```
   PORT=8080
   WEATHER_API_KEY=your_weather_api_key
   ```

## Passo 8: Testes Unitários
1. Criar `internal/services/cep_service_test.go`:
   - Testes para validação de CEP
   - Testes para consulta da API ViaCEP

2. Criar `internal/services/weather_service_test.go`:
   - Testes para conversão de temperaturas
   - Testes para consulta da WeatherAPI

3. Criar `internal/handlers/weather_handler_test.go`:
   - Testes para diferentes cenários de resposta

## Passo 9: Dockerização
1. Criar `Dockerfile`:
   ```dockerfile
   FROM golang:1.21-alpine AS builder
   WORKDIR /app
   COPY go.mod go.sum ./
   RUN go mod download
   COPY . .
   RUN go build -o main cmd/server/main.go

   FROM alpine:latest
   RUN apk --no-cache add ca-certificates
   WORKDIR /root/
   COPY --from=builder /app/main .
   COPY --from=builder /app/.env .
   CMD ["./main"]
   ```

2. Criar `.dockerignore`:
   ```
   .git
   .gitignore
   README.md
   Dockerfile
   .dockerignore
   ```

## Passo 10: Configuração do Google Cloud Run
1. Criar `cloudbuild.yaml`:
   ```yaml
   steps:
     - name: 'gcr.io/cloud-builders/docker'
       args: ['build', '-t', 'gcr.io/$PROJECT_ID/cep-temperature', '.']
     - name: 'gcr.io/cloud-builders/docker'
       args: ['push', 'gcr.io/$PROJECT_ID/cep-temperature']
     - name: 'gcr.io/cloud-builders/gcloud'
       args: ['run', 'deploy', 'cep-temperature', '--image', 'gcr.io/$PROJECT_ID/cep-temperature', '--region', 'us-central1', '--platform', 'managed']
   ```

2. Configurar variáveis de ambiente no Cloud Run:
   - `WEATHER_API_KEY`: sua chave da WeatherAPI

## Passo 11: Deploy e Testes
1. Build e push da imagem:
   ```bash
   gcloud builds submit --config cloudbuild.yaml
   ```

2. Testes da API:
   ```bash
   # CEP válido
   curl https://your-cloud-run-url/weather/01310-100
   
   # CEP inválido
   curl https://your-cloud-run-url/weather/123
   
   # CEP não encontrado
   curl https://your-cloud-run-url/weather/99999-999
   ```

## Passo 12: Documentação da API
1. Criar `API.md` com:
   - Endpoints disponíveis
   - Exemplos de requisições
   - Exemplos de respostas
   - Códigos de erro

## Estrutura Final do Projeto
```
CepTemperature/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── handlers/
│   │   ├── weather_handler.go
│   │   └── weather_handler_test.go
│   ├── services/
│   │   ├── cep_service.go
│   │   ├── cep_service_test.go
│   │   ├── weather_service.go
│   │   └── weather_service_test.go
│   └── models/
│       ├── weather.go
│       └── response.go
├── pkg/
│   └── weather/
├── .env
├── .env.example
├── .dockerignore
├── .gitignore
├── Dockerfile
├── cloudbuild.yaml
├── go.mod
├── go.sum
├── README.md
├── PASSO_A_PASSO.md
└── API.md
```

## Comandos Úteis
```bash
# Executar localmente
go run cmd/server/main.go

# Executar testes
go test ./...

# Build para produção
go build -o bin/server cmd/server/main.go

# Executar com Docker localmente
docker build -t cep-temperature .
docker run -p 8080:8080 --env-file .env cep-temperature
```

## Próximos Passos
1. Implementar cache para consultas de CEP
2. Adicionar métricas e monitoramento
3. Implementar rate limiting
4. Adicionar logs estruturados
5. Implementar health check endpoint
