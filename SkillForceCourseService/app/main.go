package main

import (
	"log"
	"net"
	"skillForce/config"
	"skillForce/pkg/logs"

	courseGrpcHandler "skillForce/internal/delivery/grpc/handler"
	coursepb "skillForce/internal/delivery/grpc/proto"
	courseUsecase "skillForce/internal/usecase"

	"skillForce/internal/repository"

	"google.golang.org/grpc"
)

func main() {
	config := config.LoadConfig()

	infrastructure := repository.NewCourseInfrastructure(config)
	defer infrastructure.Close()

	courseUsecase := courseUsecase.NewCourseUsecase(infrastructure)

	lis, err := net.Listen("tcp", ":8082")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(logs.GRPCLoggerInterceptor()),
	)

	courseHandler := courseGrpcHandler.NewCourseHandler(courseUsecase)
	coursepb.RegisterCourseServiceServer(grpcServer, courseHandler)

	log.Println("gRPC Server started on :8082")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
