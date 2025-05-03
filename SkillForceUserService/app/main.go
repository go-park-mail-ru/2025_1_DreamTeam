package main

import (
	"log"
	"net"
	"skillForce/config"
	"skillForce/pkg/logs"

	userGrpcHandler "skillForce/internal/delivery/grpc/handler"
	userpb "skillForce/internal/delivery/grpc/proto"
	userUsecase "skillForce/internal/usecase"

	"skillForce/internal/repository"

	"google.golang.org/grpc"
)

func main() {
	config := config.LoadConfig()

	infrastructure := repository.NewUserInfrastructure(config)
	defer infrastructure.Close()

	userUsecase := userUsecase.NewUserUsecase(infrastructure)

	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(logs.GRPCLoggerInterceptor()),
	)

	userHandler := userGrpcHandler.NewUserHandler(userUsecase)
	userpb.RegisterUserServiceServer(grpcServer, userHandler)

	log.Println("gRPC Server started on :8081")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
