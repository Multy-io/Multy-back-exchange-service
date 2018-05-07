package api


import (
	//"fmt"
	//"strconv"
	//"encoding/json"

	//"fmt"
	//"fmt"
	"github.com/Appscrunch/Multy-back-exchange-service/currencies"
)




type OkexRestApi struct {
	*RestApi
}

func NewOkexRestApi() *OkexRestApi {
	return &OkexRestApi{NewRestApi()}
}


func (p *OkexRestApi) GetTicker(pair currencies.CurrencyPair, responseCh chan <- RestApiReposponse, errorCh chan <- error)  {

	referenceCurrencyCode := pair.ReferenceCurrency.CurrencyCode()
	targetCurrencyCode := pair.TargetCurrency.CurrencyCode()

	//if  referenceCurrencyCode == "BCH" {
	//	referenceCurrencyCode = "BCC"
	//} else if targetCurrencyCode == "BCH" {
	//	targetCurrencyCode = "BCC"
	//}

	//https://www.okex.com/api/v1/ticker.do?symbol=ltc_btc
	urlStrging := "https://www.okex.com/api/v1/ticker.do?symbol=" + targetCurrencyCode + "_" + referenceCurrencyCode
	//fmt.Println(urlStrging)
	p.publicRequest(urlStrging, pair, responseCh, errorCh)


}

