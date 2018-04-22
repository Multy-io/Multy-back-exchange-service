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
	var manager = core.NewManager()
	go manager.StartListen(conf)
	e.Manager = manager

	exManager := exchangeRates.NewExchangeManager(conf)
	go exManager.StartGetingData()
	e.Exchanger = exManager
}
