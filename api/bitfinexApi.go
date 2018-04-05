package api

import (
	"net/url"
	"fmt"
	//"encoding/json"
	"github.com/gorilla/websocket"
	"log"
)

const bitfinexHost = "api.bitfinex.com"
const bitfinexPath = "/ws/2"

type biftfinexSubscription struct {
	Command string `json:"event"`
	Channel string `json:"channel"`
	Symbol string `json:"symbol"`
}

type BitfinexApi struct {
	//connection *websocket.Conn
	symbolesForSubscirbe []string
}

type ApiCurrenciesConfiguration struct {
	TargetCurrencies    []string
	ReferenceCurrencies []string
}

func NewBitfinexApi() *BitfinexApi {
	var api = BitfinexApi{}
	//api.symbolesForSubscirbe = []string{"tBTCUSD", "tETHUSD","tBTSUSD", "tSTEEMUSD", "tWAVESUSD", "tLTCUSD", "tBCHUSD", "tETCUSD", "tDASHUSD", "tEOSUSD",  "tETHBTC","tBTSBTC", "tSTEEMBTC", "tWAVESBTC", "tLTCBTC", "tBCHBTC", "tETCBTC", "tDASHBTC", "tEOSBTC"}
	return &api
}



func (b *BitfinexApi)  StartListen(apiCurrenciesConfiguration ApiCurrenciesConfiguration, callback  func(message []byte, error error)) {

	url := url.URL{Scheme: "wss", Host: bitfinexHost, Path: bitfinexPath}
	log.Printf("connecting to %s", url.String())

	connection, _, error := websocket.DefaultDialer.Dial(url.String(), nil)

	if error != nil {
		fmt.Println("Bitfinex ws error: ",error)
		callback(nil, error)
	} else if connection != nil {
		fmt.Println("Bitfinex ws connected")

		b.symbolesForSubscirbe = b.composeSymbolsForSubscirbe(apiCurrenciesConfiguration)
		for _, symbol := range b.symbolesForSubscirbe {
			subscribtion := `{"event":"subscribe","channel":"ticker","symbol": "` + symbol + `"}`
			connection.WriteMessage(websocket.TextMessage, []byte(subscribtion))
		}

		for {
			func() {
				_, messageJSON, _ := connection.ReadMessage()
				//fmt.Println("read", messageJSON)
				callback(messageJSON, error)

			}()
		}
	} else {
		fmt.Println("connection is nil")
		callback(nil, nil)
	}


}

func (b *BitfinexApi)  composeSymbolsForSubscirbe(apiCurrenciesConfiguration ApiCurrenciesConfiguration) []string {
	var smybolsForSubscirbe = []string{}
	for _, targetCurrency := range apiCurrenciesConfiguration.TargetCurrencies {
		for _, referenceCurrency := range apiCurrenciesConfiguration.ReferenceCurrencies {
			symbol := "t" + targetCurrency + referenceCurrency
			smybolsForSubscirbe = append(smybolsForSubscirbe, symbol)
		}
	}
	return smybolsForSubscirbe

}


func (b *BitfinexApi)  StopListen() {
	//fmt.Println("before close")
	//b.connection.Close()
	//fmt.Println("closed")
}

