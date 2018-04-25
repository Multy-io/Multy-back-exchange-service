package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
	"time"

	_ "github.com/KristinaEtc/slflog"
	"github.com/gorilla/websocket"
	"github.com/shopspring/decimal"
)

const host = "api2.poloniex.com"
const path = "/realm1"

const (
	origin        = "https://api2.poloniex.com/"
	pushAPIUrl    = "wss://api2.poloniex.com/realm1"
	publicAPIUrl  = "https://poloniex.com/public?command="
	tradingAPIUrl = "https://poloniex.com/tradingApi"
)

type subscription struct {
	Command string `json:"command"`
	Channel string `json:"channel"`
}

type PoloniexApi struct {
	connection *websocket.Conn
	logger     Logger
	LogBus     chan<- string
	httpClient *http.Client
}

func NewPoloniexApi() *PoloniexApi {
	var api = PoloniexApi{}
	api.httpClient = &http.Client{Timeout: time.Second * 10}
	return &api
}

func (b *PoloniexApi) connectWs() *websocket.Conn {
	url := url.URL{Scheme: "wss", Host: host, Path: path}
	log.Infof("connectWs:onnecting to %s", url.String())

	connection, _, err := websocket.DefaultDialer.Dial(url.String(), nil)

	if err != nil || connection == nil {
		log.Errorf("connectWs:Poloniex ws connection error: %v", err.Error())
		return nil
	} else {
		log.Debugf("connectWs:Poloniex ws connected")
		subs := subscription{Command: "subscribe", Channel: "1002"}
		msg, _ := json.Marshal(subs)
		connection.WriteMessage(websocket.BinaryMessage, msg)
		return connection
	}

}

func (b *PoloniexApi) StartListen(ch chan Reposponse) {

	for {
		if b.connection == nil {
			b.connection = b.connectWs()
		} else if b.connection != nil {
			func() {
				_, message, err := b.connection.ReadMessage()
				if err != nil {
					log.Errorf("StartListen:Poloniex read message error: %v", err.Error())
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

type Ticker struct {
	ID            int             `json:"id, int"`
	Last          decimal.Decimal `json:"last, string"`
	LowestAsk     decimal.Decimal `json:"lowestAsk, string"`
	HighestBid    decimal.Decimal `json:"highestBid, string"`
	PercentChange decimal.Decimal `json:"percentChange, string"`
	BaseVolume    decimal.Decimal `json:"baseVolume, string"`
	QuoteVolume   decimal.Decimal `json:"quoteVolume, string"`
	IsFrozen      int             `json:"isFrozen ,string"`
	High24hr      decimal.Decimal `json:"high24hr, string"`
	Low24hr       decimal.Decimal `json:"low24hr, string"`
}

func (p *PoloniexApi) PubReturnTickers() (tickers map[string]Ticker, err error) {

	respch := make(chan []byte)
	errch := make(chan error)

	go p.publicRequest("returnTicker", respch, errch)

	response := <-respch
	err = <-errch

	if err != nil {
		return
	}

	err = json.Unmarshal(response, &tickers)
	return
}

var (
	//Poloniex says we are allowed 6 req/s
	//but this is not true if you don't want to see
	//'nonce must be greater than' error 3 req/s is the best option.
	throttle = time.Tick(time.Second / 3)
)

func (p *PoloniexApi) publicRequest(action string, respch chan<- []byte, errch chan<- error) {

	<-throttle

	defer close(respch)
	defer close(errch)

	rawurl := publicAPIUrl + action

	req, err := http.NewRequest("GET", rawurl, nil)
	if err != nil {
		respch <- nil
		errch <- Error(RequestError)
		return
	}

	req.Header.Add("Accept", "application/json")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		respch <- nil
		errch <- Error(ConnectError)
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		respch <- body
		errch <- err
		return
	}

	respch <- body
	errch <- nil
}

func (b *PoloniexApi) StopListen() {
	//fmt.Println("before close")
	//b.connection.Close()
	//fmt.Println("closed")
}

var (
	ConnectError    = "[ERROR] Connection could not be established!"
	RequestError    = "[ERROR] NewRequest Error!"
	SetApiError     = "[ERROR] Set the API KEY and API SECRET!"
	PeriodError     = "[ERROR] Invalid Period!"
	TimePeriodError = "[ERROR] Time Period incompatibility!"
	TimeError       = "[ERROR] Invalid Time!"
	StartTimeError  = "[ERROR] Start Time Format Error!"
	EndTimeError    = "[ERROR] End Time Format Error!"
	LimitError      = "[ERROR] Limit Format Error!"
	ChannelError    = "[ERROR] Unknown Channel Name: %s"
	SubscribeError  = "[ERROR] Already Subscribed!"
	WSTickerError   = "[ERROR] WSTicker Parsing %s"
	OrderBookError  = "[ERROR] MarketUpdate OrderBook Parsing %s"
	NewTradeError   = "[ERROR] MarketUpdate NewTrade Parsing %s"
	ServerError     = "[SERVER ERROR] Response: %s"
)

func Error(msg string, args ...interface{}) error {
	if len(args) > 0 {
		return errors.New(fmt.Sprintf(msg, args))
	} else {
		return errors.New(msg)
	}
}

type Logger struct {
	isOpen bool
	Lock   *sync.Mutex
}
