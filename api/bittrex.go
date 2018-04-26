package api

import (
	//"fmt"
	"io/ioutil"
	//"strconv"
	//"encoding/json"

	"net/http"
	"time"
	//"fmt"
	"fmt"
	"github.com/Appscrunch/Multy-back-exchange-service/currencies"
)

type RestApiReposponse struct {
	Message []byte
	Pair currencies.CurrencyPair
}


type BittrexApi struct {
	//connection *websocket.Conn

	httpClient *http.Client
}


func NewBittrexApi() *BittrexApi {
	var api = BittrexApi{}
	api.httpClient = &http.Client{Timeout: time.Second * 10}
	return &api
}

func (p *BittrexApi) publicRequest(urlString string, pair currencies.CurrencyPair, responseCh chan <- RestApiReposponse, errorCh chan <- error) {

	<-throttle

	//TODO - check if close is needed
	//defer close(responseCh)
	//defer close(errorCh)


	req, err := http.NewRequest("GET", urlString, nil)
	if err != nil {
		fmt.Println("error creating request:", err)
		errorCh <- Error(RequestError)
		return
	}

	req.Header.Add("Accept", "application/json")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		fmt.Println("error sending request:", err)
		errorCh <- Error(ConnectError)
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("error reading response:", err)
		errorCh <- err
		return
	}

	restApiResponse := RestApiReposponse{body, pair}

	responseCh <- restApiResponse
	//errorCh <- nil
}



func (p *BittrexApi) GetTicker(pair currencies.CurrencyPair, responseCh chan <- RestApiReposponse, errorCh chan <- error)  {

	referenceCurrencyCode := pair.ReferenceCurrency.CurrencyCode()
	targetCurrencyCode := pair.TargetCurrency.CurrencyCode()

	if  referenceCurrencyCode == "BCH" {
		referenceCurrencyCode = "BCC"
	} else if targetCurrencyCode == "BCH" {
		targetCurrencyCode = "BCC"
	}

	urlStrging := "https://bittrex.com/api/v1.1/public/getticker?market=" + referenceCurrencyCode +"-" + targetCurrencyCode
	//fmt.Println(urlStrging)
	p.publicRequest(urlStrging, pair, responseCh, errorCh)


}

