package result

import (
	"fmt"
	"github.com/Zmey56/arbitrage/pkg/getdatahuobi"
	"github.com/Zmey56/arbitrage/pkg/getinfobinance"
	"log"
	"strings"
)

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

func CheckResultSaveSend(howone, howtwo string, boarder int, per float64, profitResult ResultP2P) {
	if (profitResult.TotalAdvBuy > 0) && (profitResult.TotalAdvSell > 0) {
		log.Printf("%+v\n\n", profitResult)
		how_market := fmt.Sprintf("3steps_%s%s", string(profitResult.User.FirstUser[0]), string(profitResult.User.ThirdUser[0]))
		SaveResultJsonFile(profitResult.Fiat, profitResult, how_market)
		if profitResult.Profit && (profitResult.ProfitPercet >= per) &&
			(profitResult.TotalAdvBuy >= boarder) && (profitResult.TotalAdvSell >= boarder) {
			log.Printf("Profit - %s, ProfitPercet - %v, per - %v, TotalAdvBuy - %v, TotalAdvSell - %v, border - %v",
				profitResult.Profit, profitResult.ProfitPercet, per, profitResult.TotalAdvBuy, profitResult.TotalAdvSell, boarder)
			FormatMessageAndSend(profitResult)
		}
	}
}

func CheckResultSaveSend2Steps(profitResult ResultP2P2steps, border int) {
	if (profitResult.AdvToalSell > 0) && (profitResult.AdvToalSell > 0) {
		//log.Printf("%+v\n\n", profitResult)
		how_market := ""
		if profitResult.Merchant {
			how_market = "2steps_merchant"
		} else {
			how_market = "2steps"
		}

		SaveResultJsonFile2steps(profitResult.FiatUnit, profitResult, how_market)

		if (profitResult.AdvToalBuy >= border) &&
			(profitResult.AdvToalSell >= border) {

			FormatMessageAndSend2steps(profitResult)
		}
	}
}
