package main

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
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

	//Start Echo server
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(200, "Hello, gRPC and Echo!")
	})

	e.Logger.Fatal(e.Start(":8080"))

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
		listener, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatalf("Failed to listen: %v", err)
		}
		log.Println("gRPC server is running on :50051")
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down the server...")

	// Stop the gRPC server
	grpcServer.GracefulStop()
}
