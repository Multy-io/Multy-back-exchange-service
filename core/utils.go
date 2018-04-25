package core

func TrimLeftChars(s string, n int) string {
	m := 0
	for i := range s {
		if m >= n {
			return s[i:]
		}
		m++
	}
	return s[:0]
}

//type MapTickers struct {
//	tickers map[string]Ticker
//}
//
//func (b MapTickers) copy() MapTickers {
//	tickers := map[string]Ticker{}
//	for k, v := range b.tickers {
//		tickers[k] = v
//	}
//	return  MapTickers{tickers}
//}