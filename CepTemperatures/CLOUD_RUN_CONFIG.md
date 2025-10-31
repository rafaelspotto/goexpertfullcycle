# Configuração do Google Cloud Run

## Pré-requisitos
- Conta no Google Cloud Platform
- Projeto criado no GCP
- Google Cloud SDK instalado (`gcloud`)
- API do Cloud Run habilitada
- API do Container Registry habilitada

## Configuração das Variáveis de Ambiente

### No Cloud Run Console:
1. Acesse o Google Cloud Console
2. Navegue para Cloud Run
3. Selecione o serviço `cep-temperature`
4. Clique em "Edit & Deploy New Revision"
5. Na aba "Variables & Secrets", adicione:
   - **WEATHER_API_KEY**: Sua chave da WeatherAPI

### Via gcloud CLI:
```bash
gcloud run services update cep-temperature \
  --region=us-central1 \
  --set-env-vars="WEATHER_API_KEY=sua_chave_da_weatherapi"
```

## Deploy

### Deploy automático via Cloud Build:
```bash
gcloud builds submit --config cloudbuild.yaml
```

### Deploy manual:
```bash
# Build da imagem
docker build -t gcr.io/$PROJECT_ID/cep-temperature .

# Push para Container Registry
docker push gcr.io/$PROJECT_ID/cep-temperature

# Deploy no Cloud Run
gcloud run deploy cep-temperature \
  --image gcr.io/$PROJECT_ID/cep-temperature \
  --region us-central1 \
  --platform managed \
  --allow-unauthenticated
```

## Configurações Recomendadas

### Recursos:
- **CPU**: 1 vCPU
- **Memória**: 512 MiB
- **Concurrency**: 1000 requests por instância
- **Timeout**: 300 segundos

### Rede:
- **Porta**: 8080 (configurada via variável PORT)
- **Acesso**: Permitir tráfego não autenticado

### Monitoramento:
- Logs automáticos habilitados
- Métricas de performance habilitadas
- Alertas configurados para erros 5xx
