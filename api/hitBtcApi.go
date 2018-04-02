package api

import (
	"net/url"
	"log"
	"github.com/gorilla/websocket"
	//"fmt"
	"fmt"
)


const hitBtcHost = "api.hitbtc.com"
const hitBtcPath = "/api/2/ws"

type HitBtcApi struct {
	connection *websocket.Conn
}




func (b HitBtcApi)  StartListen(callback func(message []byte, error error)) {
	url := url.URL{Scheme: "wss", Host: hitBtcHost, Path: hitBtcPath}
	log.Printf("connecting to %s", url.String())

	connection, _, error := websocket.DefaultDialer.Dial(url.String(), nil)

	if error != nil {
		fmt.Println("HitBtc ws error: ",error)
		callback(nil, error)
	} else if connection != nil {
		fmt.Println("HitBtc ws connected")

		//TODO: get symbols from exhange
		subscribtionBTCUSD := `{"method":"subscribeTicker","params":{"symbol": "BTCUSD"},"id": 10000}`
		subscribtionETHBTC := `{"method":"subscribeTicker","params":{"symbol":"ETHBTC"},"id": 10000}`
		subscribtionETHUSD := `{"method":"subscribeTicker","params":{"symbol":"ETHUSD"},"id": 10000}`

		subscribtionBCHUSD := `{"method":"subscribeTicker","params":{"symbol":"BCHUSD"},"id": 10000}`
		subscribtionLTCUSD := `{"method":"subscribeTicker","params":{"symbol":"LTCUSD"},"id": 10000}`
		subscribtionXMRUSD := `{"method":"subscribeTicker","params":{"symbol":"XMRUSD"},"id": 10000}`
		subscribtionDASHUSD := `{"method":"subscribeTicker","params":{"symbol":"DASHUSD"},"id": 10000}`
		subscribtionEOSUSD := `{"method":"subscribeTicker","params":{"symbol":"EOSUSD"},"id": 10000}`
		subscribtionXRPUSD := `{"method":"subscribeTicker","params":{"symbol":"XRPUSD"},"id": 10000}`
		subscribtionZECUSD := `{"method":"subscribeTicker","params":{"symbol":"ZECUSD"},"id": 10000}`


		connection.WriteMessage(websocket.TextMessage, []byte(subscribtionBTCUSD))
		connection.WriteMessage(websocket.TextMessage, []byte(subscribtionETHBTC))
		connection.WriteMessage(websocket.TextMessage, []byte(subscribtionETHUSD))

		connection.WriteMessage(websocket.TextMessage, []byte(subscribtionBCHUSD))
		connection.WriteMessage(websocket.TextMessage, []byte(subscribtionLTCUSD))
		connection.WriteMessage(websocket.TextMessage, []byte(subscribtionXMRUSD))
		connection.WriteMessage(websocket.TextMessage, []byte(subscribtionDASHUSD))
		connection.WriteMessage(websocket.TextMessage, []byte(subscribtionEOSUSD))
		connection.WriteMessage(websocket.TextMessage, []byte(subscribtionXRPUSD))
		connection.WriteMessage(websocket.TextMessage, []byte(subscribtionZECUSD))


		for {
			func() {
				_, message, error := connection.ReadMessage()
				callback(message, error)
			}()
		}
	} else {
		fmt.Println("connection is nil")
		callback(nil, nil)
	}


}

func (b HitBtcApi)  StopListen() {
	//fmt.Println("before close")
	//b.connection.Close()
	//fmt.Println("closed")
}