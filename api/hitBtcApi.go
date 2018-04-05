package api

import (
	"net/url"
	"log"
	"github.com/gorilla/websocket"
	//"fmt"
	"fmt"
	"encoding/json"
)


const hitBtcHost = "api.hitbtc.com"
const hitBtcPath = "/api/2/ws"

type HitBtcApi struct {
	connection *websocket.Conn
}

type HitBtcSubscription struct {
	Method string `json:"method"`
	Params HitBtcSubscriptionParams `json:"params"`
	ID int `json:"id"`
}

type HitBtcSubscriptionParams struct {
	Symbol string `json:"symbol"`
}



func (b *HitBtcApi)  StartListen(apiCurrenciesConfiguration ApiCurrenciesConfiguration, callback func(message []byte, error error)) {
	url := url.URL{Scheme: "wss", Host: hitBtcHost, Path: hitBtcPath}
	log.Printf("connecting to %s", url.String())

	connection, _, error := websocket.DefaultDialer.Dial(url.String(), nil)

	if error != nil {
		fmt.Println("HitBtc ws error: ",error)
		callback(nil, error)
	} else if connection != nil {
		fmt.Println("HitBtc ws connected")


		productsIds :=  b.composeSymbolsForSubscirbe(apiCurrenciesConfiguration)

		for _, productId := range  productsIds {

			hitBtcSubscriptionParams := HitBtcSubscriptionParams{}
			hitBtcSubscriptionParams.Symbol = productId

			subscribtion := HitBtcSubscription{}
			subscribtion.Method = "subscribeTicker"
			subscribtion.ID = 10000
			subscribtion.Params = hitBtcSubscriptionParams

			msg, _ := json.Marshal(subscribtion)
			connection.WriteMessage(websocket.TextMessage, msg)
		}

		for {
			func() {
				_, message, error := connection.ReadMessage()
				callback(message, error)
			}()
		}
	} else {
		fmt.Println("connection is nil")
		callback(nil, nil)
	}


}

func (b *HitBtcApi)  StopListen() {
	//fmt.Println("before close")
	//b.connection.Close()
	//fmt.Println("closed")
}

func (b *HitBtcApi)  composeSymbolsForSubscirbe(apiCurrenciesConfiguration ApiCurrenciesConfiguration) []string {
	var smybolsForSubscirbe = []string{}
	for _, targetCurrency := range apiCurrenciesConfiguration.TargetCurrencies {
		for _, referenceCurrency := range apiCurrenciesConfiguration.ReferenceCurrencies {

			if targetCurrency == referenceCurrency {
				continue
			}

			symbol := targetCurrency + referenceCurrency
			smybolsForSubscirbe = append(smybolsForSubscirbe, symbol)
		}
	}
	return smybolsForSubscirbe

}
