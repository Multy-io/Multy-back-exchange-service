
package api

import (
//"fmt"
//"strconv"
//"encoding/json"

//"fmt"
//"fmt"
"github.com/Appscrunch/Multy-back-exchange-service/currencies"
 "strings"
)




type HuobiApi struct {
	*RestApi
}

func NewHuobiApi() *HuobiApi {
	return &HuobiApi{NewRestApi()}
}


func (p *HuobiApi) GetTicker(pair currencies.CurrencyPair, responseCh chan <- RestApiReposponse, errorCh chan <- error)  {

	referenceCurrencyCode := pair.ReferenceCurrency.CurrencyCode()
	targetCurrencyCode := pair.TargetCurrency.CurrencyCode()

	//if  referenceCurrencyCode == "BCH" {
	//	referenceCurrencyCode = "BCC"
	//} else if targetCurrencyCode == "BCH" {
	//	targetCurrencyCode = "BCC"
	//}

//http://api.huobipro.com/market/detail/merged?symbol=btcusdt
	urlStrging := "http://api.huobipro.com/market/detail/merged?symbol=" + strings.ToLower(targetCurrencyCode) + strings.ToLower(referenceCurrencyCode)
	//fmt.Println(urlStrging)
	p.publicRequest(urlStrging, pair, responseCh, errorCh)


}

