package server


//var configString = `{
//		"targetCurrencies" : ["BTC", "ETH", "GOLOS", "BTS", "STEEM", "WAVES", "LTC", "BCH", "ETC", "DASH", "EOS"],
//		"referenceCurrencies" : ["USD", "BTC"],
//		"exchanges": ["Binance","Bitfinex","Gdax","HitBtc","Okex","Poloniex"],
//		"refreshInterval" : "3"
//		}`

import (
	"flag"
	"fmt"
	"log"
	"net"
	"google.golang.org/grpc"
	"time"
	"Multy-back-exchange-service/currencies"
)

var (
	port = flag.Int("port", 10000, "The server port")
)

type StreamTickerCollection struct {
	TimpeStamp time.Time
	Tickers []StreamTicker
}

type StreamTicker struct {
	TargetCurrency currencies.Currency
	ReferenceCurrency currencies.Currency
	Symbol 	string
	Rate	string
}


type Server struct {
	ServerHandler func(*map[string]StreamTickerCollection)
	RefreshInterval time.Duration
}

func (s *Server) Tickers(whoAreYouParams *WhoAreYouParams, stream TickerGRPCServer_TickersServer) error {

	for range time.Tick(s.RefreshInterval * time.Second) {
		if streemError := stream.Context().Err(); streemError != nil  {
			println("error getting contex from client: ", streemError)
			break
		}

		var allTickers = make(map[string]StreamTickerCollection)
		s.ServerHandler(&allTickers)
		//fmt.Println(allTickers)



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
				//nodeTicker.Symbol = ticker.Symbol

				nodeTicker.Referrence = ticker.ReferenceCurrency.CurrencyCode()
				nodeTicker.Target = ticker.TargetCurrency.CurrencyCode()


				nodeTickers = append(nodeTickers, &nodeTicker)
			}

			exhangeTickers.Tickers = nodeTickers
			streamTickers.ExchangeTickers = append(streamTickers.ExchangeTickers, &exhangeTickers)

		}

		func() {
			//fmt.Println(streamTickers)
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