# Rate Limiter

Um rate limiter em Go que pode ser configurado para limitar o número máximo de requisições por segundo com base em um endereço IP específico ou em um token de acesso.

## Características

- **Limitação por IP**: Restringe o número de requisições recebidas de um único endereço IP
- **Limitação por Token**: Permite diferentes limites para diferentes tokens de acesso
- **Configuração flexível**: Configuração via variáveis de ambiente ou arquivo `.env`
- **Persistência Redis**: Armazena dados de limitação no Redis com fallback para memória
- **Middleware HTTP**: Fácil integração com servidores web
- **Estratégia de armazenamento**: Interface que permite trocar facilmente o Redis por outro mecanismo
- **Testes automatizados**: Cobertura completa de testes unitários e de integração

## Arquitetura

O projeto está organizado da seguinte forma:

```
├── cmd/                    # Ponto de entrada da aplicação
├── internal/
│   ├── config/            # Configuração e carregamento de variáveis de ambiente
│   ├── limiter/           # Lógica principal do rate limiter
│   ├── middleware/        # Middleware HTTP para Gin
│   ├── server/            # Servidor HTTP
│   └── storage/           # Interfaces e implementações de armazenamento
├── test/                  # Testes automatizados
├── scripts/               # Scripts de teste e load testing
├── docker-compose.yml     # Configuração Docker para Redis e aplicação
├── Dockerfile            # Imagem Docker da aplicação
└── config.env            # Arquivo de configuração de exemplo
```

## Como Funciona

### Limitação por IP
- Cada IP tem um limite configurável de requisições por segundo
- Quando o limite é excedido, o IP é bloqueado por um período configurável
- O contador de requisições é resetado a cada segundo

### Limitação por Token
- Tokens específicos podem ter limites diferentes dos IPs
- Tokens são identificados pelo header `API_KEY`
- Se um token não for reconhecido, o sistema usa os limites do IP
- Tokens têm seus próprios períodos de bloqueio

### Estratégia de Armazenamento
- **Redis**: Armazenamento principal com persistência
- **Memória**: Fallback para testes ou quando Redis não está disponível
- Interface `Storage` permite fácil substituição por outros mecanismos

## Configuração

### Variáveis de Ambiente

```bash
# Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# Rate Limiter Settings
RATE_LIMIT_IP_REQUESTS_PER_SECOND=5
RATE_LIMIT_IP_BLOCK_DURATION_MINUTES=5

# Token Rate Limits (format: TOKEN_LIMIT_<TOKEN>=<REQUESTS_PER_SECOND>:<BLOCK_DURATION_MINUTES>)
TOKEN_LIMIT_abc123=10:5
TOKEN_LIMIT_def456=20:10
TOKEN_LIMIT_ghi789=50:15

# Server Configuration
SERVER_PORT=8080
```

### Exemplo de Configuração

```bash
# IP pode fazer 5 requisições por segundo, bloqueado por 5 minutos se exceder
RATE_LIMIT_IP_REQUESTS_PER_SECOND=5
RATE_LIMIT_IP_BLOCK_DURATION_MINUTES=5

# Token abc123 pode fazer 10 requisições por segundo, bloqueado por 5 minutos
TOKEN_LIMIT_abc123=10:5

# Token def456 pode fazer 20 requisições por segundo, bloqueado por 10 minutos
TOKEN_LIMIT_def456=20:10
```

## Uso

### Executando com Docker Compose

```bash
# Iniciar Redis e aplicação
docker-compose up -d

# Verificar logs
docker-compose logs -f

# Parar serviços
docker-compose down
```

### Executando Localmente

```bash
# Instalar dependências
go mod download

# Executar aplicação (Redis deve estar rodando)
go run cmd/main.go

# Ou executar com configuração personalizada
REDIS_HOST=localhost go run cmd/main.go
```

## Endpoints da API

### `GET /health`
Verificação de saúde da aplicação (sem rate limiting).

**Resposta:**
```json
{
  "status": "ok",
  "time": "2024-01-01T12:00:00Z"
}
```

### `GET /api/test`
Endpoint de teste com rate limiting aplicado.

**Headers:**
- `API_KEY` (opcional): Token de acesso

**Resposta (200):**
```json
{
  "message": "Request successful",
  "ip": "192.168.1.1",
  "token": "abc123",
  "time": "2024-01-01T12:00:00Z"
}
```

**Resposta (429 - Rate Limited):**
```json
{
  "error": "you have reached the maximum number of requests or actions allowed within a certain time frame",
  "reason": "Rate limit exceeded: 5 requests per second"
}
```

### `GET /api/status`
Verifica o status atual do rate limiter para o IP/token atual.

**Headers:**
- `API_KEY` (opcional): Token de acesso

**Resposta:**
```json
{
  "remaining_requests": 3,
  "blocked": false,
  "ip": "192.168.1.1",
  "token": "abc123"
}
```

### `POST /admin/unblock`
Desbloqueia um IP ou token (endpoint administrativo).

**Body:**
```json
{
  "type": "ip",
  "key": "192.168.1.1"
}
```

ou

```json
{
  "type": "token",
  "key": "abc123"
}
```

## Testes

### Executando Testes Unitários

```bash
go test -v ./test/...
```

### Executando Testes de Integração

```bash
# Iniciar Redis
docker-compose up -d redis

# Executar testes
go test -v -tags=integration ./test/...

# Limpar
docker-compose down
```

### Executando Script de Teste Completo

```bash
./scripts/test.sh
```

### Executando Load Test

```bash
./scripts/load_test.sh
```

## Exemplos de Uso

### Exemplo 1: Limitação por IP

```bash
# Fazer 6 requisições rapidamente (limite é 5 por segundo)
for i in {1..6}; do
  curl -w "Request $i: HTTP %{http_code}\n" http://localhost:8080/api/test
done
```

Resultado esperado:
- Requisições 1-5: HTTP 200
- Requisição 6: HTTP 429

### Exemplo 2: Limitação por Token

```bash
# Fazer 11 requisições com token abc123 (limite é 10 por segundo)
for i in {1..11}; do
  curl -w "Request $i: HTTP %{http_code}\n" -H "API_KEY: abc123" http://localhost:8080/api/test
done
```

Resultado esperado:
- Requisições 1-10: HTTP 200
- Requisição 11: HTTP 429

### Exemplo 3: Token Desconhecido

```bash
# Token desconhecido usa limite do IP
curl -H "API_KEY: unknown-token" http://localhost:8080/api/test
```

## Desenvolvimento

### Adicionando Nova Estratégia de Armazenamento

1. Implemente a interface `Storage` em `internal/storage/`
2. Adicione lógica de seleção em `internal/server/server.go`

```go
// Exemplo: implementar armazenamento em arquivo
type FileStorage struct {
    // implementação
}

func (f *FileStorage) GetRequestCount(ctx context.Context, key string) (int, error) {
    // implementação
}
```

### Adicionando Novos Tipos de Limitação

1. Estenda a configuração em `internal/config/config.go`
2. Atualize a lógica em `internal/limiter/limiter.go`
3. Adicione testes correspondentes

## Monitoramento

### Headers de Resposta

- `X-RateLimit-Remaining`: Número de requisições restantes no período atual

### Logs

A aplicação registra:
- Início e parada do servidor
- Erros de conexão com Redis
- Fallback para armazenamento em memória

## Troubleshooting

### Redis não Conecta

A aplicação automaticamente usa armazenamento em memória se Redis não estiver disponível. Verifique:
- Redis está rodando
- Configurações de host/porta estão corretas
- Firewall não está bloqueando a conexão

### Rate Limiter não Funciona

Verifique:
- Configurações de rate limit estão corretas
- Headers `API_KEY` estão sendo enviados corretamente
- IP está sendo detectado corretamente (verifique logs)

### Performance

Para alta performance:
- Use Redis com configurações otimizadas
- Configure pools de conexão adequados
- Monitore uso de memória do Redis

## Licença

Este projeto é de código aberto e está disponível sob a licença MIT.
