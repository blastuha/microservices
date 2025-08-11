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
	port   int
}

func NewServer(port int) *Server {
	return &Server{
		server: grpc.NewServer(),
		port:   port,
	}
}

func (s *Server) RegisterServices(svc tasks.TasksService, cl Client) {
	tasksHandler := NewHandler(svc, cl)
	taskspb.RegisterTasksServiceServer(s.server, tasksHandler)
}

func (s *Server) Start() error {
	ls, err := net.Listen("tcp", strconv.Itoa(s.port))

	fmt.Println("gRPC tasks server is running on port", s.port)

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
