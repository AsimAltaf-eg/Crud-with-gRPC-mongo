package main

import (
	"context"
	pb "grpc-mongo/proto"
	"log"
	"net"
	"os"
	"os/signal"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"google.golang.org/grpc"
)

var Blogdb *mongo.Collection
var MongoConn *options.ClientOptions
var Db *mongo.Client
var Ctx context.Context

const (
	port = ":8080"
)

type helloServer struct {
	pb.Blog_ServiceServer
}

func main() {

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to Listen at %v", port)
	}

	grpcServer := grpc.NewServer()

	pb.RegisterBlog_ServiceServer(grpcServer, &helloServer{})

	//Connecting to MongoDB

	Ctx = context.Background()
	MongoConn = options.Client().ApplyURI("mongodb://localhost:27017")
	Db, err = mongo.Connect(Ctx, MongoConn)
	if err != nil {
		log.Fatal("Applicaton Failed to Connect to the Database")
	}
	err = Db.Ping(Ctx, readpref.Primary())
	if err != nil {
		log.Fatal("Applicaton Failed to Connect to the Database")
	}

	Blogdb = Db.Database("grpc").Collection("blog")

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()
	log.Printf("Server succesfully started on port %v", port)

	// Right way to stop the server using a SHUTDOWN HOOK
	// Create a channel to receive OS signals
	c := make(chan os.Signal, 1)

	// Relay os.Interrupt to our channel (os.Interrupt = CTRL+C)
	// Ignore other incoming signals
	signal.Notify(c, os.Interrupt)

	// Block main routine until a signal is received
	// As long as user doesn't press CTRL+C a message is not passed and our main routine keeps running
	<-c

	// After receiving CTRL+C Properly stop the server
	log.Printf("Stopping the server...")
	grpcServer.Stop()
	lis.Close()
	log.Printf("Closing MongoDB connection")
	Db.Disconnect(Ctx)
}
