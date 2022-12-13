package main

import (
	"fmt"
	"github.com/Zmey56/arbitrage/interact"
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

	tmp := interact.InputCommandLine(fiat)

	for {
		start := time.Now()

		working.P2P3steps(fiat, tmp)

		fmt.Println(time.Since(start), "\n")

		time.Sleep(60 * time.Second)
	}

}
