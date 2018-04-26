package core

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/Appscrunch/Multy-back-exchange-service/api"
	"github.com/Appscrunch/Multy-back-exchange-service/currencies"
)

type BitfinexManager struct {
	BasicManager
	bitfinexTickers map[int]BitfinexTicker
	api             *api.BitfinexApi
}

type BitfinexTicker struct {
	ChanID     int    `json:"chanId"`
	Channel    string `json:"channel"`
	Event      string `json:"event"`
	Pair       string `json:"pair"`
	Symbol     string `json:"symbol"`
	Rate       string
	TimpeStamp time.Time
}

func (bitfinexTicker BitfinexTicker) IsFilled() bool {
	return (len(bitfinexTicker.Symbol) > 0 && len(bitfinexTicker.Rate) > 0)
}

func (b *BitfinexTicker) getCurriences() currencies.CurrencyPair {
	if len(b.Symbol) > 0 {
		var symbol = b.Symbol
		var damagedSymbol = TrimLeftChars(symbol, 2)
		for _, referenceCurrency := range currencies.DefaultReferenceCurrencies {
			//fmt.Println(damagedSymbol, referenceCurrency.CurrencyCode())

			referenceCurrencyCode := referenceCurrency.CurrencyCode()

			if referenceCurrencyCode == "USDT" {
				referenceCurrencyCode = "USD"
			}

			//fmt.Println(damagedSymbol)
			//fmt.Println(referenceCurrencyCode)
			//fmt.Println(strings.Contains(damagedSymbol, referenceCurrencyCode))

			if strings.Contains(damagedSymbol, referenceCurrencyCode) {
				//fmt.Println(damagedSymbol)

				//fmt.Println("2",symbol, referenceCurrency.CurrencyCode())
				targetCurrencyStringWithT := strings.TrimSuffix(symbol, referenceCurrencyCode)
				targetCurrencyString := TrimLeftChars(targetCurrencyStringWithT, 1)
				//fmt.Println("targetCurrencyString", targetCurrencyString)
				var targetCurrency = currencies.NewCurrencyWithCode(targetCurrencyString)
				return currencies.CurrencyPair{ targetCurrency, referenceCurrency}
			}
		}

	}
	return currencies.CurrencyPair{currencies.NotAplicable, currencies.NotAplicable}
}

func (b *BitfinexManager) StartListen(exchangeConfiguration ExchangeConfiguration, resultChan chan Result) {
	b.bitfinexTickers = make(map[int]BitfinexTicker)
	b.api = api.NewBitfinexApi()

	var apiCurrenciesConfiguration = api.ApiCurrenciesConfiguration{}
	apiCurrenciesConfiguration.TargetCurrencies = exchangeConfiguration.TargetCurrencies
	apiCurrenciesConfiguration.ReferenceCurrencies = exchangeConfiguration.ReferenceCurrencies

	ch := make(chan api.Reposponse)

	go b.api.StartListen(apiCurrenciesConfiguration, ch)

	go b.startSendingDataBack(exchangeConfiguration, resultChan)

	for {
		select {
		case response := <-ch:

			//fmt.Println(0)
			if *response.Err != nil {
				log.Errorf("StartListen *response.Err: %v", response.Err)
				resultChan <- Result{exchangeConfiguration.Exchange.String(), nil, response.Err}
			} else if *response.Message != nil {
				//fmt.Printf("%s \n", response.Message)
				//fmt.Println(1)
				b.addMessage(*response.Message)
				//fmt.
			} else {
				log.Errorf("StartListen :error parsing Bitfinex ticker")
			}

		}
	}

}

func (b *BitfinexManager) startSendingDataBack(exchangeConfiguration ExchangeConfiguration, resultChan chan Result) {

	for range time.Tick(1 * time.Second) {
		func() {
			tickers := []Ticker{}

			b.Lock()
			tempTickers := map[int]BitfinexTicker{}
			for k, v := range b.bitfinexTickers {
				tempTickers[k] = v
			}
			b.Unlock()


			for _, value := range tempTickers {
				if value.TimpeStamp.After(time.Now().Add(-maxTickerAge * time.Second)) {
					var ticker = Ticker{}
					ticker.Rate, _ = strconv.ParseFloat(value.Rate, 64)

					ticker.Pair = value.getCurriences()
					tickers = append(tickers, ticker)
				}
			}

			var tickerCollection = TickerCollection{}
			tickerCollection.TimpeStamp = time.Now()
			tickerCollection.Tickers = tickers
			//fmt.Println(tickerCollection)
			if len(tickerCollection.Tickers) > 0 {
				resultChan <- Result{exchangeConfiguration.Exchange.String(), &tickerCollection, nil}
			}
		}()
	}
}

func (b *BitfinexManager) addMessage(message []byte) {

	var bitfinexTicker BitfinexTicker
	json.Unmarshal(message, &bitfinexTicker)

	if bitfinexTicker.ChanID > 0 {
		//fmt.Println(bitfinexTicker)
		b.Lock()
		b.bitfinexTickers[bitfinexTicker.ChanID] = bitfinexTicker
		b.Unlock()
	} else {
		var unmarshaledTickerMessage []interface{}
		json.Unmarshal(message, &unmarshaledTickerMessage)
		if len(unmarshaledTickerMessage) > 1 {
			var chanId = int(unmarshaledTickerMessage[0].(float64))
			//var unmarshaledTicker []interface{}
			if v, ok := unmarshaledTickerMessage[1].([]interface{}); ok {
				b.Lock()
				var sub = b.bitfinexTickers[chanId]
				sub.Rate = strconv.FormatFloat(v[0].(float64), 'f', 8, 64)
				sub.TimpeStamp = time.Now()
				b.bitfinexTickers[chanId] = sub
				b.Unlock()
			}
		}
	}

}

//func (b PoloniexManager) convertArgsToTicker(args []interface{}) (wsticker PoloniexTicker, err error) {
//	wsticker.CurrencyPair = b.channelsByID[strconv.FormatFloat(args[0].(float64), 'f', 0, 64)]
//	wsticker.Last = args[1].(string)
//	return
//}
