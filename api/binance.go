package api

import (
	"fmt"
	"net/url"
	"time"

	"github.com/KristinaEtc/slf"
	_ "github.com/KristinaEtc/slflog"
	"github.com/gorilla/websocket"
)

var log = slf.WithContext("api")

const binanceHost = "stream.binance.com:9443"
const tickerPath = "/ws/!ticker@arr"

type BinanceApi struct {
	connection *websocket.Conn
}

type Reposponse struct {
	Message *[]byte
	Err     *error
}

func (b *BinanceApi) connectWs() *websocket.Conn {
	url := url.URL{Scheme: "wss", Host: binanceHost, Path: tickerPath}
	log.Infof("connectWs:connecting to %s", url.String())

	connection, _, err := websocket.DefaultDialer.Dial(url.String(), nil)
	if err != nil || connection == nil {
		fmt.Errorf("connectWs:Binance ws connection error: %v", err.Error())
		return nil
	} else {
		log.Debugf("connectWs:Binance ws connected")
		return connection
	}
}

func (b *BinanceApi) StartListen(ch chan Reposponse) {

	for range time.Tick(1 * time.Second) {
		if b.connection == nil {
			b.connection = b.connectWs()
		} else if b.connection != nil {
			func() {
				_, message, err := b.connection.ReadMessage()
				if err != nil {
					log.Errorf("StartListen:Binance read message error: %v", err.Error())
					b.connection.Close()
					b.connection = nil
				} else {
					//fmt.Printf("%s \n", message)
					ch <- Reposponse{Message: &message, Err: &err}
				}
			}()
		}
	}
}

func (b *BinanceApi) StopListen() {
	if b.connection != nil {
		b.connection.Close()
		b.connection = nil
	}
	log.Debugf("StopListen:Binance ws closed")
}
