package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/sigit14ap/personal-finance/auth-service/internal/handler"
	pb "github.com/sigit14ap/personal-finance/auth-service/internal/proto"
	"github.com/sigit14ap/personal-finance/auth-service/internal/repositories"
	"github.com/sigit14ap/personal-finance/auth-service/internal/service"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	log.Info("[Authentication Service] Start")

	// Load environment from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")

	//Initialize Database connection
	mongoClient, err := mongo.Connect(context.Background(), options.Client().ApplyURI(dbHost))

	if err != nil {
		log.Fatal(err)
	}

	db := mongoClient.Database(dbName)

	// Initialize repositories
	_ = repositories.NewRepositories(db)

	//Initialize gRPC server
	grpcServer := grpc.NewServer()
	pb.RegisterAuthServiceServer(grpcServer, handler.NewAuthHandler(service.NewAuthService(*repositories.NewUserRepository(db))))

	reflection.Register(grpcServer)

	// Start gRPC server
	go func() {
		grpcPort := fmt.Sprintf(":%s", os.Getenv("GRPC_PORT"))
		listener, err := net.Listen("tcp", grpcPort)
		if err != nil {
			log.Fatalf("Failed to listen: %v", err)
		}

		log.Info(fmt.Sprintf("gRPC server is running on %s", grpcPort))
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down the server...")

	// Stop the gRPC server
	grpcServer.GracefulStop()
}
