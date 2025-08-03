package grpc

import (
	"context"

	userpb "github.com/blastuha/test-service-proto/gen/user"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/your-org/users-service/internal/user"
	"github.com/your-org/users-service/internal/web/users"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	svc user.UsersService
	userpb.UnimplementedUserServiceServer
}

func NewHandler(svc user.UsersService) *Handler {
	return &Handler{svc: svc}
}

// CreateUser создает нового пользователя
func (h *Handler) CreateUser(ctx context.Context, req *userpb.CreateUserRequest) (*userpb.CreateUserResponse, error) {
	// Валидируем запрос
	if err := validateCreateUserRequest(req.Email, req.Password); err != nil {
		return nil, handleValidationError(err)
	}

	// Конвертируем gRPC запрос в внутренний формат
	createReq := users.CreateUserRequest{
		Email:    openapi_types.Email(req.Email),
		Password: &req.Password,
	}

	// Создаем пользователя через сервис
	createdUser, err := h.svc.CreateUser(createReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
	}

	// Конвертируем результат в gRPC ответ
	response := &userpb.CreateUserResponse{
		User: &userpb.User{
			Id:       uint32(createdUser.ID),
			Email:    createdUser.Email,
			Password: createdUser.Password,
		},
	}

	return response, nil
}

// GetUser получает пользователя по ID
func (h *Handler) GetUser(ctx context.Context, req *userpb.GetUserRequest) (*userpb.User, error) {
	if err := validateUserID(req.Id); err != nil {
		return nil, handleValidationError(err)
	}

	// Получаем пользователя через сервис
	userObj, err := h.svc.GetUserByID(req.Id)
	if err != nil {
		if err == user.ErrUserNoFound {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get user: %v", err)
	}

	// Конвертируем результат в gRPC ответ
	response := &userpb.User{
		Id:       uint32(userObj.ID),
		Email:    userObj.Email,
		Password: userObj.Password,
	}

	return response, nil
}

// UpdateUser обновляет пользователя
func (h *Handler) UpdateUser(ctx context.Context, req *userpb.UpdateUserRequest) (*userpb.User, error) {
	if err := validateUpdateUserRequest(req.Id, req.Email, req.Password); err != nil {
		return nil, handleValidationError(err)
	}

	// Конвертируем gRPC запрос в внутренний формат
	// Валидация уже выполнена в validateUpdateUserRequest
	updateReq := users.UpdateUserRequest{}

	if req.Email != nil {
		email := openapi_types.Email(*req.Email)
		updateReq.Email = &email
	}
	if req.Password != nil {
		updateReq.Password = req.Password
	}

	// Обновляем пользователя через сервис
	updatedUser, err := h.svc.UpdateUser(req.Id, updateReq)
	if err != nil {
		if err == user.ErrUserNoFound {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to update user: %v", err)
	}

	// Конвертируем результат в gRPC ответ
	response := &userpb.User{
		Id:       uint32(updatedUser.ID),
		Email:    updatedUser.Email,
		Password: updatedUser.Password,
	}

	return response, nil
}

// ListUsers получает список всех пользователей
func (h *Handler) ListUsers(ctx context.Context, req *userpb.ListUsersRequest) (*userpb.ListUsersResponse, error) {
	// Получаем всех пользователей через сервис
	users, err := h.svc.GetAllUsers()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get users: %v", err)
	}

	// Конвертируем результат в gRPC ответ
	response := &userpb.ListUsersResponse{
		Users: make([]*userpb.User, len(users)),
	}

	for i, u := range users {
		response.Users[i] = &userpb.User{
			Id:       uint32(u.ID),
			Email:    u.Email,
			Password: u.Password,
		}
	}

	return response, nil
}

// DeleteUser удаляет пользователя
func (h *Handler) DeleteUser(ctx context.Context, req *userpb.DeleteUserRequest) (*userpb.DeleteUserResponse, error) {
	if err := validateUserID(req.Id); err != nil {
		return nil, handleValidationError(err)
	}

	// Удаляем пользователя через сервис
	err := h.svc.DeleteUser(req.Id)
	if err != nil {
		if err == user.ErrUserNoFound {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to delete user: %v", err)
	}

	// Возвращаем успешный ответ
	response := &userpb.DeleteUserResponse{
		Success: true,
	}

	return response, nil
}
