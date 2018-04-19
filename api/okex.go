package api

import (
	"encoding/json"
	"net/url"

	_ "github.com/KristinaEtc/slflog"
	"github.com/gorilla/websocket"
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

func (b *OkexApi) connectWs(apiCurrenciesConfiguration ApiCurrenciesConfiguration) *websocket.Conn {
	url := url.URL{Scheme: "wss", Host: okexHost, Path: ""}
	log.Infof("connecting to %s", url.String())

	connection, _, err := websocket.DefaultDialer.Dial(url.String(), nil)

	if err != nil || connection == nil {
		log.Errorf("Okex ws connection error: %v", err.Error())
		return nil
	} else {
		log.Debugf("Okex ws connected")

		productsIds := b.composeSymbolsForSubscirbe(apiCurrenciesConfiguration)

		for _, productId := range productsIds {
			subscribtion := OkexSubscription{}
			subscribtion.Event = "addChannel"
			subscribtion.Channel = productId

			msg, _ := json.Marshal(subscribtion)
			connection.WriteMessage(websocket.TextMessage, msg)
		}

		return connection
	}
}

func (b *OkexApi) StartListen(apiCurrenciesConfiguration ApiCurrenciesConfiguration, callback func(message []byte, err error)) {
	for {
		if b.connection == nil {
			b.connection = b.connectWs(apiCurrenciesConfiguration)
		} else if b.connection != nil {

			func() {
				_, message, err := b.connection.ReadMessage()
				if err != nil {
					log.Errorf("okex read message error: %v", err.Error())
					b.connection.Close()
					b.connection = nil
				} else {
					callback(message, err)
				}
			}()
		}
	}

}

func (b *OkexApi) StopListen() {
	if b.connection != nil {
		b.connection.Close()
		b.connection = nil
	}
	log.Debugf("Okex ws closed")
}

func (b *OkexApi) composeSymbolsForSubscirbe(apiCurrenciesConfiguration ApiCurrenciesConfiguration) []string {
	var smybolsForSubscirbe = []string{}
	for _, targetCurrency := range apiCurrenciesConfiguration.TargetCurrencies {

		symbol := "ok_sub_futureusd_" + targetCurrency + "_ticker_this_week"
		smybolsForSubscirbe = append(smybolsForSubscirbe, symbol)

	}
	return smybolsForSubscirbe

}
