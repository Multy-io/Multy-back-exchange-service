package api

import (
	"net/url"
	"log"
	"github.com/gorilla/websocket"
	//"fmt"
	"fmt"
)


const okexHost = "real.okex.com:10440"
const okexPath = "/websocket/okexapi"

type OkexApi struct {
	connection *websocket.Conn
}




func (b *OkexApi)  StartListen(callback func(message []byte, error error)) {
	url := url.URL{Scheme: "wss", Host: okexHost, Path: ""}
	log.Printf("connecting to %s", url.String())

	connection, _, error := websocket.DefaultDialer.Dial(url.String(), nil)

	if error != nil {
		fmt.Println("Okex ws error: ",error)
		callback(nil, error)
	} else if connection != nil {
		fmt.Println("Okex ws connected")

		//TODO: get symbols from exhange
		subscribtion0 := `{"event":"addChannel","channel":"ok_sub_futureusd_btc_ticker_this_week"}`
		subscribtion1 := `{"event":"addChannel","channel":"ok_sub_futureusd_ltc_ticker_this_week"}`
		subscribtion2 := `{"event":"addChannel","channel":"ok_sub_futureusd_eth_ticker_this_week"}`
		subscribtion3 := `{"event":"addChannel","channel":"ok_sub_futureusd_etc_ticker_this_week"}`
		subscribtion4 := `{"event":"addChannel","channel":"ok_sub_futureusd_bch_ticker_this_week"}`
		connection.WriteMessage(websocket.TextMessage, []byte(subscribtion0))
		connection.WriteMessage(websocket.TextMessage, []byte(subscribtion1))
		connection.WriteMessage(websocket.TextMessage, []byte(subscribtion2))
		connection.WriteMessage(websocket.TextMessage, []byte(subscribtion3))
		connection.WriteMessage(websocket.TextMessage, []byte(subscribtion4))



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

func (b *OkexApi)  StopListen() {
	//fmt.Println("before close")
	//b.connection.Close()
	//fmt.Println("closed")
}