package server


import (
	"flag"
	"fmt"
	"log"
	"net"
	"google.golang.org/grpc"
	"time"
	"Multy-back-exchange-service/api"
	"Multy-back-exchange-service/stream/streamDescription"
)

var (
	port = flag.Int("port", 10000, "The server port")
)


type Server struct {
	ServerHandler func(*map[string]api.TickerCollection)
}

func (s *Server) Tickers(whoAreYouParams *streamDescription.WhoAreYouParams, stream streamDescription.TickerGRPCServer_TickersServer) error {

	for range time.Tick(1 * time.Second) {
		if streemError := stream.Context().Err(); streemError != nil  {
			println("error getting contex from client: ", streemError)
			break
		}

		var allTickers = make(map[string]api.TickerCollection)

		s.ServerHandler(&allTickers)

		var streamTickers = streamDescription.Tickers{}
		streamTickers.ExchangeTickers = []*streamDescription.ExchangeTickers{}

		for exchange, tickers := range allTickers {
			var exhangeTickers = streamDescription.ExchangeTickers{}
			exhangeTickers.Exchange = exchange
			exhangeTickers.TimpeStamp = tickers.TimpeStamp.Unix()

			var nodeTickers = []*streamDescription.Ticker{}
			for _, ticker := range tickers.Tickers {
				var nodeTicker = streamDescription.Ticker{}
				nodeTicker.Rate = ticker.Rate
				nodeTicker.Symbol = ticker.Symbol

				nodeTickers = append(nodeTickers, &nodeTicker)
			}

			exhangeTickers.Tickers = nodeTickers
			streamTickers.ExchangeTickers = append(streamTickers.ExchangeTickers, &exhangeTickers)

		}

		func() {
			if error := stream.Send(&streamTickers); error != nil {
				fmt.Println("error sending to stream: ", error)
			}
		}()
	}
	return nil
}


func (s *Server) StartServer() {
	fmt.Println("starting sever")
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)
	streamDescription.RegisterTickerGRPCServerServer(grpcServer, s)
	grpcServer.Serve(lis)
}