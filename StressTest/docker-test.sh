#!/bin/bash

echo "=== Testes Docker do StressTest ==="
echo ""

echo "1. Teste básico (10 requests, 3 concorrência):"
docker run stresstest --url=http://google.com --requests=10 --concurrency=3
echo ""

echo "2. Teste de média carga (50 requests, 10 concorrência):"
docker run stresstest --url=http://google.com --requests=50 --concurrency=10
echo ""

echo "3. Teste de alta carga (100 requests, 20 concorrência):"
docker run stresstest --url=http://google.com --requests=100 --concurrency=20
echo ""

echo "4. Teste com serviço de teste (pode retornar 503):"
docker run stresstest --url=https://httpbin.org/get --requests=30 --concurrency=5
echo ""

echo "=== Informações da imagem Docker ==="
docker images stresstest
