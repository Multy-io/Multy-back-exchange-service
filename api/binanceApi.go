package api

import (
	"net/url"
	"log"
	"github.com/gorilla/websocket"
	//"fmt"
	"time"
	"fmt"
)

const binanceHost = "stream.binance.com:9443"
const tickerPath = "/ws/!ticker@arr"

type BinanceApi struct {
	connection *websocket.Conn
}

func (b BinanceApi)  StartListen(callback func(message []byte, error error)) {
	url := url.URL{Scheme: "wss", Host: binanceHost, Path: tickerPath}
	log.Printf("connecting to %s", url.String())

	connection, _, error := websocket.DefaultDialer.Dial(url.String(), nil)
	b.connection = connection
	if error != nil {
		callback(nil, error)
	} else if connection != nil {
		fmt.Println("Binance ws connected")
		func() {
			for range time.Tick(1 * time.Second) {
				_, message, error := b.connection.ReadMessage()
				callback(message, error)
			}
		}()
	} else {
			fmt.Println("connection is nil")
			callback(nil, nil)
	}



}

func (b BinanceApi)  StopListen() {
	//fmt.Println("before close")
	//b.connection.Close()
	//fmt.Println("closed")
}