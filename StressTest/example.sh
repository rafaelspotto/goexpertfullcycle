#!/bin/bash

echo "=== Exemplos de uso do StressTest ==="
echo ""

echo "1. Teste básico com Google:"
echo "./stresstest --url=http://google.com --requests=10 --concurrency=3"
echo ""

echo "2. Teste de alta carga:"
echo "./stresstest --url=http://google.com --requests=100 --concurrency=10"
echo ""

echo "3. Teste com serviço de teste:"
echo "./stresstest --url=https://httpbin.org/get --requests=50 --concurrency=5"
echo ""

echo "4. Teste com Docker (quando disponível):"
echo "docker run stresstest --url=http://google.com --requests=1000 --concurrency=20"
echo ""

echo "Executando exemplo 1..."
./stresstest --url=http://google.com --requests=10 --concurrency=3
