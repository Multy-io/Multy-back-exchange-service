package main

import (
	"flag"
	"io"
	"log"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"Multy-back-exchange-service/stream/server"
	"fmt"
)

var (
	serverAddr         = flag.String("server_addr", "127.0.0.1:10000", "The server address in the format of host:port")
)


func printAllTickers(client server.TickerGRPCServerClient) {
	//log.Printf("Looking for features within %v")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	stream, err := client.Tickers(ctx, &server.WhoAreYouParams{})
	if err != nil {
		log.Fatalf("%v.ListFeatures(_) = _, %v", client, err)
	}
	for {
		feature, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("%v.ListFeatures(_) = _, %v", client, err)
		}
		fmt.Println(feature)
	}
}


func main() {
	log.Println("starting client")
	flag.Parse()
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := server.NewTickerGRPCServerClient(conn)

	printAllTickers(client)
}
