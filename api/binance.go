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

func (b *BinanceApi)  connectWs() *websocket.Conn {
	url := url.URL{Scheme: "wss", Host: binanceHost, Path: tickerPath}
	log.Printf("connecting to %s", url.String())

	connection, _, error := websocket.DefaultDialer.Dial(url.String(), nil)
	if error != nil || connection == nil  {
		fmt.Println("Binance ws connection error: ",error)
		return nil
	} else {
		fmt.Println("Binance ws connected")
		return connection
	}
}

func (b *BinanceApi)  StartListen(callback func(message []byte, error error)) {

	for range time.Tick(1 * time.Second) {
		if b.connection == nil {
			b.connection = b.connectWs()
		} else if b.connection != nil {
			func() {
				_, message, error := b.connection.ReadMessage()
				if error != nil {
					fmt.Println("Binance read message error:", error)
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

func (b *BinanceApi)  StopListen() {
	if b.connection != nil {
		b.connection.Close()
		b.connection = nil
	}
	fmt.Println("Binance ws closed")
}