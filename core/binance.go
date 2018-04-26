package core

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/Appscrunch/Multy-back-exchange-service/api"
	"github.com/Appscrunch/Multy-back-exchange-service/currencies"
	//"fmt"
	"strconv"
)

type BinanceTicker struct {
	Symbol             string  `json:"s"`
	Rate               string  `json:"c"`
	EventTime          float64 `json:"E"` // field is not needed but it's a workaround because unmarshal is case insensitive and without this filed json can't be parsed
	StatisticCloseTime float64 `json:"C"` // field is not needed but it's a workaround because unmarshal is case insensitive and without this filed json can't be parsed
}

func (b *BinanceTicker) getCurriences() currencies.CurrencyPair {

	if len(b.Symbol) > 0 {
		var symbol = b.Symbol
		var damagedSymbol = TrimLeftChars(symbol, 1)
		for _, referenceCurrency := range currencies.DefaultReferenceCurrencies {
			//fmt.Println(damagedSymbol, referenceCurrency.CurrencyCode())

			if strings.Contains(damagedSymbol, referenceCurrency.CurrencyCode()) {

				//fmt.Println("2",symbol, referenceCurrency.CurrencyCode())
				targetCurrencyString := strings.TrimSuffix(symbol, referenceCurrency.CurrencyCode())
				//fmt.Println(targetCurrencyString)
				var targetCurrency = currencies.NewCurrencyWithCode(targetCurrencyString)
				return currencies.CurrencyPair{targetCurrency, referenceCurrency}
			}
		}

	}
	return currencies.CurrencyPair{currencies.NotAplicable, currencies.NotAplicable}
}

type BinanceManager struct {
	BasicManager
	binanceApi     *api.BinanceApi
	symbolsToParse map[string]bool
}

func NewBinanceManager() *BinanceManager {
	var manger = BinanceManager{}
	manger.symbolsToParse = map[string]bool{}
	manger.binanceApi = &api.BinanceApi{}
	return &manger
}

func (b *BinanceManager) StartListen(exchangeConfiguration ExchangeConfiguration, resultChan chan Result) {
	log.Debugf("StartListen:start binance manager listen")
	b.symbolsToParse = b.composeSybolsToParse(exchangeConfiguration)
	ch := make(chan api.Reposponse)
	go b.binanceApi.StartListen(ch)

	for {
		select {
		case response := <-ch:

			if *response.Err != nil {
				log.Errorf("StartListen: binance error:%v", *response.Err)
				resultChan <- Result{exchangeConfiguration.Exchange.String(), nil, response.Err}
			} else if *response.Message != nil {
				//fmt.Printf("%s", message)
				var binanceTickers []BinanceTicker
				json.Unmarshal(*response.Message, &binanceTickers)

				var tickers = []Ticker{}

				for _, binanceTicker := range binanceTickers {
					if b.symbolsToParse[binanceTicker.Symbol] {
						var ticker = Ticker{}
						ticker.Rate, _ = strconv.ParseFloat(binanceTicker.Rate, 64)
						ticker.Pair = binanceTicker.getCurriences()
						tickers = append(tickers, ticker)
						//fmt.Println(binanceTicker.Symbol ,targetCurrency.CurrencyName(), referenceCurrency.CurrencyName())
					}
				}

				var tickerCollection TickerCollection
				tickerCollection.TimpeStamp = time.Now()
				tickerCollection.Tickers = tickers
				resultChan <- Result{exchangeConfiguration.Exchange.String(), &tickerCollection, nil}
			} else {
				log.Errorf("StartListen: Binance mesage is nil")
			}
		}
	}

}

func (b *BinanceManager) composeSybolsToParse(exchangeConfiguration ExchangeConfiguration) map[string]bool {
	var symbolsToParse = map[string]bool{}
	for _, targetCurrency := range exchangeConfiguration.TargetCurrencies {
		for _, referenceCurrency := range exchangeConfiguration.ReferenceCurrencies {

			if referenceCurrency == "USD" {
				referenceCurrency = "USDT"
			} else if referenceCurrency == targetCurrency {
				//fmt.Println(referenceCurrency, targetCurrency)
				continue
			}

			symbol := targetCurrency + referenceCurrency
			symbolsToParse[symbol] = true
		}
	}
	return symbolsToParse

}
