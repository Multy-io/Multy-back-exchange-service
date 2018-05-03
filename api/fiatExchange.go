package api


import (
	//"fmt"
	//"strconv"
	//"encoding/json"

	//"fmt"
	//"fmt"
	"github.com/Appscrunch/Multy-back-exchange-service/currencies"
)




type FiatApi struct {
	*RestApi
}

func NewFiatApi() *FiatApi {
	return &FiatApi{NewRestApi()}
}


func (p *FiatApi) GetTicker(pair currencies.CurrencyPair, responseCh chan <- RestApiReposponse, errorCh chan <- error)  {

	referenceCurrencyCode := pair.ReferenceCurrency.CurrencyCode()
	targetCurrencyCode := pair.TargetCurrency.CurrencyCode()

	if  referenceCurrencyCode == "USDT" {
		referenceCurrencyCode = "USD"
	}
//fmt.Printf(pair.Symbol())
//http://free.currencyconverterapi.com/api/v5/convert?q=KRW_USD&compact=y
 	urlStrging := "http://free.currencyconverterapi.com/api/v5/convert?q="+targetCurrencyCode+"_"+referenceCurrencyCode+"&compact=y"
	//fmt.Println(urlStrging)
	p.publicRequest(urlStrging, pair, responseCh, errorCh)


}


