package exchanger

import (
	"github.com/Appscrunch/Multy-back-exchange-service/core"
	"github.com/Appscrunch/Multy-back-exchange-service/exchange-rates"
)

type Exchanger struct {
	Manager   *core.Manager
	Exchanger *exchangeRates.ExchangeManager
}

func (e *Exchanger) InitExchanger(conf core.ManagerConfiguration) {
	//var exchangeManger *exchangeRates.ExchangeManager
	//
	//var manager = core.NewManager()
	//go manager.StartListen(conf)
	//
	//exchangeManger = exchangeRates.NewExchangeManager()
	//go exchangeManger.StartGetingData(conf)
	//
	//
	//waitGroup.Add(len(configuration.Exchanges) + 5)
	//
	//go manager.StartListen(configuration)
	//
	//exchangeManger = exchangeRates.NewExchangeManager()
	//go exchangeManger.StartGetingData(configuration)
	//
	//waitGroup.Wait()
}
