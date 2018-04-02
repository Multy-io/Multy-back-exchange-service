package api

import (
	"net/url"
	"fmt"
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
)

const host = "api2.poloniex.com"
const path = "/realm1"

type subscription struct {
	Command string `json:"command"`
	Channel string `json:"channel"`
}

type PoloniexApi struct {
	//connection *websocket.Conn
}



func (b PoloniexApi)  StartListen(callback func(message []byte, error error)) {
	url := url.URL{Scheme: "wss", Host: host, Path: path}
	log.Printf("connecting to %s", url.String())

	connection, _, error := websocket.DefaultDialer.Dial(url.String(), nil)

	if error != nil {
		fmt.Println("Poloniex ws error: ",error)
		callback(nil, error)
	} else if connection != nil {
		fmt.Println("Poloniex ws connected")

		subs := subscription{Command: "subscribe", Channel: "1002"}
		msg, _ := json.Marshal(subs)
		connection.WriteMessage(websocket.BinaryMessage, msg)

		for {
			func() {
				_, messageJSON, _ := connection.ReadMessage()
					//fmt.Println(messageJSON)
					callback(messageJSON, error)

			}()
		}
	} else {
		fmt.Println("connection is nil")
		callback(nil, nil)
	}


}

func (b PoloniexApi)  StopListen() {
	//fmt.Println("before close")
	//b.connection.Close()
	//fmt.Println("closed")
}

