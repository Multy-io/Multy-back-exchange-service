package api


import (
	//"fmt"
	//"strconv"
	//"encoding/json"

	//"fmt"
	//"fmt"
	"github.com/Appscrunch/Multy-back-exchange-service/currencies"
)




type BithumbApi struct {
	*RestApi
}

func NewBithumbApi() *BithumbApi {
	return &BithumbApi{NewRestApi()}
}


func (p *BithumbApi) GetTicker(pair currencies.CurrencyPair, responseCh chan <- RestApiReposponse, errorCh chan <- error)  {

// https://api.bithumb.com/public/ticker/all
 urlStrging := "https://api.bithumb.com/public/ticker/all"
	//fmt.Println(urlStrging)
	p.publicRequest(urlStrging, pair, responseCh, errorCh)


}


