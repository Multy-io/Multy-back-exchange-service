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


func (b *OkexApi)  connectWs(apiCurrenciesConfiguration ApiCurrenciesConfiguration) *websocket.Conn {
	url := url.URL{Scheme: "wss", Host: okexHost, Path: ""}
	log.Printf("connecting to %s", url.String())

	connection, _, error := websocket.DefaultDialer.Dial(url.String(), nil)

	if error != nil || connection == nil {
		fmt.Println("Okex ws connection error: ",error)
		return nil
	} else  {
		fmt.Println("Okex ws connected")

		productsIds :=  b.composeSymbolsForSubscirbe(apiCurrenciesConfiguration)

		for _, productId := range  productsIds {
			subscribtion := OkexSubscription{}
			subscribtion.Event = "addChannel"
			subscribtion.Channel = productId

			msg, _ := json.Marshal(subscribtion)
			connection.WriteMessage(websocket.TextMessage, msg)
		}

		return connection
	}
}


func (b *OkexApi)  StartListen(apiCurrenciesConfiguration ApiCurrenciesConfiguration, callback func(message []byte, error error)) {

	for {
		if b.connection == nil {
			b.connection = b.connectWs(apiCurrenciesConfiguration)
		} else if b.connection != nil {

			func() {
				_, message, error := b.connection.ReadMessage()
				if error != nil {
					fmt.Println("okex read message error:", error)
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

func (b *OkexApi)  StopListen() {
	//fmt.Println("before close")
	if b.connection != nil {
		b.connection.Close()
		b.connection = nil
	}
	fmt.Println("Okex ws closed")
}

func (b *OkexApi)  composeSymbolsForSubscirbe(apiCurrenciesConfiguration ApiCurrenciesConfiguration) []string {
	var smybolsForSubscirbe = []string{}
	for _, targetCurrency := range apiCurrenciesConfiguration.TargetCurrencies {

		symbol := "ok_sub_futureusd_" +  targetCurrency + "_ticker_this_week"
		smybolsForSubscirbe = append(smybolsForSubscirbe, symbol)

	}
	return smybolsForSubscirbe

}
