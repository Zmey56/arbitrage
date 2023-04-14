package main

import (
	"github.com/Zmey56/arbitrage/pkg/workingbinance"
	"github.com/Zmey56/arbitrage/pkg/workinghuobi"
	"github.com/Zmey56/arbitrage/pkg/workingokx"
	"log"
)

func main() {

	//fiat := "USD"
	//
	//getinfookx.GetCoinOKX(fiat)
	//getinfookx.GetPeymontMethodsOKX(fiat)
	//count := 0

	//fiats := []string{"AED", "EUR", "GEL", "KZT", "RUB", "TRY", "UAH", "USD"}
	fiats := []string{"RUB", "UAH"}

	for _, fiat := range fiats {
		log.Println("FIAT", fiat, "Market Binance")

		workingbinance.InputCommandLineBinance(fiat)

		log.Println("FIAT", fiat, "Market Huobi")

		workinghuobi.InputCommandLineHuobi(fiat)

		log.Println("FIAT", fiat, "Market OKX")

		workingokx.InputCommandLineOKX(fiat)

	}
}
