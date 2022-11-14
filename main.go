package main

import (
	"fmt"
	"github.com/Zmey56/arbitrage/inputvalue"
	"github.com/Zmey56/arbitrage/p2pbinance"
)

const (
	fiat       = "GEL"
	asset      = "USDT"
	tradeTypeB = "Buy"
	tradeTypeS = "Sell"
)

func main() {

	payTypes := inputvalue.InputCommandLine(fiat)
	fmt.Println("\n")
	fmt.Println(payTypes)
	fmt.Println("\n")

	tmp := p2pbinance.AdvertiserAdv{}

	tmp = p2pbinance.GetDataP2P(asset, fiat, tradeTypeB, payTypes)

	fmt.Println(len(tmp.Advertisers))
	fmt.Println(len(tmp.Advs))
	////
	////fmt.Printf("%+v\n\n", tmp.Advertisers)
	////fmt.Printf("%+v\n\n", tmp.Advs)

}
