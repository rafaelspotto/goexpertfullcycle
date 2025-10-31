#!/bin/bash

# Script de Deploy usando Cloud Build com código fonte
# Este script faz upload do código e usa Cloud Build para construir e fazer deploy

set -e

PROJECT_ID="starry-journal-459919-u7"
SERVICE_NAME="cep-temperature"
REGION="us-central1"

echo "🚀 Iniciando deploy do código Go para Cloud Run..."

# Configurar projeto
echo "📋 Configurando projeto..."
gcloud config set project $PROJECT_ID

# Criar um arquivo temporário com o código
echo "📦 Preparando código para upload..."
cd /home/rafaelspotto/goexpertfullcycle/CepTemperatures

# Criar um arquivo tar com o código fonte
tar -czf /tmp/cep-temperature-source.tar.gz \
  --exclude='.git' \
  --exclude='node_modules' \
  --exclude='*.log' \
  --exclude='.env' \
  .

echo "📤 Fazendo upload do código e iniciando build..."
gcloud builds submit \
  --config cloudbuild.yaml \
  /tmp/cep-temperature-source.tar.gz

echo "✅ Deploy concluído!"
echo "🔗 URL do serviço: https://cep-temperature-auvoswhfkq-uc.a.run.app"
echo ""
echo "🧪 Teste a API:"
echo "curl https://cep-temperature-auvoswhfkq-uc.a.run.app/health"
echo "curl https://cep-temperature-auvoswhfkq-uc.a.run.app/weather/01310100"

# Limpar arquivo temporário
rm -f /tmp/cep-temperature-source.tar.gz

