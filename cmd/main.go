package main

import (
	"fmt"
	"log"
	"net"

	"github.com/IAmFutureHokage/HL-BufferService/internal/app/services"
	pb "github.com/IAmFutureHokage/HL-BufferService/internal/proto"
	"google.golang.org/grpc"
)

func main() {

	fmt.Println("gRPC server running ...")

	port := 50052
	if port == 0 {
		log.Fatal("Server port is not set in the config file")
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	hydrologyBufferService := services.NewHydrologyBufferService()

	s := grpc.NewServer()
	pb.RegisterHydrologyBufferServiceServer(s, hydrologyBufferService)

	log.Printf("Server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve : %v", err)
	}
}
