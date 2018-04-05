package api

import (
"net/url"
"log"
"github.com/gorilla/websocket"
//"fmt"
"fmt"
	"encoding/json"
	//"net"
	//"bytes"
)


const gdaxHost = "ws-feed.gdax.com"
//const gdaxPath = "/api/2/ws"

type GdaxApi struct {
	connection *websocket.Conn
}

type GdaxSubscription struct {
	Type     string `json:"type"`
	Channels []Channel `json:"channels"`
}

type Channel struct {
Name       string   `json:"name"`
ProductIds []string `json:"product_ids"`
}



func (b *GdaxApi)  StartListen(apiCurrenciesConfiguration ApiCurrenciesConfiguration, callback func(message []byte, error error)) {
	url := url.URL{Scheme: "wss", Host: gdaxHost, Path: ""}
	log.Printf("connecting to %s", url.String())

	connection, _, error := websocket.DefaultDialer.Dial(url.String(), nil)

	if error != nil {
		fmt.Println("Gdax ws error: ",error)
		callback(nil, error)
	} else if connection != nil {
		fmt.Println("Gdax ws connected")

		productsIds :=  b.composeSymbolsForSubscirbe(apiCurrenciesConfiguration)

		for _, productId := range  productsIds {

			channel := Channel{}
			channel.Name = "ticker"
			channel.ProductIds = []string{productId}

			subscribtion := GdaxSubscription{}
			subscribtion.Type = "subscribe"
			subscribtion.Channels = []Channel{channel}

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

func (b *GdaxApi)  StopListen() {
	//fmt.Println("before close")
	//b.connection.Close()
	//fmt.Println("closed")
}

func (b *GdaxApi)  composeSymbolsForSubscirbe(apiCurrenciesConfiguration ApiCurrenciesConfiguration) []string {
	var smybolsForSubscirbe = []string{}
	for _, targetCurrency := range apiCurrenciesConfiguration.TargetCurrencies {
		for _, referenceCurrency := range apiCurrenciesConfiguration.ReferenceCurrencies {

			if targetCurrency == referenceCurrency {
				continue
			}

			symbol := targetCurrency + "-" + referenceCurrency
			smybolsForSubscirbe = append(smybolsForSubscirbe, symbol)
		}
	}
	return smybolsForSubscirbe

}
