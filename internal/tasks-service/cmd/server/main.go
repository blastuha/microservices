package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/your-org/tasks-service/internal/database"
	"github.com/your-org/tasks-service/internal/tasks"
	"github.com/your-org/tasks-service/internal/transport/grpc"
)

const (
	dsn              = "postgres://postgres:aZAz1998@localhost:5432/postgres?sslmode=disable"
	tasksServicePort = 50052             // на каком порту слушает tasks сервис
	userServiceAddr  = "localhost:50051" // адрес gRPC user-service
)

func main() {
	// Контекст с отменой по Ctrl+C/SIGTERM
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// БД
	db, err := database.NewDB(dsn)
	if err != nil {
		log.Fatalf("db connect failed: %v", err)
	}

	// Репо/сервис
	repo := tasks.NewTaskRepo(db.Db)
	svc := tasks.NewTasksService(repo)

	// gRPC-клиент к user-service
	userClient, cleanup, err := grpc.NewClient(ctx, userServiceAddr)
	if err != nil {
		log.Fatalf("user client dial failed: %v", err)
	}
	defer cleanup()

	// gRPC-сервер задач
	server := grpc.NewServer(tasksServicePort)
	server.RegisterServices(svc, userClient)

	// Стартуем сервер в отдельной горутине
	errCh := make(chan error, 1)
	go func() {
		if err := server.Start(); err != nil {
			errCh <- err
		}
	}()

	// Ожидаем либо ошибку, либо сигнал завершения
	select {
	case <-ctx.Done():
		// graceful stop
		server.Stop()
		// даём немного времени освободить ресурсы
		time.Sleep(200 * time.Millisecond)
	case err := <-errCh:
		log.Fatalf("server error: %v", err)
	}

}
