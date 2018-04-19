package api

import (
	"encoding/json"
	"net/url"

	_ "github.com/KristinaEtc/slflog"
	"github.com/gorilla/websocket"
)

const gdaxHost = "ws-feed.gdax.com"

//const gdaxPath = "/api/2/ws"

type GdaxApi struct {
	connection *websocket.Conn
}

type GdaxSubscription struct {
	Type     string    `json:"type"`
	Channels []Channel `json:"channels"`
}

type Channel struct {
	Name       string   `json:"name"`
	ProductIds []string `json:"product_ids"`
}

func (b *GdaxApi) connectWs(apiCurrenciesConfiguration ApiCurrenciesConfiguration) *websocket.Conn {
	url := url.URL{Scheme: "wss", Host: gdaxHost, Path: ""}
	log.Infof("connecting to %s", url.String())

	connection, _, err := websocket.DefaultDialer.Dial(url.String(), nil)

	if err != nil || connection == nil {
		log.Errorf("Gdax ws connection error:%v", err.Error())
		return nil
	} else {
		log.Debugf("Gdax ws connected")

		productsIds := b.composeSymbolsForSubscirbe(apiCurrenciesConfiguration)

		for _, productId := range productsIds {

			channel := Channel{}
			channel.Name = "ticker"
			channel.ProductIds = []string{productId}

			subscribtion := GdaxSubscription{}
			subscribtion.Type = "subscribe"
			subscribtion.Channels = []Channel{channel}

			msg, _ := json.Marshal(subscribtion)
			connection.WriteMessage(websocket.TextMessage, msg)
		}

		return connection
	}
}

func (b *GdaxApi) StartListen(apiCurrenciesConfiguration ApiCurrenciesConfiguration, callback func(message []byte, err error)) {
	for {
		if b.connection == nil {
			b.connection = b.connectWs(apiCurrenciesConfiguration)
		} else if b.connection != nil {

			func() {
				_, message, err := b.connection.ReadMessage()
				if err != nil {
					log.Errorf("Gdax read message error:", err.Error())
					b.connection.Close()
					b.connection = nil
				} else {
					callback(message, err)
				}
			}()
		}
	}
}

func (b *GdaxApi) StopListen() {
	if b.connection != nil {
		b.connection.Close()
		b.connection = nil
	}
	log.Debugf("Gdax ws closed")
}

func (b *GdaxApi) composeSymbolsForSubscirbe(apiCurrenciesConfiguration ApiCurrenciesConfiguration) []string {
	var smybolsForSubscirbe = []string{}
	for _, targetCurrency := range apiCurrenciesConfiguration.TargetCurrencies {
		for _, referenceCurrency := range apiCurrenciesConfiguration.ReferenceCurrencies {

			if targetCurrency == referenceCurrency {
				continue
			}

			if referenceCurrency == "USDT" {
				referenceCurrency = "USD"
			}

			symbol := targetCurrency + "-" + referenceCurrency
			smybolsForSubscirbe = append(smybolsForSubscirbe, symbol)
		}
	}
	return smybolsForSubscirbe

}
