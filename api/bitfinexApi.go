package api

import (
	"net/url"
	"fmt"
	//"encoding/json"
	"github.com/gorilla/websocket"
	//"log"
)

const bitfinexHost = "api.bitfinex.com"
const bitfinexPath = "/ws/2"

type biftfinexSubscription struct {
	Command string `json:"event"`
	Channel string `json:"channel"`
	Symbol string `json:"symbol"`
}

type BitfinexApi struct {
	connection *websocket.Conn
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


func (b *BitfinexApi)  connectWs(apiCurrenciesConfiguration ApiCurrenciesConfiguration) *websocket.Conn {
	url := url.URL{Scheme: "wss", Host: bitfinexHost, Path: bitfinexPath}
	//log.Printf("connecting to %s", url.String())

	connection, _, error := websocket.DefaultDialer.Dial(url.String(), nil)

	if error != nil || connection == nil {
		//fmt.Println("Bitfinex ws connection error: ",error)
		return nil
	} else  {
		fmt.Println("Bitfinex ws connected")

		b.symbolesForSubscirbe = b.composeSymbolsForSubscirbe(apiCurrenciesConfiguration)
		for _, symbol := range b.symbolesForSubscirbe {
			subscribtion := `{"event":"subscribe","channel":"ticker","symbol": "` + symbol + `"}`
			connection.WriteMessage(websocket.TextMessage, []byte(subscribtion))
		}

		return connection
	}
}


func (b *BitfinexApi)  StartListen(apiCurrenciesConfiguration ApiCurrenciesConfiguration, callback  func(message []byte, error error)) {

	for {
		if b.connection == nil {
			b.connection = b.connectWs(apiCurrenciesConfiguration)
		} else if b.connection != nil {

			func() {
				_, message, error := b.connection.ReadMessage()
				if error != nil {
					fmt.Println("Bitfinex read message error:", error)
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
	if b.connection != nil {
		b.connection.Close()
		b.connection = nil
	}
	fmt.Println("Bitfinex ws closed")
}

