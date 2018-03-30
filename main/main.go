package main

import (
	"time"
	"binanceParser/Api"
	"log"
	stream "binanceParser/Stream/Server"
	"sync"
	"fmt"
)

var binanceManager = Api.BinanceManager{}
var hitBtcManager = Api.HitBtcManager{}
var server = stream.Server{}

var allTickers = make(map[string]Api.TickerCollection)
var waitGroup sync.WaitGroup

func main() {

	waitGroup.Add(3)

	go binanceManager.StartListen( func(tickerCollection Api.TickerCollection, error error) {
		if error != nil {
			log.Println("error:", error)
		} else {
			//fmt.Println(tickerCollection)
			add(tickerCollection, "Binance")
		}
	} )

	go hitBtcManager.StartListen( func(tickerCollection Api.TickerCollection, error error) {
		if error != nil {
			log.Println("error:", error)
		} else {
			//fmt.Println(tickerCollection)
			add(tickerCollection, "HitBtc")
		}
	} )

	go server.StartServer()
	server.ServerHandler =  func(allTickers *map[string]Api.TickerCollection) {
		*allTickers = getTickers(time.Now().Add(-3 * time.Second))
	}

	waitGroup.Wait()

}


func add(tickerCollection Api.TickerCollection, forExchange string) {
	allTickers[forExchange] = tickerCollection
}

func getTickers(startDate time.Time) map[string]Api.TickerCollection {
	var filteredTickers = make(map[string]Api.TickerCollection)
	for exhange, tickerCollection := range allTickers {
		if tickerCollection.TimpeStamp.After(startDate) {
			filteredTickers[exhange] = tickerCollection
		}
	}

	return filteredTickers
}

