package main

import (
	"fmt"
	"github.com/sunil206b/customer_api/src/customerpb"
	"github.com/sunil206b/customer_api/src/driver"
	"github.com/sunil206b/customer_api/src/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
	"os/signal"
)

func main() {
	// If there is any error we will get the file name and line number
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Starting the Customer GRP Server...")

	port := os.Getenv("PORT")
	if port == "" {
		port = "50051"
	}
	list, err := net.Listen("tcp", "0.0.0.0:"+port)
	if err != nil {
		log.Fatalf("Failed to listen on the port %s: %v\n", port, err)
		return
	}

	client, ctx, err := driver.MongoDB()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	collection := client.Database("customerdb").Collection("customer")

	s := grpc.NewServer()
	server := service.NewServer(collection)
	customerpb.RegisterCustomerServiceServer(s, server)
	reflection.Register(s)

	go func() {
		fmt.Println("Starting Server...")
		if err = s.Serve(list); err != nil {
			log.Fatalf("Failed to start serve: %v\n", err)
		}
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	<-ch
	log.Println("Stopping the server...")
	s.Stop()
	log.Println("Closing the listener...")
	list.Close()
	log.Println("End of Program...")
}
