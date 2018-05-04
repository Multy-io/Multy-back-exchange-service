package api


import (
	//"fmt"
	//"strconv"
	//"encoding/json"

	//"fmt"
	//"fmt"
	"github.com/Appscrunch/Multy-back-exchange-service/currencies"
)




type KrakenApi struct {
	*RestApi
}

func NewKrakenApi() *KrakenApi {
	return &KrakenApi{NewRestApi()}
}


func (p *KrakenApi) GetTicker(pair currencies.CurrencyPair, responseCh chan <- RestApiReposponse, errorCh chan <- error)  {

	referenceCurrencyCode := pair.ReferenceCurrency.CurrencyCode()
	targetCurrencyCode := pair.TargetCurrency.CurrencyCode()

	if  referenceCurrencyCode == "BTC" {
		referenceCurrencyCode = "XBT"
	} else if targetCurrencyCode == "BTC" {
		targetCurrencyCode = "XBT"
	}

	if  referenceCurrencyCode == "USDT" {
		referenceCurrencyCode = "USD"
	}

	//https://api.kraken.com/0/public/Ticker?pair=XBTUSD
	urlStrging := "https://api.kraken.com/0/public/Ticker?pair="+targetCurrencyCode+referenceCurrencyCode
	//fmt.Println(urlStrging)
	p.publicRequest(urlStrging, pair, responseCh, errorCh)


}


