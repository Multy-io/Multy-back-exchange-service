package api

import (
	"net/url"

	"fmt"

	_ "github.com/KristinaEtc/slflog"
	"github.com/gorilla/websocket"
)

const bitfinexHost = "api.bitfinex.com"
const bitfinexPath = "/ws/2"

type biftfinexSubscription struct {
	Command string `json:"event"`
	Channel string `json:"channel"`
	Symbol  string `json:"symbol"`
}

type BitfinexApi struct {
	connection           *websocket.Conn
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

func (b *BitfinexApi) connectWs(apiCurrenciesConfiguration ApiCurrenciesConfiguration) *websocket.Conn {
	url := url.URL{Scheme: "wss", Host: bitfinexHost, Path: bitfinexPath}
	//log.Printf("connecting to %s", url.String())

	connection, _, err := websocket.DefaultDialer.Dial(url.String(), nil)

	if err != nil || connection == nil {
		log.Errorf("connectWs:Bitfinex ws connection error: ", err)
		return nil
	} else {
		log.Debugf("connectWs:Bitfinex ws connected")
		b.symbolesForSubscirbe = b.composeSymbolsForSubscirbe(apiCurrenciesConfiguration)
		for _, symbol := range b.symbolesForSubscirbe {
			subscribtion := `{"event":"subscribe","channel":"ticker","symbol": "` + symbol + `"}`
			//fmt.Println(subscribtion)
			connection.WriteMessage(websocket.TextMessage, []byte(subscribtion))
		}
		return connection
	}
}

func (b *BitfinexApi) StartListen(apiCurrenciesConfiguration ApiCurrenciesConfiguration, ch chan Reposponse) {
	fmt.Println("StartListen:Start listen Bitfinex")
	for {
		if b.connection == nil {
			b.connection = b.connectWs(apiCurrenciesConfiguration)
		} else if b.connection != nil {

			func() {
				_, message, err := b.connection.ReadMessage()
				if err != nil {
					log.Errorf("StartListen:Bitfinex read message error: %v", err.Error())
					b.connection.Close()
					b.connection = nil
				} else {
					//fmt.Printf("%f", message)
					ch <- Reposponse{Message: &message, Err: &err}
				}
			}()
		}
	}

}

func (b *BitfinexApi) composeSymbolsForSubscirbe(apiCurrenciesConfiguration ApiCurrenciesConfiguration) []string {
	var smybolsForSubscirbe = []string{}
	for _, targetCurrency := range apiCurrenciesConfiguration.TargetCurrencies {
		for _, referenceCurrency := range apiCurrenciesConfiguration.ReferenceCurrencies {

			if targetCurrency == referenceCurrency {
				continue
			}

			if targetCurrency == "DASH" {
				targetCurrency = "DSH"
			}

			if referenceCurrency == "USDT" {
				referenceCurrency = "USD"
			}
			symbol := "t" + targetCurrency + referenceCurrency
			smybolsForSubscirbe = append(smybolsForSubscirbe, symbol)
		}
	}
	return smybolsForSubscirbe

}

func (b *BitfinexApi) StopListen() {
	if b.connection != nil {
		b.connection.Close()
		b.connection = nil
	}
	log.Debugf("Bitfinex ws closed")
}
