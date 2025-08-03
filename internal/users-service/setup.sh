#!/bin/bash

# Скрипт автоматической настройки Users Service
set -e

echo "🚀 Настройка Users Service..."

# Проверяем наличие Go
if ! command -v go &> /dev/null; then
    echo "❌ Go не установлен. Установите Go 1.24+"
    exit 1
fi

echo "✅ Go найден: $(go version)"

# Устанавливаем oapi-codegen если не установлен
if ! command -v oapi-codegen &> /dev/null; then
    echo "📦 Устанавливаем oapi-codegen..."
    go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
else
    echo "✅ oapi-codegen уже установлен"
fi

# Устанавливаем migrate если не установлен
if ! command -v migrate &> /dev/null; then
    echo "📦 Устанавливаем migrate..."
    go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
else
    echo "✅ migrate уже установлен"
fi

# Генерируем код из OpenAPI
echo "🔧 Генерируем код из OpenAPI спецификации..."
make gen-users

# Проверяем, что все зависимости установлены
echo "🔍 Проверяем зависимости..."
go mod verify

echo "✅ Настройка завершена!"
echo ""
echo "📋 Следующие шаги:"
echo "1. Настройте переменную окружения DB_DSN"
echo "2. Запустите миграции: make migrate"
echo "3. Запустите сервис: make run" 