package main

import (
	"sync"

	"time"

	"github.com/Appscrunch/Multy-back-exchange-service/core"
	"github.com/Appscrunch/Multy-back-exchange-service/exchange-rates"
	"github.com/KristinaEtc/config"
	_ "github.com/KristinaEtc/slflog"
)

var manager = core.NewManager()
var exchangeManger *exchangeRates.ExchangeManager
var waitGroup = &sync.WaitGroup{}
var configuration = core.ManagerConfiguration{}

//var configString = `{
//		"targetCurrencies" : ["BTC", "ETH", "GOLOS", "BTS", "STEEM", "WAVES", "LTC", "BCH", "ETC", "DASH", "EOS"],
//		"referenceCurrencies" : ["USD", "BTC"],
//		"exchanges": ["Binance","Bitfinex","Gdax","HitBtc","Okex","Poloniex"],
//		"refreshInterval" : "3"
//		}`

//const (
//	DbUser     = "postgres"
//	DbPassword = "postgres"
//	DbName     = "test"
//)

func main() {

	config.ReadGlobalConfig(&configuration, "multy exchange info service configuration")

	//, "GOLOS", "BTS", "STEEM", "WAVES", "LTC", "BCH", "ETC", "DASH", "EOS"
	//	[]string{"BTC", "ETH", "GOLOS", "BTS", "STEEM", "WAVES", "LTC", "BCH", "ETC", "DASH", "EOS"}
	//	configuration.TargetCurrencies = []string{"LTC","DASH"}
	configuration.TargetCurrencies = []string{"BTC", "ETH", "GOLOS", "BTS", "STEEM", "WAVES", "LTC", "BCH", "ETC", "DASH", "EOS"}
	configuration.ReferenceCurrencies = []string{"USDT", "BTC"}
	configuration.Exchanges = []string{"BINANCE", "BITFINEX"}
	configuration.RefreshInterval = 1

	configuration.HistoryApiKey = "A502B3C1-9C40-446F-9831-CA12EC039AB8"
	historyStartDate, _ := time.Parse(
		time.RFC3339,
		"2016-11-01T22:08:41+00:00")
	configuration.HistoryStartDate = historyStartDate
	configuration.HistoryEndDate = time.Now().UTC().Add(-3600)

	dbConfig := core.DBConfiguration{}
	// dbConfig.User = "postgres"
	// dbConfig.Password = "postgres"
	// dbConfig.Name = "test"
	configuration.DBConfiguration = dbConfig

	waitGroup.Add(len(configuration.Exchanges) + 5)

	go manager.StartListen(configuration)

	exchangeManger = exchangeRates.NewExchangeManager(configuration)
	go exchangeManger.StartGetingData()

	waitGroup.Wait()

}
