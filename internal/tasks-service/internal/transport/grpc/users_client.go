package grpc

import (
	"context"
	"fmt"
	"time"

	userspb "github.com/blastuha/test-service-proto/gen/user"
	"github.com/your-org/tasks-service/domain"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

var (
	defaultTimeout  = 2 * time.Second
	ErrUserNotFound = fmt.Errorf("пользователь не найден")
	ErrUnavailable  = fmt.Errorf("сервис пользователей недоступен")
)

// Client определяет интерфейс для клиента сервиса пользователей.
type Client interface {
	GetUser(ctx context.Context, id uint32) (*domain.User, error)
}

type client struct {
	raw userspb.UserServiceClient
}

// New создает новый клиент сервиса пользователей, инкапсулируя логику подключения.
// Возвращает клиент, функцию для закрытия соединения и ошибку.
func New(ctx context.Context, addr string) (Client, func(), error) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	conn, err := grpc.NewClient(addr, opts...)
	if err != nil {
		return nil, nil, fmt.Errorf("не удалось подключиться к сервису пользователей: %w", err)
	}

	cleanup := func() {
		conn.Close()
	}

	c := &client{
		raw: userspb.NewUserServiceClient(conn),
	}

	return c, cleanup, nil
}

// GetUser получает пользователя по его ID.
func (c *client) GetUser(ctx context.Context, id uint32) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	resp, err := c.raw.GetUser(ctx, &userspb.GetUserRequest{Id: id})
	if err != nil {
		st, _ := status.FromError(err)
		switch st.Code() {
		case codes.NotFound:
			return nil, ErrUserNotFound
		case codes.Unavailable, codes.DeadlineExceeded:
			return nil, ErrUnavailable
		default:
			return nil, fmt.Errorf("GetUser: %w", err)
		}
	}
	return &domain.User{ID: resp.GetId(), Email: resp.GetEmail()}, nil
}
