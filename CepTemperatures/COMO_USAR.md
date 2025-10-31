# 🌐 Como Usar o Serviço no Cloud Run

## 🔗 URL do Serviço

**URL Principal:** `https://cep-temperature-47ocgrvvgq-uc.a.run.app`

## 📡 Endpoints Disponíveis

### 1. Health Check
Verifica se o serviço está funcionando.

**URL:** `https://cep-temperature-47ocgrvvgq-uc.a.run.app/health`

**Exemplo no navegador:**
- Acesse: https://cep-temperature-47ocgrvvgq-uc.a.run.app/health

**Exemplo com curl:**
```bash
curl https://cep-temperature-47ocgrvvgq-uc.a.run.app/health
```

**Resposta:**
```json
{"status":"ok"}
```

---

### 2. Consultar Temperatura por CEP
Retorna a temperatura atual da cidade correspondente ao CEP.

**URL:** `https://cep-temperature-47ocgrvvgq-uc.a.run.app/weather/{cep}`

**Formato:** CEP de 8 dígitos (apenas números, sem hífen)

**Exemplos no navegador:**
- CEP da Avenida Paulista, São Paulo: https://cep-temperature-47ocgrvvgq-uc.a.run.app/weather/01310100
- CEP do Centro de São Paulo: https://cep-temperature-47ocgrvvgq-uc.a.run.app/weather/01001000
- CEP do Rio de Janeiro: https://cep-temperature-47ocgrvvgq-uc.a.run.app/weather/20040020

**Exemplos com curl:**
```bash
# São Paulo - Avenida Paulista
curl https://cep-temperature-47ocgrvvgq-uc.a.run.app/weather/01310100

# São Paulo - Centro
curl https://cep-temperature-47ocgrvvgq-uc.a.run.app/weather/01001000

# Rio de Janeiro
curl https://cep-temperature-47ocgrvvgq-uc.a.run.app/weather/20040020
```

**Resposta de Sucesso (200):**
```json
{
  "temp_C": 16.1,
  "temp_F": 60.98,
  "temp_K": 289.25
}
```

**Campos:**
- `temp_C`: Temperatura em Celsius (°C)
- `temp_F`: Temperatura em Fahrenheit (°F)
- `temp_K`: Temperatura em Kelvin (K)

**Resposta de Erro - CEP Inválido (422):**
```json
{
  "error": "invalid zipcode"
}
```

**Resposta de Erro - CEP Não Encontrado (404):**
```json
{
  "error": "can not find zipcode"
}
```

---

## 🌐 Como Usar no Navegador

1. **Health Check:**
   - Abra seu navegador
   - Acesse: `https://cep-temperature-47ocgrvvgq-uc.a.run.app/health`

2. **Consultar Temperatura:**
   - Abra seu navegador
   - Acesse: `https://cep-temperature-47ocgrvvgq-uc.a.run.app/weather/01310100`
   - Substitua `01310100` pelo CEP desejado (8 dígitos)

---

## 💻 Como Usar no Terminal (curl)

```bash
# Testar health
curl https://cep-temperature-47ocgrvvgq-uc.a.run.app/health

# Consultar temperatura
curl https://cep-temperature-47ocgrvvgq-uc.a.run.app/weather/01310100
```

---

## 🔧 Como Usar com JavaScript (Fetch API)

```javascript
// Health check
fetch('https://cep-temperature-47ocgrvvgq-uc.a.run.app/health')
  .then(response => response.json())
  .then(data => console.log(data));

// Consultar temperatura
const cep = '01310100';
fetch(`https://cep-temperature-47ocgrvvgq-uc.a.run.app/weather/${cep}`)
  .then(response => response.json())
  .then(data => {
    console.log(`Temperatura: ${data.temp_C}°C`);
    console.log(`Temperatura: ${data.temp_F}°F`);
    console.log(`Temperatura: ${data.temp_K}K`);
  });
```

---

## 🐍 Como Usar com Python

```python
import requests

# Health check
response = requests.get('https://cep-temperature-47ocgrvvgq-uc.a.run.app/health')
print(response.json())

# Consultar temperatura
cep = '01310100'
response = requests.get(f'https://cep-temperature-47ocgrvvgq-uc.a.run.app/weather/{cep}')
data = response.json()
print(f"Temperatura: {data['temp_C']}°C")
print(f"Temperatura: {data['temp_F']}°F")
print(f"Temperatura: {data['temp_K']}K")
```

---

## 📊 Monitoramento

### Ver logs do serviço:
```bash
gcloud run services logs read cep-temperature --region us-central1 --project cep-temperature-476519
```

### Ver informações do serviço:
```bash
gcloud run services describe cep-temperature --region us-central1 --project cep-temperature-476519
```

### Ver métricas:
```bash
gcloud run services describe cep-temperature --region us-central1 --project cep-temperature-476519 --format="yaml(status)"
```

---

## 🔄 Atualizar o Serviço

Para fazer um novo deploy após mudanças no código:

```bash
./deploy-local.sh cep-temperature-476519
```

Para atualizar apenas variáveis de ambiente:

```bash
gcloud run services update cep-temperature \
  --region us-central1 \
  --project cep-temperature-476519 \
  --set-env-vars="WEATHER_API_KEY=sua_chave_aqui"
```

---

## ❓ Problemas Comuns

### Erro 404 - CEP não encontrado
- Verifique se o CEP tem exatamente 8 dígitos
- Verifique se o CEP existe
- CEPs válidos: apenas números, sem hífen (ex: 01310100, não 01310-100)

### Serviço não responde
- Verifique os logs: `gcloud run services logs read cep-temperature --region us-central1`
- Verifique se o serviço está rodando: `gcloud run services describe cep-temperature --region us-central1`

### Erro 422 - CEP inválido
- O CEP deve ter exatamente 8 dígitos
- Use apenas números (sem hífen, pontos ou espaços)

