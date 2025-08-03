package grpc

import (
	"fmt"
	"net"

	userpb "github.com/blastuha/test-service-proto/gen/user"
	"github.com/your-org/users-service/internal/user"
	"google.golang.org/grpc"
)

type Server struct {
	grpcServer *grpc.Server
	port       int
}

func NewServer(port int) *Server {
	return &Server{
		grpcServer: grpc.NewServer(),
		port:       port,
	}
}

func (s *Server) RegisterServices(userService user.UsersService) {
	// Регистрируем gRPC обработчики
	userHandler := NewHandler(userService)
	userpb.RegisterUserServiceServer(s.grpcServer, userHandler)
}

func (s *Server) Start() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	fmt.Printf("gRPC server listening on port %d\n", s.port)
	return s.grpcServer.Serve(lis)
}

func (s *Server) Stop() {
	if s.grpcServer != nil {
		s.grpcServer.GracefulStop()
	}
}
