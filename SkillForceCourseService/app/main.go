package main

import (
	"log"
	"net"

	"skillForce/config"
	courseGrpcHandler "skillForce/internal/delivery/grpc/handler"
	coursepb "skillForce/internal/delivery/grpc/proto"
	"skillForce/internal/repository"
	"skillForce/internal/usecase"
	"skillForce/pkg/logs"
	"skillForce/pkg/metrics"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"

	"google.golang.org/grpc"
)

func main() {
	cfg := config.LoadConfig()

	infrastructure := repository.NewCourseInfrastructure(cfg)
	defer infrastructure.Close()

	courseUsecase := usecase.NewCourseUsecase(infrastructure)

	metrics.Init(":9082")

	lis, err := net.Listen("tcp", ":8082")
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

	courseHandler := courseGrpcHandler.NewCourseHandler(courseUsecase)
	coursepb.RegisterCourseServiceServer(grpcServer, courseHandler)
	grpc_prometheus.Register(grpcServer)

	log.Println("gRPC Server started on :8082")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
