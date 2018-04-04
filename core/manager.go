package core

import (
	stream "Multy-back-exchange-service/stream/server"
	"sync"
	//"time"
	"log"

	"time"
)

type Manager struct {
	binanceManager *BinanceManager
	hitBtcManager *HitBtcManager
	poloniexManager *PoloniexManager
	bitfinexManager *BitfinexManager
	gdaxManager *GdaxManager
	okexManager *OkexManager
	server *stream.Server
	agregator *Agregator

	waitGroup sync.WaitGroup
}

func NewManager() *Manager {
	var manger = Manager{}
	manger.binanceManager = &BinanceManager{}
	manger.hitBtcManager = &HitBtcManager{}
	manger.poloniexManager = &PoloniexManager{}
	manger.bitfinexManager = &BitfinexManager{}
	manger.gdaxManager = &GdaxManager{}
	manger.okexManager = &OkexManager{}
	manger.server = &stream.Server{}
	manger.agregator = NewAgregator()

	return &manger
}

func (b *Manager) StartListen() {

	b.waitGroup.Add(7)

	go b.binanceManager.StartListen( func(tickerCollection TickerCollection, error error) {
		if error != nil {
			log.Println("error:", error)
		} else {
			//fmt.Println(tickerCollection)
			b.agregator.add(tickerCollection, "Binance")
		}
	} )

	go b.hitBtcManager.StartListen( func(tickerCollection TickerCollection, error error) {
		if error != nil {
			log.Println("error:", error)
		} else {
			//fmt.Println(tickerCollection)
			b.agregator.add(tickerCollection, "HitBtc")
		}
	} )
	//
	go b.poloniexManager.StartListen( func(tickerCollection TickerCollection, error error) {
		if error != nil {
			log.Println("error:", error)
		} else {
			//fmt.Println(tickerCollection)
			b.agregator.add(tickerCollection, "Poloniex")
		}
	} )

	go b.bitfinexManager.StartListen( func(tickerCollection TickerCollection, error error) {
		if error != nil {
			log.Println("error:", error)
		} else {
			//fmt.Println(tickerCollection)
			b.agregator.add(tickerCollection, "Bitfinex")
		}
	} )

	go b.gdaxManager.StartListen( func(tickerCollection TickerCollection, error error) {
		if error != nil {
			log.Println("error:", error)
		} else {
			//fmt.Println(tickerCollection)
			b.agregator.add(tickerCollection, "Gdax")
		}
	} )

	go b.okexManager.StartListen( func(tickerCollection TickerCollection, error error) {
		if error != nil {
			log.Println("error:", error)
		} else {
			//fmt.Println(tickerCollection)
			b.agregator.add(tickerCollection, "Okex")
		}
	} )
	//
	go b.server.StartServer()
	b.server.ServerHandler =  func(allTickers *map[string]stream.StreamTickerCollection) {

		var tickerCollections = b.agregator.getTickers(time.Now().Add(-3 * time.Second))

		var streamTickerCollections = make(map[string]stream.StreamTickerCollection)

		for key, tickerColection := range tickerCollections {
			var streamTickerColection = b.convertToTickerCollection(tickerColection)
			streamTickerCollections[key] = streamTickerColection
		}
		allTickers = &streamTickerCollections
	}

	b.waitGroup.Wait()

}

func (b *Manager) convertToTickerCollection (tickerCollection TickerCollection) stream.StreamTickerCollection {
	var streamTickerCollection = stream.StreamTickerCollection{}
	var streamTickers = []stream.StreamTicker{}

	streamTickerCollection.TimpeStamp = tickerCollection.TimpeStamp
	for _, ticker := range tickerCollection.Tickers {
		var streamTicker = b.convertToStreamTicker(ticker)
		streamTickers = append(streamTickers, streamTicker)
	}
	streamTickerCollection.Tickers = streamTickers

	return streamTickerCollection

}

func (b *Manager) convertToStreamTicker (ticker Ticker) stream.StreamTicker {
	var streamTicker = stream.StreamTicker{}
	streamTicker.Symbol = ticker.Symbol
	streamTicker.Rate = ticker.Rate
	return streamTicker
}