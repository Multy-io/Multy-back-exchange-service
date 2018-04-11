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
	"strings"
	//"fmt"
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
	//manger.allTickers
	manger.grpcClient = NewGrpcClient()
	manger.fiatureCh = make(chan *server.Tickers)
	manger.dbManger = NewDbManager()

	//manger.binanceManager = NewBinanceManager()
	//manger.hitBtcManager = &HitBtcManager{}
	//manger.poloniexManager = &PoloniexManager{}
	//manger.bitfinexManager = &BitfinexManager{}
	//manger.gdaxManager = &GdaxManager{}
	//manger.okexManager = &OkexManager{}
	//manger.server = &stream.Server{}
	//manger.agregator = NewAgregator()
	return &manger
}

func (b *ExchangeManager) StartGetingData() {

	go b.grpcClient.printAllTickers(b.fiatureCh)
	go b.fillDb()

	for {
		select {
		case msg := <-b.fiatureCh:
			//fmt.Println("received message", msg)
			b.add(msg)
		default:
			//fmt.Println("no activity")
		}
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
	//fmt.Println(b.exchanges["BINANCE"])
	b.getRates(time.Now(), "Binance", "ETH",[]string{"USDT", "BTC"})
	b.Unlock()
}

func (b *ExchangeManager) getRates(timeStamp time.Time, exchangeName string, targetCode string, refereceCodes []string)  []*Ticker {
	exchange := b.exchanges[strings.ToUpper(exchangeName)]

	if exchange == nil {
		return nil
	}

	//fmt.Println(exchange.name)
	var tickers = []*Ticker{}
	for _, refereceCode := range refereceCodes {

		var symbol = targetCode + "-" + refereceCode

		//fmt.Println(symbol)

		var ticker = exchange.Tickers[symbol]
		//fmt.Println(ticker)
		if ticker != nil {
			tickers = append(tickers, ticker)
		}
	}


	//fmt.Println(tickers)
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
	}
}