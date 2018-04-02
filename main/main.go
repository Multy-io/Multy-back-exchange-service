package main

import (
	"time"
	"log"
	"sync"
	//as
	stream "Multy-back-exchange-service/stream/server"
	//"fmt"
	//"google.golang.org/genproto/protobuf/api"
	 exchangeApi "Multy-back-exchange-service/api"

)

var binanceManager = exchangeApi.BinanceManager{}
var hitBtcManager = exchangeApi.HitBtcManager{}
var poloniexManager = exchangeApi.PoloniexManager{}
var server = stream.Server{}

var allTickers = make(map[string]exchangeApi.TickerCollection)
var waitGroup sync.WaitGroup

func main() {

	waitGroup.Add(3)

	go binanceManager.StartListen( func(tickerCollection exchangeApi.TickerCollection, error error) {
		if error != nil {
			log.Println("error:", error)
		} else {
			//fmt.Println(tickerCollection)
			add(tickerCollection, "Binance")
		}
	} )

	go hitBtcManager.StartListen( func(tickerCollection exchangeApi.TickerCollection, error error) {
		if error != nil {
			log.Println("error:", error)
		} else {
			//fmt.Println(tickerCollection)
			add(tickerCollection, "HitBtc")
		}
	} )

	go poloniexManager.StartListen( func(tickerCollection exchangeApi.TickerCollection, error error) {
		if error != nil {
			log.Println("error:", error)
		} else {
			//fmt.Println(tickerCollection)
			add(tickerCollection, "Poloniex")
		}
	} )

	go server.StartServer()
	server.ServerHandler =  func(allTickers *map[string]exchangeApi.TickerCollection) {
		*allTickers = getTickers(time.Now().Add(-3 * time.Second))
	}

	waitGroup.Wait()

}


func add(tickerCollection exchangeApi.TickerCollection, forExchange string) {
	allTickers[forExchange] = tickerCollection
}

func getTickers(startDate time.Time) map[string]exchangeApi.TickerCollection {
	var filteredTickers = make(map[string]exchangeApi.TickerCollection)
	for exhange, tickerCollection := range allTickers {
		if tickerCollection.TimpeStamp.After(startDate) {
			filteredTickers[exhange] = tickerCollection
		}
	}

	return filteredTickers
}

