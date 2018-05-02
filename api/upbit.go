package api


import (
//"fmt"
//"strconv"
//"encoding/json"

//"fmt"
//"fmt"
"github.com/Appscrunch/Multy-back-exchange-service/currencies"
)




type UpbitApi struct {
	*RestApi
}

func NewUpbitApi() *UpbitApi {
	return &UpbitApi{NewRestApi()}
}


func (p *UpbitApi) GetTicker(pair currencies.CurrencyPair, responseCh chan <- RestApiReposponse, errorCh chan <- error)  {

	referenceCurrencyCode := pair.ReferenceCurrency.CurrencyCode()
	targetCurrencyCode := pair.TargetCurrency.CurrencyCode()

	//if  referenceCurrencyCode == "BCH" {
	//	referenceCurrencyCode = "BCC"
	//} else if targetCurrencyCode == "BCH" {
	//	targetCurrencyCode = "BCC"
	//}

	//https://crix-api.upbit.com/v1/crix/trades/ticks?code=CRIX.UPBIT.USDT-BTC&count=1
	urlStrging := "https://crix-api.upbit.com/v1/crix/trades/ticks?code=CRIX.UPBIT."+referenceCurrencyCode+"-"+targetCurrencyCode+"&count=1"
	//fmt.Println(urlStrging)
	p.publicRequest(urlStrging, pair, responseCh, errorCh)


}


