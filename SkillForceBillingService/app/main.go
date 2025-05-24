package main

import (
	"log"
	"net"

	"skillForce/config"
	billingGrpcHandler "skillForce/internal/delivery/grpc/handler"
	billingpb "skillForce/internal/delivery/grpc/proto"
	"skillForce/internal/repository"
	"skillForce/internal/usecase"
	"skillForce/pkg/logs"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"

	"google.golang.org/grpc"
)

func main() {
	cfg := config.LoadConfig()

	infrastructure := repository.NewBillingInfrastructure(cfg)
	defer infrastructure.Close()

	billingUsecase := usecase.NewBillingUsecase(infrastructure)

	// metrics.Init(":9084")

	lis, err := net.Listen("tcp", ":8084")
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

	billingHandler := billingGrpcHandler.NewBillingHandler(billingUsecase)
	billingpb.RegisterBillingServiceServer(grpcServer, billingHandler)
	grpc_prometheus.Register(grpcServer)

	log.Println("gRPC Server started on :8084")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
