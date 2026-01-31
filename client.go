package main

import (
	"fmt"
	"log"
	"time"
	"context"
	"google.golang.org/grpc"
	"github.com/0ppliger/open-asset-gateway/api"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	fmt.Println("Hello world !")
	cnx, err := grpc.NewClient(
		"localhost:9000",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal("Failed to connect to gRPC server")
	}
	defer cnx.Close()

	c := api.NewGreetingClient(cnx)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	out, err := c.SayHello(ctx, &api.SayHelloInput{Value: "Julien"})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%s\n", out.Value)
}
