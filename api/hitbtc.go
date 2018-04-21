package api

import (
	"encoding/json"
	"net/url"

	_ "github.com/KristinaEtc/slflog"
	"github.com/gorilla/websocket"
)

const hitBtcHost = "api.hitbtc.com"
const hitBtcPath = "/api/2/ws"

type HitBtcApi struct {
	connection *websocket.Conn
}

type HitBtcSubscription struct {
	Method string                   `json:"method"`
	Params HitBtcSubscriptionParams `json:"params"`
	ID     int                      `json:"id"`
}

type HitBtcSubscriptionParams struct {
	Symbol string `json:"symbol"`
}

func (b *HitBtcApi) connectWs(apiCurrenciesConfiguration ApiCurrenciesConfiguration) *websocket.Conn {
	url := url.URL{Scheme: "wss", Host: hitBtcHost, Path: hitBtcPath}
	log.Infof("connecting to %s", url.String())

	connection, _, err := websocket.DefaultDialer.Dial(url.String(), nil)

	if err != nil || connection == nil {
		log.Errorf("HitBtc ws connection error: %v", err.Error())
		return nil
	} else {
		log.Debugf("HitBtc ws connected")

		productsIds := b.composeSymbolsForSubscirbe(apiCurrenciesConfiguration)

		for _, productId := range productsIds {

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

func (b *HitBtcApi) StartListen(apiCurrenciesConfiguration ApiCurrenciesConfiguration, ch chan Reposponse) {


	for {
		if b.connection == nil {
			b.connection = b.connectWs(apiCurrenciesConfiguration)
		} else if b.connection != nil {

			func() {
				_, message, err := b.connection.ReadMessage()
				if err != nil {
					log.Errorf("HitBtc read message error:%v", err.Error())
					b.connection.Close()
					b.connection = nil
				} else {
					ch <- Reposponse{Message:&message, Err:&err}
				}
			}()
		}
	}

}

func (b *HitBtcApi) StopListen() {
	if b.connection != nil {
		b.connection.Close()
		b.connection = nil
	}
	log.Debugf("HitBtc ws closed")
}

func (b *HitBtcApi) composeSymbolsForSubscirbe(apiCurrenciesConfiguration ApiCurrenciesConfiguration) []string {
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
