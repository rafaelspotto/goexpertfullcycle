# Sistema de Consulta de Temperatura por CEP

Sistema em Go que recebe um CEP, identifica a cidade e retorna o clima atual em graus Celsius, Fahrenheit e Kelvin.

## 📋 Requisitos

- Docker e Docker Compose instalados
- Conta na [WeatherAPI](https://www.weatherapi.com/) para obter a chave da API
- CEP válido de 8 dígitos (apenas números)

## 🚀 Instalação e Configuração

### 1. Clonar o repositório (se aplicável)

```bash
cd CepTemperatures
```

### 2. Criar o arquivo `.env`

Crie um arquivo `.env` na raiz do projeto com o seguinte conteúdo:

```env
PORT=8080
WEATHER_API_KEY=sua_chave_da_weatherapi
```

**Importante:** 
- Substitua `sua_chave_da_weatherapi` pela sua chave real da WeatherAPI
- Não deixe linhas duplicadas com `WEATHER_API_KEY` no arquivo
- Não adicione espaços ou aspas ao redor da chave

### 3. Obter a chave da WeatherAPI

1. Acesse [https://www.weatherapi.com/](https://www.weatherapi.com/)
2. Crie uma conta gratuita
3. Copie sua chave da API
4. Cole no arquivo `.env`

## 🏃 Como Executar

### Usando Docker Compose (Recomendado)

```bash
# Subir o container
docker compose up -d

# Ver os logs
docker compose logs -f cep-temperature

# Parar o container
docker compose down
```

### Build e Execução Manual

```bash
# Build da imagem
docker compose build

# Executar
docker compose up -d

# Reiniciar após mudanças no .env
docker compose restart
```

## 📡 Como Usar a API

### Endpoint Principal

**GET** `/weather/:cep`

Consulta a temperatura atual da cidade correspondente ao CEP informado.

#### Parâmetros

- `cep` (string): CEP de 8 dígitos (apenas números)

#### Exemplos de Requisição

```bash
# CEP de São Paulo
curl http://localhost:8080/weather/01001000

# CEP da Avenida Paulista
curl http://localhost:8080/weather/01310100

# CEP do Rio de Janeiro
curl http://localhost:8080/weather/20040020
```

#### Resposta de Sucesso (200)

```json
{
  "temp_C": 22.4,
  "temp_F": 72.32,
  "temp_K": 295.55
}
```

**Campos:**
- `temp_C`: Temperatura em Celsius (°C)
- `temp_F`: Temperatura em Fahrenheit (°F)
- `temp_K`: Temperatura em Kelvin (K)

#### Resposta de Erro - CEP Inválido (422)

Quando o CEP não tem 8 dígitos:

```json
{
  "error": "invalid zipcode"
}
```

#### Resposta de Erro - CEP Não Encontrado (404)

Quando o CEP não existe ou não foi encontrado:

```json
{
  "error": "can not find zipcode"
}
```

### Health Check

**GET** `/health`

Verifica se o servidor está funcionando.

```bash
curl http://localhost:8080/health
```

**Resposta:**
```json
{
  "status": "ok"
}
```

## 🧪 Executar Testes

```bash
# Executar todos os testes
go test ./...

# Executar testes com detalhes
go test ./... -v

# Executar testes de um pacote específico
go test ./internal/services
go test ./internal/handlers
```

## 📦 Estrutura do Projeto

```
CepTemperatures/
├── cmd/
│   └── server/
│       └── main.go          # Ponto de entrada da aplicação
├── internal/
│   ├── handlers/
│   │   └── weather_handler.go    # Handlers HTTP
│   ├── models/
│   │   ├── weather.go            # Modelos de dados
│   │   └── response.go           # Modelos de resposta
│   └── services/
│       ├── cep_service.go         # Serviço de consulta CEP (ViaCEP)
│       └── weather_service.go     # Serviço de consulta clima (WeatherAPI)
├── docker-compose.yml
├── Dockerfile
├── go.mod
├── go.sum
└── .env                        # Variáveis de ambiente (criar manualmente)
```

## 🔧 Tecnologias Utilizadas

- **Go 1.25+**: Linguagem de programação
- **Gin**: Framework web HTTP
- **ViaCEP API**: Consulta de CEPs brasileiros
- **WeatherAPI**: Consulta de dados climáticos
- **Docker**: Containerização

## 📝 Fórmulas de Conversão

- **Celsius para Fahrenheit**: `F = C * 1.8 + 32`
- **Celsius para Kelvin**: `K = C + 273.15`

## 🐛 Solução de Problemas

### Erro 401 Unauthorized

- Verifique se a chave da WeatherAPI está correta no arquivo `.env`
- Certifique-se de que não há linhas duplicadas no `.env`
- Reinicie o container após atualizar o `.env`: `docker compose restart`

### Erro "can not find zipcode"

- Verifique se o CEP possui 8 dígitos
- Confirme que o CEP existe no Brasil
- O CEP pode não estar disponível na base da ViaCEP

### Container não inicia

- Verifique se a porta 8080 está disponível
- Execute `docker compose logs` para ver os erros
- Confirme que o Docker está rodando

## 📄 Licença

Este projeto foi desenvolvido para fins educacionais.

## 👤 Autor

Desenvolvido como parte do curso Go Expert Full Cycle.
