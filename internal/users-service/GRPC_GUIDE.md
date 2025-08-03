# Руководство по gRPC в Users Service

## Архитектура

Ваш микросервис использует **гибридную архитектуру**, что является **правильным подходом** для микросервисов:

### 🚀 HTTP API (Echo) - для внешних клиентов
- **Назначение**: REST API для веб-приложений, мобильных приложений
- **Протокол**: HTTP/JSON
- **Генерация**: OpenAPI + oapi-codegen
- **Порт**: Обычно 8080

### 🔗 gRPC - для межсервисного взаимодействия
- **Назначение**: Внутреннее взаимодействие между микросервисами
- **Протокол**: HTTP/2 + Protocol Buffers
- **Производительность**: Выше чем REST
- **Порт**: Обычно 9090

## Структура проекта

```
internal/users-service/
├── internal/
│   ├── user/           # Бизнес-логика (пакет: user)
│   │   ├── model.go    # Модели данных
│   │   ├── service.go  # Сервисный слой
│   │   ├── repository.go # Слой данных
│   │   └── errors.go   # Ошибки
│   ├── transport/
│   │   └── grpc/       # gRPC транспорт (пакет: grpc)
│   │       ├── handler.go # gRPC обработчики
│   │       └── server.go  # gRPC сервер
│   └── web/
│       └── users/      # HTTP API (Echo) (пакет: users)
│           └── api.gen.go # Сгенерированный код
├── api/
│   └── users-service/
│       └── openapi.yaml # OpenAPI спецификация
└── go.mod
```

### 📦 Правила именования пакетов

В Go **имя пакета должно соответствовать имени директории**:

- `internal/user/` → пакет `user`
- `internal/transport/grpc/` → пакет `grpc`
- `internal/web/users/` → пакет `users`

Это обеспечивает:
- **Читаемость кода**: легко понять, где находится код
- **Правильные импорты**: `import "github.com/your-org/users-service/internal/transport/grpc"`
- **Соблюдение Go conventions**: стандартные практики Go

### 🔢 Унификация типов ID

**Важно**: Все ID пользователей везде используют тип `uint32`:

- **gRPC**: `uint32` ID
- **Сервис**: `uint32` ID
- **Репозиторий**: `uint32` ID
- **База данных**: `SERIAL` (автоматически конвертируется)

**Преимущества:**
- ✅ Нет лишних конвертаций
- ✅ Консистентность типов
- ✅ Лучшая производительность
- ✅ Меньше ошибок

**Пример:**
```go
// gRPC обработчик
func (h *Handler) GetUser(ctx context.Context, req *userpb.GetUserRequest) (*userpb.User, error) {
    // req.Id уже uint32 - никаких конвертаций!
    userObj, err := h.svc.GetUserByID(req.Id)
    // ...
}

// Сервис
func (u *usersService) GetUserByID(id uint32) (*User, error) {
    // id уже uint32 - никаких конвертаций!
    user, err := u.repo.GetUserByID(id)
    // ...
}

// Репозиторий
func (repo *usersRepo) GetUserByID(id uint32) (*User, error) {
    // GORM автоматически конвертирует uint32 в int для базы
    var u User
    if err := repo.db.First(&u, id).Error; err != nil {
        // ...
    }
    return &u, nil
}
```

## Реализованные gRPC методы

### 1. CreateUser
```protobuf
rpc CreateUser(CreateUserRequest) returns (CreateUserResponse)
```

**Запрос:**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Ответ:**
```json
{
  "user": {
    "id": 1,
    "email": "user@example.com",
    "password": "password123"
  }
}
```

### 2. GetUser
```protobuf
rpc GetUser(GetUserRequest) returns (User)
```

**Запрос:**
```json
{
  "id": 1
}
```

**Ответ:**
```json
{
  "id": 1,
  "email": "user@example.com",
  "password": "password123"
}
```

### 3. UpdateUser
```protobuf
rpc UpdateUser(UpdateUserRequest) returns (User)
```

**Запрос:**
```json
{
  "id": 1,
  "email": "newemail@example.com",
  "password": "newpassword123"
}
```

### 4. ListUsers
```protobuf
rpc ListUsers(ListUsersRequest) returns (ListUsersResponse)
```

**Ответ:**
```json
{
  "users": [
    {
      "id": 1,
      "email": "user1@example.com",
      "password": "password123"
    },
    {
      "id": 2,
      "email": "user2@example.com",
      "password": "password456"
    }
  ]
}
```

### 5. DeleteUser
```protobuf
rpc DeleteUser(DeleteUserRequest) returns (DeleteUserResponse)
```

**Запрос:**
```json
{
  "id": 1
}
```

**Ответ:**
```json
{
  "success": true
}
```

## Использование gRPC сервера

### Запуск сервера

```go
package main

import (
    "log"
    grpcTransport "github.com/your-org/users-service/internal/transport/grpc"
    "github.com/your-org/users-service/internal/user"
    "github.com/your-org/users-service/internal/database"
)

func main() {
    // Инициализация базы данных
    db, err := database.NewConnection()
    if err != nil {
        log.Fatal(err)
    }

    // Создание слоев
    userRepo := user.NewUsersRepo(db)
    userService := user.NewUsersService(userRepo)

    // Создание и запуск gRPC сервера
    grpcServer := grpcTransport.NewServer(9090)
    grpcServer.RegisterServices(userService)

    log.Println("Starting gRPC server on port 9090")
    if err := grpcServer.Start(); err != nil {
        log.Fatal(err)
    }
}
```

### Клиентский код (пример)

```go
package main

import (
    "context"
    "log"
    userpb "github.com/blastuha/test-service-proto/gen/user"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
)

func main() {
    // Подключение к gRPC серверу
    conn, err := grpc.Dial("localhost:9090", grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()

    // Создание клиента
    client := userpb.NewUserServiceClient(conn)

    // Создание пользователя
    createResp, err := client.CreateUser(context.Background(), &userpb.CreateUserRequest{
        Email:    "test@example.com",
        Password: "password123",
    })
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("Created user: %+v", createResp.User)

    // Получение списка пользователей
    listResp, err := client.ListUsers(context.Background(), &userpb.ListUsersRequest{})
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("Users: %+v", listResp.Users)
}
```

## Обработка ошибок

### gRPC коды состояния

gRPC использует стандартные коды состояния из пакета `google.golang.org/grpc/codes`:

- **`codes.OK` (0)**: Успешная операция
- **`codes.InvalidArgument` (3)**: Неверные параметры запроса (аналог HTTP 400)
- **`codes.NotFound` (5)**: Ресурс не найден (аналог HTTP 404)
- **`codes.Internal` (13)**: Внутренняя ошибка сервера (аналог HTTP 500)

### Почему нужна валидация?

**Protocol Buffers НЕ валидируют автоматически:**

```protobuf
message CreateUserRequest {
    string email = 1;    // Может быть пустой строкой!
    string password = 2; // Может быть пустой строкой!
}
```

**Protocol Buffers проверяют только:**
- ✅ Типы данных (string, int, bool)
- ✅ Структуру сообщения
- ✅ Кодировку/декодирование

**Protocol Buffers НЕ проверяют:**
- ❌ Пустые строки
- ❌ Формат email
- ❌ Длину пароля
- ❌ Бизнес-логику

**Поэтому нужна дополнительная валидация:**
```go
func validateEmail(email string) error {
    if email == "" {
        return ValidationError{Field: "email", Message: "email is required"}
    }
    
    emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
    if !emailRegex.MatchString(email) {
        return ValidationError{Field: "email", Message: "invalid email format"}
    }
    
    return nil
}
```

### 🎯 Логика валидации

**Принцип: Валидируем один раз, используем везде**

```go
// 1. Валидация в начале метода
func (h *Handler) UpdateUser(ctx context.Context, req *userpb.UpdateUserRequest) (*userpb.User, error) {
    // Валидируем ВСЕ поля сразу
    if err := validateUpdateUserRequest(req.Id, req.Email, req.Password); err != nil {
        return nil, handleValidationError(err)
    }
    
    // 2. После валидации - никаких дополнительных проверок!
    updateReq := users.UpdateUserRequest{}
    
    if req.Email != nil {
        // req.Email уже валидирован - можно безопасно использовать
        email := openapi_types.Email(*req.Email)
        updateReq.Email = &email
    }
    
    // 3. Передаем в сервис
    updatedUser, err := h.svc.UpdateUser(req.Id, updateReq)
    // ...
}
```

**Преимущества:**
- ✅ Нет дублирования валидации
- ✅ Четкое разделение ответственности
- ✅ Лучшая производительность
- ✅ Легче поддерживать код

### ✅ Финальная проверка

**Все исправлено и правильно:**

1. **Унификация типов ID**: Везде `uint32`
2. **Консистентная валидация**: Используем функции из `validation.go`
3. **Нет дублирования**: Валидируем один раз в начале метода
4. **Правильные пакеты**: `package grpc` в `transport/grpc/`
5. **Четкая архитектура**: gRPC → Сервис → Репозиторий → База данных
6. **Чистый код**: Убрали неиспользуемые функции и дублирование

**Структура валидации:**
```go
// validation.go - все функции валидации
validateEmail(email string) error
validatePassword(password string) error  
validateUserID(id uint32) error
validateCreateUserRequest(email, password string) error
validateUpdateUserRequest(id uint32, email, password *string) error
handleValidationError(err error) error

// handler.go - используем валидацию
func (h *Handler) UpdateUser(...) {
    if err := validateUpdateUserRequest(...); err != nil {
        return nil, handleValidationError(err)
    }
    // После валидации - никаких проверок!
}

// service.go - никакой валидации!
func (u *usersService) UpdateUser(...) {
    // Валидация уже выполнена в gRPC handler
    // Просто обновляем данные
}
```

**Принцип: Валидация только на границе (gRPC handler)**

## Преимущества гибридной архитектуры

1. **Гибкость**: HTTP API для внешних клиентов, gRPC для внутренних
2. **Производительность**: gRPC быстрее для межсервисного взаимодействия
3. **Типобезопасность**: Protocol Buffers обеспечивают строгую типизацию
4. **Масштабируемость**: gRPC поддерживает стриминг и бидирекциональную связь
5. **Совместимость**: HTTP API остается доступным для существующих клиентов

## Следующие шаги

1. **Добавить аутентификацию** в gRPC методы
2. **Реализовать middleware** для логирования и метрик
3. **Добавить валидацию** запросов
4. **Настроить TLS** для безопасного соединения
5. **Добавить health checks** для gRPC сервера 