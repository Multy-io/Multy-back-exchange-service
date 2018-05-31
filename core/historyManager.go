package core

import (
	SDK "github.com/CoinAPI/coinapi-sdk/go-rest/v1"
	"time"
	"github.com/Appscrunch/Multy-back-exchange-service/currencies"
)

type HistoryConfiguration struct {
	HistoryStartDate    time.Time
	HistoryEndDate		time.Time
	Pairs []currencies.CurrencyPair
	Exchanges           []Exchange        `json:"exchanges"`
	ApiKey string
}

type HistoryResponse struct {
	OhlcvData []SDK.Ohlcv_data
	Pair currencies.CurrencyPair
	Exchange          Exchange
}


type HistoryManager struct {
	BasicManager
	sdk *SDK.SDK
	//bittrexApi    *api.BittrexApi
}


func (b *HistoryManager) StartCollectHistory(configuration HistoryConfiguration, responseCh chan HistoryResponse) {
	b.sdk = SDK.NewSDK(configuration.ApiKey)



	duration := (time.Duration(len(configuration.Pairs)) * time.Second) + 5

	go func() {
		for _, exchange := range configuration.Exchanges {
			time.Sleep(duration)
			go b.collectHistoryForExchange(exchange, configuration, responseCh)
		}
	}()

}

func (b *HistoryManager) collectHistoryForExchange(exchange Exchange, configuration HistoryConfiguration, responseCh chan HistoryResponse) {

	go func() {
		for _, paiar := range configuration.Pairs {
			time.Sleep(1 * time.Second)
			go b.collectHistoryFor(exchange, paiar, configuration.HistoryStartDate, configuration.HistoryEndDate, responseCh)
		}
	}()
}


func (b *HistoryManager) collectHistoryFor(exchange Exchange, pair currencies.CurrencyPair, startDate time.Time, endDate time.Time, responseCh chan HistoryResponse) {

	if exchange.CoinApiString() == "" {
		return
	}

	referenceCurrencyCode := pair.ReferenceCurrency.CurrencyCode()
	if (exchange == Bitfinex || exchange == Kraken || exchange == Gdax) && referenceCurrencyCode == "USDT" {
		referenceCurrencyCode = "USD"
	}

	symbolId := exchange.CoinApiString()+"_SPOT_"+pair.TargetCurrency.CurrencyCode()+"_"+referenceCurrencyCode


	Ohlcv_historic_data_with_time_end_and_limit, _ := b.sdk.Ohlcv_historic_data_with_time_end_and_limit(symbolId, "6HRS", startDate, endDate, 1)
	//fmt.Println("Ohlcv_historic_data_with_time_end_and_limit:")
	//fmt.Println("number:", len(Ohlcv_historic_data_with_time_end_and_limit))
	//Ohlcv_historic_data_with_time_end_and_limit_item, _ := json.MarshalIndent(&Ohlcv_historic_data_with_time_end_and_limit, "", "  ")
	//fmt.Println("first items:", string(Ohlcv_historic_data_with_time_end_and_limit_item))

	historyResponse := HistoryResponse{Ohlcv_historic_data_with_time_end_and_limit, pair, exchange}

	responseCh <- historyResponse
}
