package api

import (
"net/url"
"log"
"github.com/gorilla/websocket"
//"fmt"
"fmt"
)


const gdaxHost = "ws-feed.gdax.com"
//const gdaxPath = "/api/2/ws"

type GdaxApi struct {
	connection *websocket.Conn
}




func (b GdaxApi)  StartListen(callback func(message []byte, error error)) {
	url := url.URL{Scheme: "wss", Host: gdaxHost, Path: ""}
	log.Printf("connecting to %s", url.String())

	connection, _, error := websocket.DefaultDialer.Dial(url.String(), nil)

	if error != nil {
		fmt.Println("Gdax ws error: ",error)
		callback(nil, error)
	} else if connection != nil {
		fmt.Println("Gdax ws connected")

		//TODO: get symbols from exhange
		subscribtion := `{"type":"subscribe","channels":[{"name": "ticker", "product_ids":["BTC-USD", "ETH-USD", "LTC-USD", "BCH-USD"]}]}`
		connection.WriteMessage(websocket.TextMessage, []byte(subscribtion))



		for {
			func() {
				_, message, _ := connection.ReadMessage()
				//fmt.Printf("%s \n", message)
				callback(message, error)
			}()
		}
	} else {
		fmt.Println("connection is nil")
		callback(nil, nil)
	}


}

func (b GdaxApi)  StopListen() {
	//fmt.Println("before close")
	//b.connection.Close()
	//fmt.Println("closed")
}