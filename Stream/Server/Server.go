package Server


import (
	"flag"
	"fmt"
	"log"
	"net"
	"google.golang.org/grpc"
	"time"
	"binanceParser/Api"
	"binanceParser/Stream/StreamDescription"
)

var (
	port = flag.Int("port", 10000, "The Server port")
)


type Server struct {
	ServerHandler func(*map[string]Api.TickerCollection)
}

func (s *Server) Tickers(whoAreYouParams *StreamDescription.WhoAreYouParams, stream StreamDescription.TickerGRPCServer_TickersServer) error {

	for range time.Tick(1 * time.Second) {
		if streemError := stream.Context().Err(); streemError != nil  {
			println("error getting contex from client: ", streemError)
			break
		}

		var allTickers = make(map[string]Api.TickerCollection)

		s.ServerHandler(&allTickers)

		var streamTickers = StreamDescription.Tickers{}
		streamTickers.ExchangeTickers = []*StreamDescription.ExchangeTickers{}

		for exchange, tickers := range allTickers {
			var exhangeTickers = StreamDescription.ExchangeTickers{}
			exhangeTickers.Exchange = exchange
			exhangeTickers.TimpeStamp = tickers.TimpeStamp.Unix()

			var nodeTickers = []*StreamDescription.Ticker{}
			for _, ticker := range tickers.Tickers {
				var nodeTicker = StreamDescription.Ticker{}
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
	StreamDescription.RegisterTickerGRPCServerServer(grpcServer, s)
	grpcServer.Serve(lis)
}