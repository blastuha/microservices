package main

import (
	"log"

	"github.com/your-org/users-service/internal/database"
	"github.com/your-org/users-service/internal/transport/grpc"
	"github.com/your-org/users-service/internal/user"
)

const dsn = "postgres://postgres:aZAz1998@localhost:5432/users_db?sslmode=disable"

func main() {
	db, err := database.NewDB(dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	userRepo := user.NewUsersRepo(db.Db)
	userService := user.NewUsersService(userRepo)

	// Создаем gRPC сервер
	server := grpc.NewServer(50051)
	server.RegisterServices(userService)

	startErr := server.Start()
	if startErr != nil {
		log.Fatalf("Failed to start gRPC server: %v", err)
	}
}
