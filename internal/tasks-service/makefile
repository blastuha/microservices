# переменные
DB_DSN := "postgres://postgres:aZAz1998@localhost:5432/postgres?sslmode=disable"
MIGRATE := migrate -path ./migrations -database $(DB_DSN)

.PHONY: run migrate migrate-down migrate-new

# Таргет для создания новой миграции
migrate-new:
	migrate create -ext sql -dir ./migrations ${NAME}

# Применение миграций
migrate:
	$(MIGRATE) up

# Откат миграций
migrate-down:
	$(MIGRATE) down

# Откат миграций с указанием количества шагов
migrate-down-steps:
	$(MIGRATE) down $(n)

# Принудительное выставление версии
migrate-force:
	$(MIGRATE) force $(v)

# Запуск приложения
run:
	go run cmd/server/main.go