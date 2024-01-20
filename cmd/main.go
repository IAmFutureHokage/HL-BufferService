package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/IAmFutureHokage/HL-BufferService/internal/app/migration"
	"github.com/IAmFutureHokage/HL-BufferService/internal/app/services"
	pb "github.com/IAmFutureHokage/HL-BufferService/internal/proto"
	"github.com/IAmFutureHokage/HL-BufferService/pkg/database"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

var dbConfig database.Config

func init() {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	viper.SetConfigName(env)
	viper.AddConfigPath("./config")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	dbConfig = database.Config{
		Host:     viper.GetString("database.host"),
		Port:     viper.GetInt("database.port"),
		User:     viper.GetString("database.user"),
		Password: viper.GetString("database.password"),
		DBName:   viper.GetString("database.dbname"),
		PoolSize: viper.GetInt("database.poolsize"),
	}
}

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

	dbPool, err := database.ConnectDB(dbConfig)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}

	defer database.CloseDB(dbPool)

	if _, err := dbPool.Exec(context.Background(), migration.CreateTablesTelegramAndPhenomenia); err != nil {
		log.Fatalf("Failed to execute migration: %v", err)
	}

	hydrologyBufferService := services.NewHydrologyBufferService()

	s := grpc.NewServer()
	pb.RegisterHydrologyBufferServiceServer(s, hydrologyBufferService)

	log.Printf("Server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve : %v", err)
	}
}
