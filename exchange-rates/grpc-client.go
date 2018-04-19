package exchangeRates

import (
	"flag"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/Appscrunch/Multy-back-exchange-service/stream/server"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type GrpcClient struct {
	serverAddr *string
	client     server.TickerGRPCServerClient
}

func NewGrpcClient() *GrpcClient {
	grpcClient := GrpcClient{}
	grpcClient.serverAddr = flag.String("server_addr", "127.0.0.1:10000", "The server address in the format of host:port")

	log.Println("starting client")
	flag.Parse()
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial(*grpcClient.serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	//defer conn.Close()
	grpcClient.client = server.NewTickerGRPCServerClient(conn)

	return &grpcClient
}

func (b *GrpcClient) connectToServer() (server.TickerGRPCServer_TickersClient, error) {
	log.Printf("Connecting to GRPC server")
	ctx, _ := context.WithCancel(context.Background())
	stream, error := b.client.Tickers(ctx, &server.Empty{})
	return stream, error
	//defer cancel()
}

func (b *GrpcClient) printAllTickers(ch chan *server.Tickers) {

	stream, err := b.connectToServer()
	if err != nil {
		fmt.Println(err)
	}
	for range time.Tick(1 * time.Second) {

		if stream == nil {
			stream, err = b.connectToServer()
			fmt.Println(err)
		} else {
			feature, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				//log.Fatalf("%v.ListFeatures(_) = _, %v", b.client, err)
			}
			//fmt.Println("before sendign")
			ch <- feature
			//for _, exchangeTicker := range feature.ExchangeTickers {
			//	fmt.Println(exchangeTicker.Exchange, exchangeTicker.TimpeStamp, exchangeTicker.Tickers)
			//}
		}
	}
}

func main() {

	//log.Println("starting client")
	//flag.Parse()
	//var opts []grpc.DialOption
	//opts = append(opts, grpc.WithInsecure())
	//conn, err := grpc.Dial(*serverAddr, opts...)
	//if err != nil {
	//	log.Fatalf("fail to dial: %v", err)
	//}
	//defer conn.Close()
	//client := server.NewTickerGRPCServerClient(conn)
	//
	//printAllTickers(client)
}
