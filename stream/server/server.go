package server


import (
	"flag"
	"fmt"
	"log"
	"net"
	"google.golang.org/grpc"
	"time"
)

var (
	port = flag.Int("port", 10000, "The server port")
)

type StreamTickerCollection struct {
	TimpeStamp time.Time
	Tickers []StreamTicker
}

type StreamTicker struct {
	Symbol 	string
	Rate	string
}


type Server struct {
	ServerHandler func(*map[string]StreamTickerCollection)
}

func (s *Server) Tickers(whoAreYouParams *WhoAreYouParams, stream TickerGRPCServer_TickersServer) error {

	for range time.Tick(1 * time.Second) {
		if streemError := stream.Context().Err(); streemError != nil  {
			println("error getting contex from client: ", streemError)
			break
		}

		var allTickers = make(map[string]StreamTickerCollection)

		s.ServerHandler(&allTickers)

		var streamTickers = Tickers{}
		streamTickers.ExchangeTickers = []*ExchangeTickers{}

		for exchange, tickers := range allTickers {
			var exhangeTickers = ExchangeTickers{}
			exhangeTickers.Exchange = exchange
			exhangeTickers.TimpeStamp = tickers.TimpeStamp.Unix()

			var nodeTickers = []*Ticker{}
			for _, ticker := range tickers.Tickers {
				var nodeTicker = Ticker{}
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
	RegisterTickerGRPCServerServer(grpcServer, s)
	grpcServer.Serve(lis)
}