package main

import (
	"fmt"
	"github.com/Zmey56/arbitrage/working"
)

const (
	fiat       = "RUB"  //chose user
	asset      = "USDT" //chose from available
	tradeTypeB = "Buy"
	tradeTypeS = "Sell"
)

type TransAmountCollect int

func main() {

	payTypes, transAmoount := working.InputCommandLine(fiat)
	fmt.Println("\n")
	fmt.Println(payTypes)
	fmt.Println("\n")

	working.P2P2steps(fiat, payTypes, transAmoount)

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
