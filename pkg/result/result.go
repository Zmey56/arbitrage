package result

import (
	"fmt"
	"github.com/Zmey56/arbitrage/pkg/getdatahuobi"
	"github.com/Zmey56/arbitrage/pkg/getinfobinance"
	"strings"
	"time"
)

type ResultP2P struct {
	Market struct {
		First  string
		Second string
		Third  string
	}
	Merchant struct {
		FirstMerch  bool
		SecondMerch bool
		ThirdMerch  bool
	}
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
	TotalAdvBuy     int
	TotalAdvSell    int
	AdvNoBuy        string
	AdvNoSell       string
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

func ReturnLinkMarketHuobi(a, p string) string {
	var pair string
	if strings.HasPrefix(p, a) {
		pair = a + "_" + p[len(a):]
	} else {
		pair = p[:(len(p)-len(a))] + "_" + a
	}
	return fmt.Sprintf("https://www.huobi.com/en-us/exchange/%v", pair)
}

func PaymentMetods(a getinfobinance.AdvertiserAdv) []string {
	payMethods := []string{}
	for _, tm := range a.Data[0].Adv.TradeMethods {
		payMethods = append(payMethods, tm.TradeMethodName)
	}
	return payMethods
}

func PaymentMetodsHuobi(a getdatahuobi.Huobi) []string {
	payMethods := []string{}
	for _, tm := range a.Data[0].PayMethods {
		payMethods = append(payMethods, tm.Name)
	}
	return payMethods
}
