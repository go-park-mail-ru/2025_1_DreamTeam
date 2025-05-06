package main

import (
	"log"
	"net"
	"skillForce/config"
	"skillForce/pkg/logs"
	"skillForce/pkg/metrics"

	userGrpcHandler "skillForce/internal/delivery/grpc/handler"
	userpb "skillForce/internal/delivery/grpc/proto"
	userUsecase "skillForce/internal/usecase"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"

	"skillForce/internal/repository"

	"google.golang.org/grpc"
)

func main() {
	config := config.LoadConfig()

	infrastructure := repository.NewUserInfrastructure(config)
	defer infrastructure.Close()

	userUsecase := userUsecase.NewUserUsecase(infrastructure)

	metrics.Init(":9081")

	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				logs.GRPCLoggerInterceptor(),
				grpc_prometheus.UnaryServerInterceptor,
			),
		),
	)

	userHandler := userGrpcHandler.NewUserHandler(userUsecase)
	userpb.RegisterUserServiceServer(grpcServer, userHandler)
	grpc_prometheus.Register(grpcServer)

	log.Println("gRPC Server started on :8081")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
