# Multy-back-exchange-service
Exchange service for Multy back

1 install postgres 10;
2 run shell srcript createDB.sh that will create DB exampple: "sh createDB.sh "postgres" "postgres" "test""
3 run sql script from sql/setupDB file and check that all tables,views,indexes and 2 functions are created;
4 init Exchanger with func InitExchanger(conf core.ManagerConfiguration)
5 there are 2 api methods:
	5.1 Exchanger.Exchanger.Subscribe(ch chan []*Exchange, refreshInterval time.Duration, targetCodes []string, referenceCode string):
	example:
	ch := make(chan []*Exchange)
	go exchanger.Exchanger.Subscribe(ch, 5, []string{"DASH", "ETC", "EOS", "WAVES", "STEEM", "BTS", "ETH"}, "USDT")
	for {
		select {
		case ex := <-ch:
			for _, v := range exx.Tickers {
					fmt.Println(exx.name, v.symbol(), v.Rate)
			}
		}
	}

	5.2 xchanger.Exchanger. GetRates(timeStamp time.Time, exchangeName string, targetCode string, referecies []string) []*Ticker:
	example: v := b.GetRates(time.Now().Add(-4 * time.Minute), "BINANCE", "BTS", []string{"BTC", "USDT"})
