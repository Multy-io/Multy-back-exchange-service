package exchangeRates

import (
	"Multy-back-exchange-service/stream/server"
	//"fmt"
	//"fmt"
	//"fmt"
	"sync"
	"time"
	"Multy-back-exchange-service/currencies"
	//"fmt"
	//"strings"
	//"fmt"
	//"fmt"
	"strconv"

	"fmt"
)

type Exchange struct {
	name string
	Tickers map[string]*Ticker
}


type Ticker struct {
	TargetCurrency currencies.Currency
	ReferenceCurrency currencies.Currency
	Rate	string
	TimpeStamp time.Time
}

func (b *Ticker) symbol() string {
	return b.TargetCurrency.CurrencyCode() + "-" + b.ReferenceCurrency.CurrencyCode()
}


type ExchangeManager struct {
	sync.Mutex
	exchanges map[string]*Exchange
	grpcClient *GrpcClient
	fiatureCh chan *server.Tickers
	dbManger *DbManager
}



func NewExchangeManager() *ExchangeManager {
	var manger = ExchangeManager{}
	manger.exchanges = map[string]*Exchange{}
	manger.grpcClient = NewGrpcClient()
	manger.fiatureCh = make(chan *server.Tickers)
	manger.dbManger = NewDbManager()

	return &manger
}

func (b *ExchangeManager) StartGetingData() {

	go b.grpcClient.printAllTickers(b.fiatureCh)
	go b.fillDb()


	ch := make(chan []*Exchange)
	go b.subscribe(ch, 5, []string{"BTC", "ETH"}, "USDT")

	for {
		select {
		case msg := <-b.fiatureCh:
			//fmt.Println("received message", msg)
			b.add(msg)
		case ex := <-ch:


			for _, exx := range ex {
				fmt.Println("received ex", exx)
			}

		default:
			//fmt.Println("no activity")
		}
	}
}


func  (b *ExchangeManager) subscribe(ch chan []*Exchange, refreshInterval time.Duration, targetCodes  []string, referenceCode string) {

	for range time.Tick(refreshInterval * time.Second) {

		var newExchanges= []*Exchange{}
		for _, exchange := range b.exchanges {
			var newTickers= map[string]*Ticker{}
			for _, ticker := range exchange.Tickers {

				//fmt.Println(exchange.name)
				//fmt.Println(ticker.symbol())
				if ticker.ReferenceCurrency.CurrencyCode() == referenceCode && contains(targetCodes, ticker.TargetCurrency.CurrencyCode()) {
					newTickers[ticker.symbol()] = ticker
				}
			}
			if len(newTickers) > 0 {
				var newExchange= Exchange{}
				newExchange.name = exchange.name
				newExchange.Tickers = newTickers
				newExchanges = append(newExchanges, &newExchange)
			}
		}
		ch <- newExchanges
	}

}

func (b *ExchangeManager) add(tikers *server.Tickers) {
	b.Lock()
	for _, exchangeTicker := range tikers.ExchangeTickers {

		if b.exchanges[exchangeTicker.Exchange] == nil {
			var ex = Exchange{}
			ex.name = exchangeTicker.Exchange
			b.exchanges[exchangeTicker.Exchange] = &ex
		}

		for _, value := range exchangeTicker.Tickers {
			var ticker = Ticker{}
			ticker.TimpeStamp = time.Now()
			ticker.TargetCurrency = currencies.NewCurrencyWithCode(value.Target)
			ticker.ReferenceCurrency = currencies.NewCurrencyWithCode(value.Referrence)
			ticker.Rate = value.Rate

			if b.exchanges[exchangeTicker.Exchange].Tickers == nil {
				b.exchanges[exchangeTicker.Exchange].Tickers = map[string]*Ticker{}
			}
			b.exchanges[exchangeTicker.Exchange].Tickers[ticker.symbol()] = &ticker
		}
	}
	b.Unlock()
}

func (b *ExchangeManager) getRates(timeStamp time.Time, exchangeName string, targetCode string, referecies []string)  []*Ticker {

	var dbRates = b.dbManger.getRates(timeStamp, exchangeName, targetCode, referecies)

	var tickers = []*Ticker{}

	for _, dbRate := range dbRates {

		var ticker= Ticker{}
		ticker.TargetCurrency = currencies.NewCurrencyWithCode(dbRate.targetCode)
		ticker.ReferenceCurrency = currencies.NewCurrencyWithCode(dbRate.referenceCode)
		ticker.TimpeStamp = dbRate.timeStamp
		ticker.Rate = strconv.FormatFloat(dbRate.rate, 'f', 8, 64)
		tickers = append(tickers, &ticker)
	}
	return tickers
}


func (b *ExchangeManager) fillDb() {

	for range time.Tick(5 * time.Second) {

		dbExchanges := []*DbExchange{}

		for _, value := range b.exchanges {
			dbExchange := DbExchange{}
			dbExchange.name = value.name
			dbExchange.Tickers = []*DbTicker{}

			for _, ticker := range value.Tickers {
				dbTicker := DbTicker{}
				dbTicker.TimpeStamp = ticker.TimpeStamp
				dbTicker.ReferenceCurrency = ticker.ReferenceCurrency
				dbTicker.TargetCurrency = ticker.TargetCurrency
				dbTicker.Rate = ticker.Rate
				dbExchange.Tickers = append(dbExchange.Tickers, &dbTicker)
			}
			dbExchanges = append(dbExchanges, &dbExchange)
		}

		b.dbManger.FillDb(dbExchanges)

		//v := b.getRates(time.Now().Add(-4 * time.Minute), "BINANCE", "ETH", []string{"BTC", "USDT"})
		//
		//for _,value := range v {
		//	fmt.Println("rasult:", value)
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