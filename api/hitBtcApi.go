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

func (b *HitBtcApi)  connectWs(apiCurrenciesConfiguration ApiCurrenciesConfiguration) *websocket.Conn {
	url := url.URL{Scheme: "wss", Host: hitBtcHost, Path: hitBtcPath}
	log.Printf("connecting to %s", url.String())

	connection, _, error := websocket.DefaultDialer.Dial(url.String(), nil)

	if error != nil || connection == nil {
		fmt.Println("HitBtc ws connection error: ",error)
		return nil
	} else  {
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

		return connection
	}
}



func (b *HitBtcApi)  StartListen(apiCurrenciesConfiguration ApiCurrenciesConfiguration, callback func(message []byte, error error)) {

	for {
		if b.connection == nil {
			b.connection = b.connectWs(apiCurrenciesConfiguration)
		} else if b.connection != nil {

			func() {
				_, message, error := b.connection.ReadMessage()
				if error != nil {
					fmt.Println("HitBtc read message error:", error)
					b.connection.Close()
					b.connection = nil
				} else {
					//fmt.Printf("%s \n", message)
					callback(message, error)
				}
			}()
		}
	}


}

func (b *HitBtcApi)  StopListen() {
	//fmt.Println("before close")
	if b.connection != nil {
		b.connection.Close()
		b.connection = nil
	}
	fmt.Println("HitBtc ws closed")
}

func (b *HitBtcApi)  composeSymbolsForSubscirbe(apiCurrenciesConfiguration ApiCurrenciesConfiguration) []string {
	var smybolsForSubscirbe = []string{}
	for _, targetCurrency := range apiCurrenciesConfiguration.TargetCurrencies {
		for _, referenceCurrency := range apiCurrenciesConfiguration.ReferenceCurrencies {

			if targetCurrency == referenceCurrency {
				continue
			}

			if referenceCurrency == "USDT" {
				referenceCurrency = "USD"
			}

			symbol := targetCurrency + referenceCurrency
			smybolsForSubscirbe = append(smybolsForSubscirbe, symbol)
		}
	}
	return smybolsForSubscirbe

}
