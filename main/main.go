package main

import (
	core "Multy-back-exchange-service/core"
)

var manager = core.NewManager()


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
	configuration.Exchanges = []string{"Poloniex"}
	//configuration.Exchanges = []string{"Binance","Bitfinex","Gdax","HitBtc","Okex","Poloniex"}
	configuration.RefreshInterval = "3"

	manager.StartListen(configuration)

}

