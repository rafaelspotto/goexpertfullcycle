#!/bin/bash

# Script de Deploy para Google Cloud Run
# Uso: ./deploy.sh [PROJECT_ID]

set -e

# Verificar se PROJECT_ID foi fornecido
if [ -z "$1" ]; then
    echo "❌ Erro: PROJECT_ID é obrigatório"
    echo "Uso: ./deploy.sh [PROJECT_ID]"
    echo "Exemplo: ./deploy.sh meu-projeto-gcp"
    exit 1
fi

PROJECT_ID=$1
SERVICE_NAME="cep-temperature"
REGION="us-central1"

echo "🚀 Iniciando deploy do serviço $SERVICE_NAME no projeto $PROJECT_ID..."

# Verificar se gcloud está instalado
if ! command -v gcloud &> /dev/null; then
    echo "❌ Erro: Google Cloud SDK não está instalado"
    echo "Instale em: https://cloud.google.com/sdk/docs/install"
    exit 1
fi

# Configurar projeto
echo "📋 Configurando projeto..."
gcloud config set project $PROJECT_ID

# Verificar se as APIs estão habilitadas
echo "🔍 Verificando APIs necessárias..."
gcloud services enable cloudbuild.googleapis.com
gcloud services enable run.googleapis.com
gcloud services enable containerregistry.googleapis.com

# Executar build e deploy
echo "🏗️ Executando build e deploy..."
gcloud builds submit --config cloudbuild.yaml

# Obter URL do serviço
echo "🔗 Obtendo URL do serviço..."
SERVICE_URL=$(gcloud run services describe $SERVICE_NAME --region=$REGION --format="value(status.url)")

echo "✅ Deploy concluído com sucesso!"
echo "🌐 URL do serviço: $SERVICE_URL"
echo ""
echo "🧪 Teste a API:"
echo "curl $SERVICE_URL/health"
echo "curl $SERVICE_URL/weather/01310100"
echo ""
echo "📊 Monitoramento:"
echo "gcloud run services describe $SERVICE_NAME --region=$REGION"
