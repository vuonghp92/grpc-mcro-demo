package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/joho/godotenv"
	pbProject "github.com/vuonghp92/grpc-mcro-demo/proto/project"
	pbUser "github.com/vuonghp92/grpc-mcro-demo/proto/user"
	"github.com/vuonghp92/grpc-mcro-demo/shared/interceptor"
	"google.golang.org/grpc"
)

const port = ":50061"

func main() {
	godotenv.Load()
	projectConn, err := grpc.Dial(os.Getenv("PROJECT_SERVICE_ADDR"), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to dial project service: %s", err)
	}
	projectClient := pbProject.NewProjectServiceClient(projectConn)
	srv := grpc.NewServer(grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
		interceptor.XTraceID(),
		interceptor.Logging(),
	)))
	pbUser.RegisterUserServiceServer(srv, &UserService{
		store:         NewStoreOnMemory(),
		projectClient: projectClient,
	})
	go func() {
		listener, err := net.Listen("tcp", port)
		if err != nil {
			log.Fatalf("failed to create listener: %s",
				err)
		}
		log.Println("start server on port", port)
		if err := srv.Serve(listener); err != nil {
			log.Println("failed to exit serve: ", err)
		}
	}()
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGTERM)
	<-sigint
	log.Println("received a signal of graceful shutdown")
	stopped := make(chan struct{})
	go func() {
		srv.GracefulStop()
		close(stopped)
	}()
	ctx, cancel := context.WithTimeout(
		context.Background(), 1*time.Minute)
	select {
	case <-ctx.Done():
		srv.Stop()
	case <-stopped:
		cancel()
	}
	log.Println("completed graceful shutdown")
}
