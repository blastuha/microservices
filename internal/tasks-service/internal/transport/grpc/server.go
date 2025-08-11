package grpc

import (
	"fmt"
	"net"
	"strconv"

	taskspb "github.com/blastuha/test-service-proto/gen/task"
	"github.com/your-org/tasks-service/internal/tasks"
	"google.golang.org/grpc"
)

type Server struct {
	server *grpc.Server
	Port   int
}

func NewServer(port int) *Server {
	return &Server{
		server: grpc.NewServer(),
		Port:   port,
	}
}

func (s *Server) RegisterServices(svc tasks.TasksService, cl Client) {
	tasksHandler := NewHandler(svc, cl)
	taskspb.RegisterTasksServiceServer(s.server, tasksHandler)
}

func (s *Server) Start() error {
	addr := net.JoinHostPort("", strconv.Itoa(s.Port)) // ":50051"
	ls, err := net.Listen("tcp", addr)

	fmt.Println("gRPC tasks server is running on port", s.Port)

	if err != nil {
		return fmt.Errorf("failed to listen on port %d: %w", s.Port, err)
	}

	if err := s.server.Serve(ls); err != nil {
		return fmt.Errorf("failed to serve gRPC server: %w", err)
	}

	return nil
}

func (s *Server) Stop() {
	s.server.GracefulStop()
	fmt.Println("gRPC server is stopped on port", s.Port)
}
