package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/IAmFutureHokage/HL-BufferService/internal/app/migration"
	"github.com/IAmFutureHokage/HL-BufferService/internal/app/repository"
	"github.com/IAmFutureHokage/HL-BufferService/internal/app/services"
	pb "github.com/IAmFutureHokage/HL-BufferService/internal/proto"
	"github.com/IAmFutureHokage/HL-BufferService/pkg/database"
	"github.com/IAmFutureHokage/HL-BufferService/pkg/kafka"
	"github.com/Shopify/sarama"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

var dbConfig database.Config
var kafkaConfig kafka.KafkaConfig
var kafkaProducer sarama.SyncProducer

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

	kafkaConfig = kafka.KafkaConfig{
		BrokerList: viper.GetStringSlice("kafka.broker_list"),
		Topic:      viper.GetString("kafka.topic"),
	}

	var err error
	kafkaProducer, err = kafka.NewKafkaProducer(kafkaConfig)
	if err != nil {
		log.Fatalf("Error creating Kafka producer: %v", err)
	}
}

func main() {

	fmt.Println("gRPC server running ...")

	port := viper.GetInt("server.port")
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

	hydrologyBufferRepository := repository.NewHydrologyBufferRepository(dbPool)
	hydrologyBufferService := services.NewHydrologyBufferService(hydrologyBufferRepository, kafkaProducer)
	hydrologyBufferService.SetKafkaConfig(kafkaConfig)

	s := grpc.NewServer()
	pb.RegisterHydrologyBufferServiceServer(s, hydrologyBufferService)

	log.Printf("Server listening at %v", lis.Addr())

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve : %v", err)
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	<-shutdown

	fmt.Println("Shutting down server...")
	s.GracefulStop()
	if err := kafkaProducer.Close(); err != nil {
		log.Printf("Error closing Kafka producer: %v", err)
	}
	fmt.Println("Server gracefully stopped")
}
