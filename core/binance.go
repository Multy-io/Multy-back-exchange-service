package core

import (
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/Appscrunch/Multy-back-exchange-service/api"
	"github.com/Appscrunch/Multy-back-exchange-service/currencies"
)

type BinanceTicker struct {
	Symbol             string  `json:"s"`
	Rate               string  `json:"c"`
	EventTime          float64 `json:"E"` // field is not needed but it's a workaround because unmarshal is case insensitive and without this filed json can't be parsed
	StatisticCloseTime float64 `json:"C"` // field is not needed but it's a workaround because unmarshal is case insensitive and without this filed json can't be parsed
}

func (b *BinanceTicker) getCurriences() (currencies.Currency, currencies.Currency) {

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
				return targetCurrency, referenceCurrency
			}
		}

	}
	return currencies.NotAplicable, currencies.NotAplicable
}

type BinanceManager struct {
	binanceApi     *api.BinanceApi
	symbolsToParse map[string]bool
}

func NewBinanceManager() *BinanceManager {
	var manger = BinanceManager{}
	manger.symbolsToParse = map[string]bool{}
	manger.binanceApi = &api.BinanceApi{}
	return &manger
}

func (b *BinanceManager) StartListen(exchangeConfiguration ExchangeConfiguration, callback func(tickerCollection TickerCollection, err error)) {
	b.symbolsToParse = b.composeSybolsToParse(exchangeConfiguration)
	b.binanceApi.StartListen(func(message []byte, err error) {
		if err != nil {
			log.Println("binance error:", err)
			callback(TickerCollection{}, err)
		} else if message != nil {
			//fmt.Printf("%s", message)
			var binanceTickers []BinanceTicker
			json.Unmarshal(message, &binanceTickers)

			var tickers = []Ticker{}

			for _, binanceTicker := range binanceTickers {
				if b.symbolsToParse[binanceTicker.Symbol] {
					var ticker = Ticker{}
					targetCurrency, referenceCurrency := binanceTicker.getCurriences()
					ticker.Symbol = binanceTicker.Symbol
					ticker.Rate = binanceTicker.Rate
					ticker.TargetCurrency = targetCurrency
					ticker.ReferenceCurrency = referenceCurrency
					tickers = append(tickers, ticker)
					//fmt.Println(binanceTicker.Symbol ,targetCurrency.CurrencyName(), referenceCurrency.CurrencyName())
				}
			}

			var tickerCollection TickerCollection
			tickerCollection.TimpeStamp = time.Now()
			tickerCollection.Tickers = tickers
			callback(tickerCollection, nil)
		}
	})

}

func (b *BinanceManager) composeSybolsToParse(exchangeConfiguration ExchangeConfiguration) map[string]bool {
	var symbolsToParse = map[string]bool{}
	for _, targetCurrency := range exchangeConfiguration.TargetCurrencies {
		for _, referenceCurrency := range exchangeConfiguration.ReferenceCurrencies {

			if referenceCurrency == "USD" {
				referenceCurrency = "USDT"
			}

			symbol := targetCurrency + referenceCurrency
			symbolsToParse[symbol] = true
		}
	}
	return symbolsToParse

}
