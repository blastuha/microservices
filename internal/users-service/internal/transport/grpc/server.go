package grpc

import (
	"fmt"
	"net"

	userpb "github.com/blastuha/test-service-proto/gen/user"
	"github.com/your-org/users-service/internal/user"
	"google.golang.org/grpc"
)

type Server struct {
	server *grpc.Server
	port   int
}

func NewServer(port int) *Server {
	return &Server{
		server: grpc.NewServer(),
		port:   port,
	}
}

func (s *Server) RegisterServices(userService user.UsersService) {
	// Регистрируем gRPC обработчики
	userHandler := NewHandler(userService)
	userpb.RegisterUserServiceServer(s.server, userHandler)
}

func (s *Server) Start() error {
	ls, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	defer fmt.Println("gRPC server is running on port", s.port)

	if err != nil {
		return fmt.Errorf("failed to listen on port %d: %w", s.port, err)
	}

	if err := s.server.Serve(ls); err != nil {
		return fmt.Errorf("failed to serve gRPC server: %w", err)
	}

	return nil
}

func (s *Server) Stop() {
	s.server.GracefulStop()
	fmt.Println("gRPC server is stopped on port", s.port)
}
