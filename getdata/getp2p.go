package getdata

import "net/url"

// link for get p2p data
const linkp2pgetdata = "https://p2p.binance.com/bapi/c2c/v2/friendly/c2c/adv/search"

// data for url linkp2pgetdata
var datagetp2p = url.Values{
	"asset":         {},
	"fiat":          {},
	"merchantCheck": {"False"},
	"page":          {},
	"payTypes":      {},
	"publisherType": {"None"},
	"rows":          {"10"},
	"tradeType":     {},
}

// New type for names of fiats and crypto
type Fiat string
type Crypto string

// Function for creating header
func GetDataP2P(f Fiat, c Crypto) {
	p := datagetp2p
}
