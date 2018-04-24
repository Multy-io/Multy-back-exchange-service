package exchangeRates

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Appscrunch/Multy-back-exchange-service/currencies"
	"github.com/KristinaEtc/slf"
	_ "github.com/KristinaEtc/slflog"
	_ "github.com/lib/pq"
)

var log = slf.WithContext("exchange-rates")

type DBConfiguration struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type DbExchange struct {
	name    string
	Tickers []DbTicker
}

type DbTicker struct {
	TargetCurrency    currencies.Currency
	ReferenceCurrency currencies.Currency
	Rate              float64
	TimpeStamp        time.Time
	isCalculated 		bool
}

type DbManager struct {
	db *sql.DB
}

type DbRate struct {
	exchangeTitle string
	targetCode    string
	referenceCode string
	timeStamp     time.Time
	rate          float64
}

func NewDbManager(configuration DBConfiguration) *DbManager {
	manager := DbManager{}
	manager.db = manager.connectDb(configuration)
	return &manager
}

func (b *DbManager) connectDb(configuration DBConfiguration) *sql.DB {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		configuration.User, configuration.Password, configuration.Name)
	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		log.Errorf("connectDb:DbManager:sql.Open %v", err.Error())
	} else {
		log.Infof("Db connected")
	}
	return db
	//defer db.Close()
}

func (b *DbManager) FillDb(withExchanges []*DbExchange) {

	for _, exchange := range withExchanges {
		for _, ticker := range exchange.Tickers {
			b.insertSaRate(exchange.name, ticker.TargetCurrency, ticker.ReferenceCurrency, strconv.FormatFloat(ticker.Rate, 'f', 8, 64), ticker.isCalculated)
		}
	}
	b.fillRateFromSA()
}

//
//func (b *DbManager) insert(exchange *DbExchange) {
//	//fmt.Println("# Inserting values")
//
//	_, err := b.db.Exec("INSERT INTO exchanges(title,create_date) VALUES($1,$2) ON CONFLICT DO NOTHING;", exchange.name, time.Now())
//	//rows.Close()
//	checkErr(err)
//	//b.db.
//	//fmt.Println("inserted rows:", rows)
//}
//
//func (b *DbManager) insertCurrency(currency currencies.Currency) {
//	//fmt.Println("# Inserting values")
//
//	_, err := b.db.Exec("INSERT INTO currencies(code, title, create_date, native_id) VALUES($1,$2,$3,$4) ON CONFLICT DO NOTHING;", currency.CurrencyCode(), currency.CurrencyName(), time.Now(), currency)
//	//rows.Close()
//	checkErr(err)
//	//b.db.
//	//fmt.Println("inserted rows:", rows)
//}
//
func (b *DbManager) insertSaRate(exchange_title string, target_currency currencies.Currency, reference_currency currencies.Currency, rateString string, isCalculated bool) {
	rate, _ := strconv.ParseFloat(rateString, 64)
	_, err := b.db.Exec("INSERT INTO sa_rates(exchange_title, target_title, target_code, target_native_id, reference_title, reference_code, reference_native_id, time_stamp, rate, is_calculated) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10);", exchange_title, target_currency.CurrencyName(), target_currency.CurrencyCode(), target_currency, reference_currency.CurrencyName(), reference_currency.CurrencyCode(), reference_currency, time.Now(), rate, isCalculated)
	if err != nil {
		log.Errorf("DbManager:insertSaRate:b.db.Exec %v", err.Error())
	}
	//b.db.
}

func (b *DbManager) fillRateFromSA() {
	_, err := b.db.Exec("SELECT fill_rates()")
	if err != nil {
		log.Errorf("DbManager:fillRateFromSA:b.db.Exec %v", err.Error())
	}
}

func (b *DbManager) getRates(timeStamp time.Time, exchangeTitle string, targetCode string, refereciesCodes []string) []DbRate {
	var s = StringSlice{}
	s = refereciesCodes
	rows, err := b.db.Query("SELECT * from getRates($1, $2, $3, $4)", timeStamp, exchangeTitle, targetCode, s)
	if err != nil {
		log.Errorf("DbManager:getRates:b.db.Query %v", err.Error())
	}

	var dbRates = []DbRate{}

	for rows.Next() {
		dbRate := DbRate{}
		var exchange_title string
		var target_code string
		var reference_code string
		var time_stamp time.Time
		var rate float64
		err = rows.Scan(&exchange_title, &target_code, &reference_code, &time_stamp, &rate)
		if err != nil {
			log.Errorf("DbManager:getRates:rows.Scan %v", err.Error())
		}
		dbRate.exchangeTitle = exchange_title
		dbRate.targetCode = target_code
		dbRate.referenceCode = reference_code
		dbRate.timeStamp = time_stamp
		dbRate.rate = rate
		//fmt.Println("exchange_title | target_code | reference_code | time_stamp | rate")
		//fmt.Println(exchange_title, target_code, reference_code, time_stamp, rate)

		dbRates = append(dbRates, dbRate)
	}
	rows.Close()
	return dbRates
}

type StringSlice []string

func (stringSlice StringSlice) Value() (driver.Value, error) {
	var quotedStrings []string
	for _, str := range stringSlice {
		quotedStrings = append(quotedStrings, strconv.Quote(str))
	}
	value := fmt.Sprintf("{ %s }", strings.Join(quotedStrings, ","))
	return value, nil
}
