package main

import (
	core "Multy-back-exchange-service/core"
	"Multy-back-exchange-service/exchangeRates"
	"sync"
)

var manager = core.NewManager()
var exchangeManger *exchangeRates.ExchangeManager
var waitGroup = &sync.WaitGroup{}
//var configString = `{
//		"targetCurrencies" : ["BTC", "ETH", "GOLOS", "BTS", "STEEM", "WAVES", "LTC", "BCH", "ETC", "DASH", "EOS"],
//		"referenceCurrencies" : ["USD", "BTC"],
//		"exchanges": ["Binance","Bitfinex","Gdax","HitBtc","Okex","Poloniex"],
//		"refreshInterval" : "3"
//		}`

func main() {




	var configuration = core.ManagerConfiguration{}
	configuration.TargetCurrencies = []string{"BTC", "ETH", "GOLOS", "BTS", "STEEM", "WAVES", "LTC", "BCH", "ETC", "DASH", "EOS"}
	configuration.ReferenceCurrencies = []string{"USD", "BTC"}
	//configuration.Exchanges = []string{"Bitfinex"}
	configuration.Exchanges = []string{"Binance","Bitfinex","Gdax","HitBtc","Okex","Poloniex"}
	configuration.RefreshInterval = 1

	waitGroup.Add(len(configuration.Exchanges)+5)

	go manager.StartListen(configuration)

	exchangeManger = exchangeRates.NewExchangeManager()
	go exchangeManger.StartGetingData()

	waitGroup.Wait()


}

