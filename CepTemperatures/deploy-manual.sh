#!/bin/bash

# Script de Deploy Manual para Google Cloud Run
# Este script usa o Cloud Run diretamente sem Cloud Build

set -e

PROJECT_ID="starry-journal-459919-u7"
SERVICE_NAME="cep-temperature"
REGION="us-central1"

echo "🚀 Iniciando deploy manual do serviço $SERVICE_NAME no projeto $PROJECT_ID..."

# Configurar projeto
echo "📋 Configurando projeto..."
gcloud config set project $PROJECT_ID

# Verificar se as APIs estão habilitadas
echo "🔍 Verificando APIs necessárias..."
gcloud services enable run.googleapis.com

# Deploy usando uma imagem pré-construída do Go
echo "🏗️ Fazendo deploy usando imagem base do Go..."
gcloud run deploy $SERVICE_NAME \
  --image gcr.io/cloudrun/hello \
  --region $REGION \
  --platform managed \
  --allow-unauthenticated \
  --port 8080 \
  --memory 512Mi \
  --cpu 1 \
  --max-instances 10

echo "✅ Deploy concluído com sucesso!"
echo ""
echo "📝 Próximos passos:"
echo "1. Você precisará fazer upload do seu código Go para o Cloud Run"
echo "2. Configure as variáveis de ambiente necessárias"
echo "3. Teste a aplicação"
echo ""
echo "🔗 Para obter a URL do serviço:"
echo "gcloud run services describe $SERVICE_NAME --region=$REGION --format='value(status.url)'"

