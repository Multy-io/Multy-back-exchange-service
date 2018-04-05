package api

import (
	"net/url"
	"log"
	"github.com/gorilla/websocket"
	//"fmt"
	"fmt"
	"encoding/json"
)


const okexHost = "real.okex.com:10440"
const okexPath = "/websocket/okexapi"

type OkexApi struct {
	connection *websocket.Conn
}

type OkexSubscription struct {
	Event   string `json:"event"`
	Channel string `json:"channel"`
}


func (b *OkexApi)  StartListen(apiCurrenciesConfiguration ApiCurrenciesConfiguration, callback func(message []byte, error error)) {
	url := url.URL{Scheme: "wss", Host: okexHost, Path: ""}
	log.Printf("connecting to %s", url.String())

	connection, _, error := websocket.DefaultDialer.Dial(url.String(), nil)

	if error != nil {
		fmt.Println("Okex ws error: ",error)
		callback(nil, error)
	} else if connection != nil {
		fmt.Println("Okex ws connected")

		productsIds :=  b.composeSymbolsForSubscirbe(apiCurrenciesConfiguration)

		for _, productId := range  productsIds {
			subscribtion := OkexSubscription{}
			subscribtion.Event = "addChannel"
			subscribtion.Channel = productId

			msg, _ := json.Marshal(subscribtion)
			connection.WriteMessage(websocket.TextMessage, msg)
		}

		for {
			func() {
				_, message, _ := connection.ReadMessage()
				//fmt.Printf("%s \n", message)
				callback(message, error)
			}()
		}
	} else {
		fmt.Println("connection is nil")
		callback(nil, nil)
	}


}

func (b *OkexApi)  StopListen() {
	//fmt.Println("before close")
	//b.connection.Close()
	//fmt.Println("closed")
}

func (b *OkexApi)  composeSymbolsForSubscirbe(apiCurrenciesConfiguration ApiCurrenciesConfiguration) []string {
	var smybolsForSubscirbe = []string{}
	for _, targetCurrency := range apiCurrenciesConfiguration.TargetCurrencies {

		symbol := "ok_sub_futureusd_" +  targetCurrency + "_ticker_this_week"
		smybolsForSubscirbe = append(smybolsForSubscirbe, symbol)

	}
	return smybolsForSubscirbe

}
