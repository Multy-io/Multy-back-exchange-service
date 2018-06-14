package currencies

import (
	"strings"
)

type CurrencyPair struct {
	TargetCurrency   Currency
	ReferenceCurrency Currency
}

func (b *CurrencyPair) Symbol() string {
	return b.TargetCurrency.CurrencyCode() + "-" + b.ReferenceCurrency.CurrencyCode()
}

func (b *CurrencyPair) IsEqualTo(pair CurrencyPair) bool {
	return b.TargetCurrency == pair.TargetCurrency && b.ReferenceCurrency == pair.ReferenceCurrency
}


var DefaultReferenceCurrencies = []Currency{Tether, Bitcoin}

type Currency int

const (

	NotAplicable           Currency = -99
	Bitcoin                Currency = 0
	Testnet                Currency = 1
	Litecoin               Currency = 2
	Dogecoin               Currency = 3
	Reddcoin               Currency = 4
	Dash                   Currency = 5
	Peercoin               Currency = 6
	Namecoin               Currency = 7
	Feathercoin            Currency = 8
	Counterparty           Currency = 9
	Blackcoin              Currency = 10
	NuShares               Currency = 11
	NuBits                 Currency = 12
	Mazacoin               Currency = 13
	Viacoin                Currency = 14
	ClearingHouse          Currency = 15
	Rubycoin               Currency = 16
	Groestlcoin            Currency = 17
	Digitalcoin            Currency = 18
	Cannacoin              Currency = 19
	DigiByte               Currency = 20
	OpenAssets             Currency = 21
	Monacoin               Currency = 22
	Clams                  Currency = 23
	Primecoin              Currency = 24
	Neoscoin               Currency = 25
	Jumbucks               Currency = 26
	ziftrCOIN              Currency = 27
	Vertcoin               Currency = 28
	NXT                    Currency = 29
	Burst                  Currency = 30
	MonetaryUnit           Currency = 31
	Zoom                   Currency = 32
	Vpncoin                Currency = 33
	CanadaeCoin            Currency = 34
	ShadowCash             Currency = 35
	ParkByte               Currency = 36
	Pandacoin              Currency = 37
	StartCOIN              Currency = 38
	MOIN                   Currency = 39
	Expanse                Currency = 40
	Einsteinium            Currency = 41
	Decred                 Currency = 42
	NEM                    Currency = 43
	Particl                Currency = 44
	Argentum               Currency = 45
	Libertas               Currency = 46
	Poswcoin               Currency = 47
	Shreeji                Currency = 48
	GlobalCurrencyReserve  Currency = 49
	Novacoin               Currency = 50
	Asiacoin               Currency = 51
	Bitcoindark            Currency = 52
	Dopecoin               Currency = 53
	Templecoin             Currency = 54
	AIB                    Currency = 55
	EDRCoin                Currency = 56
	Syscoin                Currency = 57
	Solarcoin              Currency = 58
	Smileycoin             Currency = 59
	Ether                  Currency = 60
	EtherClassic           Currency = 61
	Pesobit                Currency = 62
	Landcoin               Currency = 63
	OpenChain              Currency = 64
	Bitcoinplus            Currency = 65
	InternetofPeople       Currency = 66
	Nexus                  Currency = 67
	InsaneCoin             Currency = 68
	OKCash                 Currency = 69
	BritCoin               Currency = 70
	Compcoin               Currency = 71
	Crown                  Currency = 72
	BelaCoin               Currency = 73
	Compcoin2              Currency = 74
	FujiCoin               Currency = 75
	MIX                    Currency = 76
	Verge                  Currency = 77
	ElectronicGulden       Currency = 78
	ClubCoin               Currency = 79
	RichCoin               Currency = 80
	Potcoin                Currency = 81
	Quarkcoin              Currency = 82
	Terracoin              Currency = 83
	Gridcoin               Currency = 84
	Auroracoin             Currency = 85
	IXCoin                 Currency = 86
	Gulden                 Currency = 87
	BitBeanv               Currency = 88
	Bata                   Currency = 89
	Myriadcoin             Currency = 90
	BitSend                Currency = 91
	Unobtanium             Currency = 92
	MasterTrader           Currency = 93
	GoldBlocks             Currency = 94
	Saham                  Currency = 95
	Chronos                Currency = 96
	Ubiquoin               Currency = 97
	Evotion                Currency = 98
	SaveTheOcean           Currency = 99
	BigUp                  Currency = 100
	GameCredits            Currency = 101
	Dollarcoins            Currency = 102
	Zayedcoin              Currency = 103
	Dubaicoin              Currency = 104
	Stratis                Currency = 105
	Shilling               Currency = 106
	MarsCoin               Currency = 107
	Ubiq                   Currency = 108
	Pesetacoin             Currency = 109
	Neurocoin              Currency = 110
	ARK                    Currency = 111
	UltimateSecureCashMain Currency = 112
	Hempcoin               Currency = 113
	Linx                   Currency = 114
	Ecoin                  Currency = 115
	Denarius               Currency = 116
	Pinkcoin               Currency = 117
	PiggyCoin              Currency = 118
	Pivx                   Currency = 119
	Flashcoin              Currency = 120
	Zencash                Currency = 121
	Putincoin              Currency = 122
	BitZeny                Currency = 123
	Unify                  Currency = 124
	StealthCoin            Currency = 125
	BreakoutCoin           Currency = 126
	Vcash                  Currency = 127
	Monerov                Currency = 128
	Voxels                 Currency = 129
	NavCoin                Currency = 130
	FactomFactoids         Currency = 131
	Factom                 Currency = 132
	Zcash                  Currency = 133
	Lisk                   Currency = 134
	Steem                  Currency = 135
	ZCoin                  Currency = 136
	Rootstock              Currency = 137
	Giftblock              Currency = 138
	RealPointCoin          Currency = 139
	LBRY                   Currency = 140
	Komodo                 Currency = 141
	bisqToken              Currency = 142
	Riecoin                Currency = 143
	Ripple                 Currency = 144
	BitcoinCash            Currency = 145
	Neblio                 Currency = 146
	ZClassic               Currency = 147
	StellarLumens          Currency = 148
	NoLimitCoin2           Currency = 149
	WhaleCoin              Currency = 150
	EuropeCoin             Currency = 151
	Diamond                Currency = 152
	Bytom                  Currency = 153
	Biocoin                Currency = 154
	Whitecoin              Currency = 155
	BitcoinGold            Currency = 156
	Bitcoin2x              Currency = 157
	SuperSkynet            Currency = 158
	TOACoin                Currency = 159
	Bitcore                Currency = 160
	Adcoin                 Currency = 161
	Bridgecoin             Currency = 162
	Ellaism                Currency = 163
	Pirl                   Currency = 164
	RaiBlocks              Currency = 165
	Vivo                   Currency = 166
	Firstcoin              Currency = 167
	Helleniccoin           Currency = 168
	BUZZ                   Currency = 169
	Ember                  Currency = 170
	Hcash                  Currency = 171
	HTMLCOIN               Currency = 172
	AskCoin                Currency = 223
	Smartcash              Currency = 224
	ZenProtocol            Currency = 258
	MemCoin                Currency = 333
	NEO                    Currency = 888
	BitcoinDiamond         Currency = 999
	Defcoin                Currency = 1337
	Cardano                Currency = 1815
	RootstockTestnet       Currency = 37310
	Golos                  Currency = 37311
	BitShares              Currency = 37312
	Waves                  Currency = 37313
	EOS                    Currency = 37314
	Tether                 Currency = 37315
	SouthKoreanWon         Currency = 57316
)

var CurrencyNames = map[Currency]string{
	NotAplicable:           "NotAplicable",
	Bitcoin:               "Bitcoin",
	Testnet:               "Testnet",
	Litecoin:              "Litecoin",
	Dogecoin:              "Dogecoin",
	Reddcoin:              "Reddcoin",
	Dash:                  "Dash",
	Peercoin:              "Peercoin",
	Namecoin:              "Namecoin",
	Feathercoin:           "Feathercoin",
	Counterparty:          "Counterparty",
	Blackcoin:             "Blackcoin",
	NuShares:              "NuShares",
	NuBits:                "NuBits",
	Mazacoin:              "Mazacoin",
	Viacoin:               "Viacoin",
	ClearingHouse:         "ClearingHouse",
	Rubycoin:              "Rubycoin",
	Groestlcoin:           "Groestlcoin",
	Digitalcoin:           "Digitalcoin",
	Cannacoin:             "Cannacoin",
	DigiByte:              "DigiByte",
	OpenAssets:            "OpenAssets",
	Monacoin:              "Monacoin",
	Clams:                 "Clams",
	Primecoin:             "Primecoin",
	Neoscoin:              "Neoscoin",
	Jumbucks:              "Jumbucks",
	ziftrCOIN:             "ziftrCOIN",
	Vertcoin:              "Vertcoin",
	NXT:                   "NXT",
	Burst:                 "Burst",
	MonetaryUnit:          "MonetaryUnit",
	Zoom:                  "Zoom",
	Vpncoin:               "Vpncoin",
	CanadaeCoin:           "CanadaeCoin",
	ShadowCash:            "ShadowCash",
	ParkByte:              "ParkByte",
	Pandacoin:             "Pandacoin",
	StartCOIN:             "StartCOIN",
	MOIN:                  "MOIN",
	Expanse:               "Expanse",
	Einsteinium:           "Einsteinium",
	Decred:                "Decred",
	NEM:                   "NEM",
	Particl:               "Particl",
	Argentum:              "Argentum",
	Libertas:              "Libertas",
	Poswcoin:              "Poswcoin",
	Shreeji:               "Shreeji",
	GlobalCurrencyReserve: "GlobalCurrencyReserve",
	Novacoin:              "Novacoin",
	Asiacoin:              "Asiacoin",
	Bitcoindark:           "Bitcoindark",
	Dopecoin:              "Dopecoin",
	Templecoin:            "Templecoin",
	AIB:                   "AIB",
	EDRCoin:               "EDRCoin",
	Syscoin:               "Syscoin",
	Solarcoin:             "Solarcoin",
	Smileycoin:            "Smileycoin",
	Ether:                 "Ether",
	EtherClassic:          "EtherClassic",
	Pesobit:               "Pesobit",
	Landcoin:              "Landcoin",
	OpenChain:             "OpenChain",
	Bitcoinplus:           "Bitcoinplus",
	InternetofPeople:      "InternetofPeople",
	Nexus:                 "Nexus",
	InsaneCoin:            "InsaneCoin",
	OKCash:                "OKCash",
	BritCoin:              "BritCoin",
	Compcoin2:             "Compcoin2",
	Crown:                 "Crown",
	BelaCoin:              "BelaCoin",
	Compcoin:              "Compcoin",
	FujiCoin:              "FujiCoin",
	MIX:                   "MIX",
	Verge:                 "Verge",
	ElectronicGulden:      "ElectronicGulden",
	ClubCoin:              "ClubCoin",
	RichCoin:              "RichCoin",
	Potcoin:               "Potcoin",
	Quarkcoin:             "Quarkcoin",
	Terracoin:             "Terracoin",
	Gridcoin:              "Gridcoin",
	Auroracoin:            "Auroracoin",
	IXCoin:                "IXCoin",
	Gulden:                "Gulden",
	BitBeanv:              "BitBeanv",
	Bata:                  "Bata",
	Myriadcoin:            "Myriadcoin",
	BitSend:               "BitSend",
	Unobtanium:            "Unobtanium",
	MasterTrader:          "MasterTrader",
	GoldBlocks:            "GoldBlocks",
	Saham:                 "Saham",
	Chronos:               "Chronos",
	Ubiquoin:              "Ubiquoin",
	Evotion:               "Evotion",
	SaveTheOcean:          "SaveTheOcean",
	BigUp:                 "BigUp",
	GameCredits:           "GameCredits",
	Dollarcoins:           "Dollarcoins",
	Zayedcoin:             "Zayedcoin",
	Dubaicoin:             "Dubaicoin",
	Stratis:               "Stratis",
	Shilling:              "Shilling",
	MarsCoin:              "MarsCoin",
	Ubiq:                  "Ubiq",
	Pesetacoin:            "Pesetacoin",
	Neurocoin:             "Neurocoin",
	ARK:                   "ARK",
	UltimateSecureCashMain: "UltimateSecureCashMain",
	Hempcoin:               "Hempcoin",
	Linx:                   "Linx",
	Ecoin:                  "Ecoin",
	Denarius:               "Denarius",
	Pinkcoin:               "Pinkcoin",
	PiggyCoin:              "PiggyCoin",
	Pivx:                   "Pivx",
	Flashcoin:              "Flashcoin",
	Zencash:                "Zencash",
	Putincoin:              "Putincoin",
	BitZeny:                "BitZeny",
	Unify:                  "Unify",
	StealthCoin:            "StealthCoin",
	BreakoutCoin:           "BreakoutCoin",
	Vcash:                  "Vcash",
	Monerov:                "Monerov",
	Voxels:                 "Voxels",
	NavCoin:                "NavCoin",
	FactomFactoids:         "FactomFactoids",
	Factom:                 "Factom",
	Zcash:                  "Zcash",
	Lisk:                   "Lisk",
	Steem:                  "Steem",
	ZCoin:                  "ZCoin",
	Rootstock:              "Rootstock",
	Giftblock:              "Giftblock",
	RealPointCoin:          "RealPointCoin",
	LBRY:                   "LBRY",
	Komodo:                 "Komodo",
	bisqToken:              "bisqToken",
	Riecoin:                "Riecoin",
	Ripple:                 "Ripple",
	BitcoinCash:            "BitcoinCash",
	Neblio:                 "Neblio",
	ZClassic:               "ZClassic",
	StellarLumens:          "StellarLumens",
	NoLimitCoin2:           "NoLimitCoin2",
	WhaleCoin:              "WhaleCoin",
	EuropeCoin:             "EuropeCoin",
	Diamond:                "Diamond",
	Bytom:                  "Bytom",
	Biocoin:                "Biocoin",
	Whitecoin:              "Whitecoin",
	BitcoinGold:            "BitcoinGold",
	Bitcoin2x:              "Bitcoin2x",
	SuperSkynet:            "SuperSkynet",
	TOACoin:                "TOACoin",
	Bitcore:                "Bitcore",
	Adcoin:                 "Adcoin",
	Bridgecoin:             "Bridgecoin",
	Ellaism:                "Ellaism",
	Pirl:                   "Pirl",
	RaiBlocks:              "RaiBlocks",
	Vivo:                   "Vivo",
	Firstcoin:              "Firstcoin",
	Helleniccoin:           "Helleniccoin",
	BUZZ:                   "BUZZ",
	Ember:                  "Ember",
	Hcash:                  "Hcash",
	HTMLCOIN:               "HTMLCOIN",
	AskCoin:                "AskCoin",
	Smartcash:              "Smartcash",
	ZenProtocol:            "ZenProtocol",
	MemCoin:                "MemCoin",
	NEO:                    "NEO",
	BitcoinDiamond:         "BitcoinDiamond",
	Defcoin:                "Defcoin",
	Cardano:                "Cardano",
	RootstockTestnet:       "RootstockTestnet",
	Golos:                  "Golos",
	BitShares:              "BitShares",
	Waves:                  "Waves",
	EOS:                    "EOS",
	Tether:                 "Tether",
	SouthKoreanWon: 		"SouthKoreanWon",
}

var CurrencyCodes = map[Currency]string{
	NotAplicable: "N.A.",
	Bitcoin:      "BTC",
	Litecoin:     "LTC",
	Dash:         "DASH",
	Ether:        "ETH",
	Golos:        "GOLOS",
	BitShares:    "BTS",
	Steem:        "STEEM",
	Waves:        "WAVES",
	BitcoinCash:  "BCH",
	EtherClassic: "ETC",
	EOS:          "EOS",
	Tether:       "USDT",
	SouthKoreanWon: "KRW",
}

func (currency Currency) CurrencyName() string {
	return CurrencyNames[currency]
}

func (currency Currency) CurrencyCode() string {
	return CurrencyCodes[currency]
}

func NewCurrency(currencyName string) Currency {
	var currencies = map[string]Currency{}
	for key, val := range CurrencyNames {
		currencies[strings.ToUpper(val)] = key
	}
	currency := currencies[strings.ToUpper(currencyName)]
	return currency
}

func NewCurrencyWithCode(currencyCodeString string) Currency {
	var currencies = map[string]Currency{}
	for key, val := range CurrencyCodes {
		currencies[strings.ToUpper(val)] = key
	}

	if currencyCodeString == "USD" {
		currencyCodeString = "USDT"
	}

	if currency, ok := currencies[strings.ToUpper(currencyCodeString)]; ok {
	 return currency
	} else {
		return NotAplicable
	}
}
