package main

import (
	"github.com/Zmey56/arbitrage/pkg/Interexchange"
	"github.com/Zmey56/arbitrage/pkg/workingbinance"
	"github.com/Zmey56/arbitrage/pkg/workinghuobi"
)

func main() {

	fiats := []string{"AED", "ARS", "EUR", "GEL", "KZT", "RUB", "TRY", "UAH", "UZS"}
	//fiat := "RUB"
	for _, fiat := range fiats {

		paramUserB := workingbinance.GetParam(fiat)
		paranUserH := workinghuobi.GetParamHuobi(fiat)

		Interexchange.P2P2stepsBinanceHuobiTT(fiat, paramUserB, paranUserH)
		Interexchange.P2P2stepsBinanceHuobiTM(fiat, paramUserB, paranUserH)
	}

}
