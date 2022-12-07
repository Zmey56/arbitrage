package main

import (
	"fmt"
	"github.com/Zmey56/arbitrage/Interact"
	"github.com/Zmey56/arbitrage/working"
	"time"
)

const (
	fiat       = "RUB" //chose user
	tradeTypeB = "Buy"
	tradeTypeS = "Sell"
	asset      = "USDT"
)

func main() {

	tmp := Interact.InputCommandLine(fiat)
	fmt.Println("\n")
	//fmt.Println(payTypes)
	fmt.Println("\n")

	//fmt.Println(tmp)
	//
	//fmt.Println(getinfobinance.GetDataP2P(asset, fiat, tradeTypeB, tmp.TransAmount, tmp.PayTypes))

	start := time.Now()
	defer func() {
		fmt.Println(time.Since(start), "\n")
	}()

	working.P2P3steps(fiat, tmp)

	//test := []string{"BNBRUB", "ETHRUB", "SHIBRUB", "USDTRUB", "BTCRUB", "BUSDRUB"}
	////
	//fmt.Println(getinfobinance.GetRatePair(test))

	//assets := getinfobinance.GetAssets(fiat)
	//
	//getinfobinance.GetListSymbols(assets, fiat)

	//payTypes, transAmoount := working.InputCommandLine(fiat)
	//fmt.Println("\n")
	//fmt.Println(payTypes)
	//fmt.Println("\n")
	//
	//working.P2P2steps(fiat, payTypes, transAmoount)

	//
	//firststep := getdata.GetDataP2P(asset, fiat, tradeTypeB, transAmount, payTypes)
	//fmt.Println(firststep.Advertisers.NickName)

	//tmp := getinfobinance.GetListSymbols(asset)
	//
	//for i, j := range tmp {
	//	fmt.Println(i, " - ", j)
	//}
	//tmp := getinfobinance.GetPeymontMethods(fiat)
	//getinfobinance.SavePaymentToJSON(tmp)

	//data.GetPaymentFromJSON(fiat)
}
