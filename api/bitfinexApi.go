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
}



func (b BitfinexApi)  StartListen(callback func(message []byte, error error)) {
	url := url.URL{Scheme: "wss", Host: bitfinexHost, Path: bitfinexPath}
	log.Printf("connecting to %s", url.String())

	connection, _, error := websocket.DefaultDialer.Dial(url.String(), nil)

	if error != nil {
		fmt.Println("Bitfinex ws error: ",error)
		callback(nil, error)
	} else if connection != nil {
		fmt.Println("Bitfinex ws connected")


		subscribtion0 := `{"event":"subscribe","channel":"ticker","symbol": "tBTCUSD"}`
		subscribtion1 := `{"event":"subscribe","channel":"ticker","symbol": "tETHUSD"}`

		subscribtion2 := `{"event":"subscribe","channel":"ticker","symbol": "tLTCUSD"}`
		subscribtion3 := `{"event":"subscribe","channel":"ticker","symbol": "tBCHUSD"}`

		subscribtion4 := `{"event":"subscribe","channel":"ticker","symbol": "tETCUSD"}`
		subscribtion5 := `{"event":"subscribe","channel":"ticker","symbol": "tEOSUSD"}`


		connection.WriteMessage(websocket.TextMessage, []byte(subscribtion0))
		connection.WriteMessage(websocket.TextMessage, []byte(subscribtion1))
		connection.WriteMessage(websocket.TextMessage, []byte(subscribtion2))
		connection.WriteMessage(websocket.TextMessage, []byte(subscribtion3))
		connection.WriteMessage(websocket.TextMessage, []byte(subscribtion4))
		connection.WriteMessage(websocket.TextMessage, []byte(subscribtion5))

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

func (b BitfinexApi)  StopListen() {
	//fmt.Println("before close")
	//b.connection.Close()
	//fmt.Println("closed")
}

