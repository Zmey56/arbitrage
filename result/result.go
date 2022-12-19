package result

import (
	"fmt"
	"github.com/Zmey56/arbitrage/getinfobinance"
	"strings"
	"time"
)

type ResultP2P struct {
	Profit          bool
	DataTime        time.Time
	Fiat            string
	AssetsBuy       string
	PriceAssetsBuy  float64
	PaymentBuy      []string
	LinkAssetsBuy   string
	Pair            string
	PricePair       float64
	LinkMarket      string
	AssetsSell      string
	PriceAssetsSell float64
	PaymentSell     []string
	LinkAssetsSell  string
	ProfitValue     float64
	ProfitPercet    float64
}

func ReturnLinkMarket(a, p string) string {
	var pair string
	if strings.HasPrefix(p, a) {
		pair = a + "_" + p[len(a):]
	} else {
		pair = p[:(len(p)-len(a))] + "_" + a
	}
	return fmt.Sprintf("https://www.binance.com/en/trade/%v?_from=markets", pair)
}

func PaymentMetods(a getinfobinance.AdvertiserAdv) []string {
	payMethods := []string{}
	for _, tm := range a.Adv.TradeMethods {
		payMethods = append(payMethods, tm.TradeMethodName)
	}
	return payMethods
}
