package core

import (
	stream "Multy-back-exchange-service/stream/server"
	"sync"
	//"time"
	"log"

	"time"
	"strings"
	"fmt"
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
	manger.binanceManager = NewBinanceManager()
	manger.hitBtcManager = &HitBtcManager{}
	manger.poloniexManager = &PoloniexManager{}
	manger.bitfinexManager = &BitfinexManager{}
	manger.gdaxManager = &GdaxManager{}
	manger.okexManager = &OkexManager{}
	manger.server = &stream.Server{}
	manger.agregator = NewAgregator()

	return &manger
}

type ManagerConfiguration struct {
	TargetCurrencies    []string `json:"targetCurrencies"`
	ReferenceCurrencies []string `json:"referenceCurrencies"`
	Exchanges           []string `json:"exchanges"`
	RefreshInterval     string   `json:"refreshInterval"`
}


type Exchange int

func NewExchange(exchangeString string ) Exchange {
	exchanges := map[string]Exchange{"BINANCE":Binance, "BITFINEX":Bitfinex, "GDAX":Gdax, "HITBTC":HitBtc, "OKEX":Okex, "POLONIEX":Poloniex}
	exchange := exchanges[strings.ToUpper(exchangeString)]
	return exchange
}

func (exchange Exchange) String() string {
	exchanges := [...]string {
		"BINANCE",
		"BITFINEX",
		"GDAX",
		"HITBTC",
		"OKEX",
		"POLONIEX"}
	return exchanges[exchange]
}
const (
	Binance 	Exchange = 0
	Bitfinex 	Exchange = 1
	Gdax 		Exchange = 2
	HitBtc 		Exchange = 3
	Okex 		Exchange = 4
	Poloniex 	Exchange = 5
)

type ExchangeConfiguration struct {
	Exchange            Exchange
	TargetCurrencies    []string
	ReferenceCurrencies []string
}



func (b *Manager)lunchExchange(exchangeConfiguration ExchangeConfiguration) {

	switch exchangeConfiguration.Exchange {
	case Binance:
		go b.binanceManager.StartListen(exchangeConfiguration, func(tickerCollection TickerCollection, error error) {
			if error != nil {
				log.Println("error:", error)
			} else {
				//fmt.Println(tickerCollection)
				b.agregator.add(tickerCollection, exchangeConfiguration.Exchange.String())
			}
		} )
	case Bitfinex:
		go b.bitfinexManager.StartListen(exchangeConfiguration, func(tickerCollection TickerCollection, error error) {
			if error != nil {
				log.Println("error:", error)
			} else {
				//fmt.Println(tickerCollection)
				b.agregator.add(tickerCollection, exchangeConfiguration.Exchange.String())
			}
		} )
	case Gdax:
		go b.gdaxManager.StartListen(exchangeConfiguration, func(tickerCollection TickerCollection, error error) {
		if error != nil {
			log.Println("error:", error)
		} else {
			//fmt.Println(tickerCollection)
			b.agregator.add(tickerCollection, exchangeConfiguration.Exchange.String())
		}
	} )
	case HitBtc:
		go b.hitBtcManager.StartListen(exchangeConfiguration, func(tickerCollection TickerCollection, error error) {
		if error != nil {
			log.Println("error:", error)
		} else {
			//fmt.Println(tickerCollection)
			b.agregator.add(tickerCollection, exchangeConfiguration.Exchange.String())
		}
	} )
	case Okex:
		go b.okexManager.StartListen(exchangeConfiguration, func(tickerCollection TickerCollection, error error) {
			if error != nil {
				log.Println("error:", error)
			} else {
				//fmt.Println(tickerCollection)
				b.agregator.add(tickerCollection, exchangeConfiguration.Exchange.String())
			}
		} )
	case Poloniex:
		go b.poloniexManager.StartListen(exchangeConfiguration, func(tickerCollection TickerCollection, error error) {
		if error != nil {
			log.Println("error:", error)
		} else {
			//fmt.Println(tickerCollection)
			b.agregator.add(tickerCollection, exchangeConfiguration.Exchange.String())
		}
	} )
	default:
		return

	}

}

func (b *Manager) StartListen(configuration ManagerConfiguration) {

	exchangesAmount := len(configuration.Exchanges)
	b.waitGroup.Add(exchangesAmount + 1)

	for _, exchangeString := range configuration.Exchanges {
		exchangeConfiguration := ExchangeConfiguration{}
		exchangeConfiguration.Exchange = NewExchange(exchangeString)
		exchangeConfiguration.TargetCurrencies = configuration.TargetCurrencies
		exchangeConfiguration.ReferenceCurrencies = configuration.ReferenceCurrencies
		b.lunchExchange(exchangeConfiguration)
	}


	go b.server.StartServer()
	b.server.ServerHandler =  func(allTickers *map[string]stream.StreamTickerCollection) {

		var tickerCollections = b.agregator.getTickers(time.Now().Add(-3 * time.Second))
		//fmt.Println(tickerCollections)
		var streamTickerCollections = make(map[string]stream.StreamTickerCollection)

		for key, tickerColection := range tickerCollections {
			var streamTickerColection = b.convertToTickerCollection(tickerColection)
			streamTickerCollections[key] = streamTickerColection
		}
		аьеюЗкштедт(ыекуфьЕшслукСщддусешщты)
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
	streamTicker.ReferenceCurrency = ticker.ReferenceCurrency
	streamTicker.TargetCurrency = ticker.TargetCurrency
	return streamTicker
}