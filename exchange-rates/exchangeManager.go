package exchangeRates

import (
	//"strconv"
	"time"

	"sync"

	"github.com/Appscrunch/Multy-back-exchange-service/core"
	"github.com/Appscrunch/Multy-back-exchange-service/currencies"
	"github.com/Appscrunch/Multy-back-exchange-service/stream/server"
	"fmt"
)

type Exchange struct {
	name    string
	Tickers map[string]Ticker
	//StraightPares
}

func (b *Exchange) containPair(pair currencies.CurrencyPair) bool {

	for _, ticker := range b.Tickers {
		if ticker.Pair.IsEqualTo(pair) {
			return true
		}
	}

	return false
}

func (b *Exchange) tickerForPair(pair currencies.CurrencyPair) *Ticker {
	for _, ticker := range b.Tickers {
		//fmt.Println(ticker.TargetCurrency.CurrencyCode(), ticker.ReferenceCurrency.CurrencyCode())
		if ticker.Pair.IsEqualTo(pair) {
			//fmt.Println("sdd", ticker.symbol())
			return &ticker
		}
	}
	return nil
}



type Ticker struct {
	//TargetCurrency    currencies.Currency
	//ReferenceCurrency currencies.Currency
	Pair         currencies.CurrencyPair
	Rate         float64
	TimpeStamp   time.Time
	isCalculated bool
}

func (b *Ticker) symbol() string {
	return b.Pair.TargetCurrency.CurrencyCode() + "-" + b.Pair.ReferenceCurrency.CurrencyCode()
}





type ExchangeManager struct {
	exchanges           map[string]Exchange
	grpcClient          *GrpcClient
	tickersCh           chan *server.Tickers
	dbManger            *DbManager
	referenceCurrencies []currencies.Currency
	configuration       core.ManagerConfiguration
	sync.Mutex
	historyManager *core.HistoryManager

}

func NewExchangeManager(configuration core.ManagerConfiguration) *ExchangeManager {
	var manger = ExchangeManager{}
	manger.configuration = configuration
	manger.exchanges = map[string]Exchange{}
	manger.grpcClient = NewGrpcClient()
	manger.tickersCh = make(chan *server.Tickers)

	dbConfig := DBConfiguration{}
	dbConfig.Name = configuration.DBConfiguration.Name
	dbConfig.Password = configuration.DBConfiguration.Password
	dbConfig.User = configuration.DBConfiguration.User

	manger.dbManger = NewDbManager(dbConfig)


	manger.historyManager = &core.HistoryManager{}


	return &manger
}

func (b *ExchangeManager) StartGetingData() {

	responseCh := make(chan core.HistoryResponse)

	go b.grpcClient.listenTickers(b.tickersCh)
	go b.fillDb()

	var exchanges []core.Exchange
	for _, exchangeString := range b.configuration.Exchanges {
		ex := core.NewExchange(exchangeString)
		exchanges = append(exchanges, ex)
	}



	//historyConfiguration := core.HistoryConfiguration{}
	//historyConfiguration.Exchanges = exchanges
	//historyConfiguration.HistoryStartDate = b.configuration.HistoryStartDate
	//historyConfiguration.HistoryEndDate = b.configuration.HistoryEndDate
	//historyConfiguration.Pairs = b.configuration.Pairs()
	//historyConfiguration.ApiKey = b.configuration.HistoryApiKey
	//go b.historyManager.StartCollectHistory(historyConfiguration, responseCh)

	//ch := make(chan []*Exchange)
	//go b.Subscribe(ch, 5, []string{"DASH", "ETC", "EOS", "WAVES", "STEEM", "BTS", "ETH"}, "USDT")

	for {
		select {
		case msg := <-b.tickersCh:
			//fmt.Println("received message", msg)
			b.add(*msg)
		//case ex := <-ch:
		//
		//	for _, exx := range ex {
		//	}	//fmt.Println("received ex", exx.name, exx.Tickers)
		//	for _, v := range exx.Tickers {
		//		if v.isCalculated {
		//			fmt.Println(exx.name, v.symbol(), v.Rate)
		//		}
		//	}

		case response := <-responseCh:
			fmt.Printf("got history %@,  %@ %@ %s \n", response.Exchange.String(), response.Pair.TargetCurrency.CurrencyCode(), response.Pair.ReferenceCurrency.CurrencyCode(), response.OhlcvData)
			if response.OhlcvData != nil {
				b.addHistoryData(response)
			}
		}
	}
}

func (b *ExchangeManager) Subscribe(ch chan []*Exchange, refreshInterval time.Duration, targetCodes []string, referenceCode string) {
	go func() {
		for range time.Tick(refreshInterval * time.Second) {
			var newExchanges = b.calculateAllTickers(targetCodes, referenceCode)
			ch <- newExchanges
		}
	}()
}

func (b *ExchangeManager) calculateAllTickers(targetCodes []string, referenceCode string) []*Exchange {

	referenceCurrency := currencies.NewCurrencyWithCode(referenceCode)
	var referenceCrossCurrency currencies.Currency
	if referenceCurrency == currencies.Bitcoin {
		referenceCrossCurrency = currencies.Tether
	} else {
		referenceCrossCurrency = currencies.Bitcoin
	}

	//fmt.Println(referenceCurrency.CurrencyCode(), referenceCrossCurrency.CurrencyCode())

	var newExchanges = []*Exchange{}

	for _, targetCode := range targetCodes {
		if targetCode == referenceCode {
			continue
		}
		var pair = currencies.CurrencyPair{}
		pair.TargetCurrency = currencies.NewCurrencyWithCode(targetCode)
		pair.ReferenceCurrency = referenceCurrency

		//fmt.Println("pair is:",pair.TargetCurrency.CurrencyCode(), pair.ReferenceCurrency.CurrencyCode())

		b.Lock()
		exchanges := map[string]Exchange{}
		for k, v := range b.exchanges {
			exchanges[k] = v
		}

		for _, exchange := range exchanges {

			var newTickers = map[string]Ticker{}

			if ticker := exchange.tickerForPair(pair); ticker != nil {
				//fmt.Println("tikers si not nil:", exchange.name, ticker.symbol())
				newTickers[ticker.symbol()] = *ticker
			} else {
				//fmt.Println("tiker is nil", exchange.name)
				crossPair := pair
				crossPair.ReferenceCurrency = referenceCrossCurrency
				//fmt.Println(crossPair.TargetCurrency.CurrencyCode(), crossPair.ReferenceCurrency.CurrencyCode())
				if crossTicker := exchange.tickerForPair(crossPair); crossTicker == nil {
					//fmt.Println("crossTiker is nil", exchange.name)
					continue
				} else {
					exchangePair := currencies.CurrencyPair{}
					isStreight := false
					if pair.ReferenceCurrency == currencies.Tether {
						exchangePair.TargetCurrency = crossPair.ReferenceCurrency
						exchangePair.ReferenceCurrency = pair.ReferenceCurrency
						isStreight = true
					} else if pair.ReferenceCurrency == currencies.Bitcoin {
						exchangePair.TargetCurrency = pair.ReferenceCurrency
						exchangePair.ReferenceCurrency = crossPair.ReferenceCurrency
					}
					//fmt.Println("crossTiker is", crossTicker.Pair.TargetCurrency.CurrencyCode(), crossTicker.Pair.ReferenceCurrency.CurrencyCode(), exchange.name)
					//fmt.Println(crossTicker.TargetCurrency, crossTicker.ReferenceCurrency)

					if exchangeTicker := exchange.tickerForPair(exchangePair); exchangeTicker != nil {
						var rate float64
						if isStreight {
							rate = crossTicker.Rate * exchangeTicker.Rate
						} else {
							rate = crossTicker.Rate / exchangeTicker.Rate
						}

						//fmt.Println(exchange.name, exchangeTicker.symbol())
						ticker := Ticker{}

						ticker.TimpeStamp = exchangeTicker.TimpeStamp
						ticker.Rate = rate
						ticker.Pair.TargetCurrency = pair.TargetCurrency
						ticker.Pair.ReferenceCurrency = pair.ReferenceCurrency
						ticker.isCalculated = true
						newTickers[ticker.symbol()] = ticker
					} else {
						log.Errorf("calculateAllTickers: exchange ticket is nil %v %v %v ", exchangePair.TargetCurrency.CurrencyCode(), exchangePair.ReferenceCurrency.CurrencyCode(), exchange.name)
					}
				}

			}

			if len(newTickers) > 0 {
				var newExchange = Exchange{}
				newExchange.name = exchange.name
				newExchange.Tickers = newTickers
				newExchanges = append(newExchanges, &newExchange)
			}

		}

		b.Unlock()
	}

	return newExchanges

}

func (b *ExchangeManager) add(tikers server.Tickers) {
	b.Lock()

	//for _, exchangeTicker := range tikers.ExchangeTickers {
	//	if _, ok := b.exchanges[exchangeTicker.Exchange]; !ok {
	//		var ex = Exchange{}
	//		ex.name = exchangeTicker.Exchange
	//		b.exchanges[exchangeTicker.Exchange] = ex
	//	}
	//
	//	for _, value := range exchangeTicker.Tickers {
	//		var ticker = Ticker{}
	//		ticker.TimpeStamp = time.Now()
	//		ticker.Pair.TargetCurrency = currencies.NewCurrencyWithCode(value.Target)
	//		ticker.Pair.ReferenceCurrency = currencies.NewCurrencyWithCode(value.Referrence)
	//		ticker.Rate, _ = strconv.ParseFloat(value.Rate, 64)
	//
	//		if v, ok := b.exchanges[exchangeTicker.Exchange]; ok {
	//			if v.Tickers == nil {
	//				v.Tickers =  map[string]Ticker{}
	//				b.exchanges[exchangeTicker.Exchange] = v
	//			}
	//		}
	//		b.exchanges[exchangeTicker.Exchange].Tickers[ticker.symbol()] = ticker
	//	}
	//
	//}

	b.Unlock()
}

func (b *ExchangeManager) addHistoryData(historyData core.HistoryResponse) {
	b.Lock()

		if _, ok := b.exchanges[historyData.Exchange.String()]; !ok {
			var ex = Exchange{}
			ex.name = historyData.Exchange.String()
			b.exchanges[historyData.Exchange.String()] = ex
		}

		for _, value := range historyData.OhlcvData {
			var ticker = Ticker{}
			ticker.TimpeStamp = value.Time_open
			ticker.Pair = historyData.Pair
			ticker.Rate, _ = value.Price_open.Float64()

			if v, ok := b.exchanges[historyData.Exchange.String()]; ok {
				if v.Tickers == nil {
					v.Tickers =  map[string]Ticker{}
					b.exchanges[historyData.Exchange.String()] = v
				}
			}
			fmt.Println(ticker.symbol()+value.Time_open.String())
			b.exchanges[historyData.Exchange.String()].Tickers[ticker.symbol()+value.Time_open.String()] = ticker
		}


	b.Unlock()
}


func (b *ExchangeManager) GetRates(timeStamp time.Time, exchangeName string, targetCode string, referecies []string) []*Ticker {

	var dbRates = b.dbManger.getRates(timeStamp, exchangeName, targetCode, referecies)

	var tickers = []*Ticker{}

	for _, dbRate := range dbRates {

		var ticker = Ticker{}
		ticker.Pair.TargetCurrency = currencies.NewCurrencyWithCode(dbRate.targetCode)
		ticker.Pair.ReferenceCurrency = currencies.NewCurrencyWithCode(dbRate.referenceCode)
		ticker.TimpeStamp = dbRate.timeStamp
		ticker.Rate = dbRate.rate
		tickers = append(tickers, &ticker)
	}
	return tickers
}

func (b *ExchangeManager) fillDb() {

	for range time.Tick(10 * time.Second) {

		dbExchanges := []*DbExchange{}

		for _, referenceCode := range b.configuration.ReferenceCurrencies {

			var newExchanges = b.calculateAllTickers(b.configuration.TargetCurrencies, referenceCode)

			//for _, ex := range newExchanges {
			//	fmt.Println(ex.name)
			//	for _, t := range ex.Tickers {
			//		fmt.Println(t.symbol(), t.Rate)
			//	}
			//}
			//fmt.Println("__________")

			for _, value := range newExchanges {
				dbExchange := DbExchange{}
				dbExchange.name = value.name
				dbExchange.Tickers = []DbTicker{}

				for _, ticker := range value.Tickers {
					dbTicker := DbTicker{}
					dbTicker.TimpeStamp = ticker.TimpeStamp
					dbTicker.ReferenceCurrency = ticker.Pair.ReferenceCurrency
					dbTicker.TargetCurrency = ticker.Pair.TargetCurrency
					dbTicker.Rate = ticker.Rate
					dbTicker.isCalculated = ticker.isCalculated
					dbExchange.Tickers = append(dbExchange.Tickers, dbTicker)
					fmt.Println(dbTicker.TargetCurrency.CurrencyCode(), dbTicker.Rate)
				}
				dbExchanges = append(dbExchanges, &dbExchange)
			}

			//fmt.Println(dbExchanges)

		}
		b.dbManger.FillDb(dbExchanges)
		b.exchanges = map[string]Exchange{}

		//v := b.GetRates(time.Now().Add(-4 * time.Minute), "BINANCE", "BTS", []string{"BTC", "USDT"})
		//
		//for _,value := range v {
		//	fmt.Println("get rates :", value.symbol(), value.TimpeStamp, value.Rate)
		//}
	}
}

func contains(currienciesCodes []string, currienceCode string) bool {
	for _, a := range currienciesCodes {
		if a == currienceCode {
			return true
		}
	}
	return false
}
