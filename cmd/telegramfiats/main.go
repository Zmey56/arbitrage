package main

import (
	"github.com/Zmey56/arbitrage/pkg/Interexchange"
	"github.com/Zmey56/arbitrage/pkg/p2pbinance"
	"github.com/Zmey56/arbitrage/pkg/p2phuobi"
	"github.com/Zmey56/arbitrage/pkg/workingbinance"
	"github.com/Zmey56/arbitrage/pkg/workinghuobi"
	"log"
	"time"
)

func main() {

	fiats := []string{"AED", "AMD", "AZN", "ARS", "EUR", "GEL", "KZT", "RUB", "TRY", "UAH", "UZS"}
	for {
		for _, fiat := range fiats {
			//time_total := time.Now()

			paramUserB := workingbinance.GetParam(fiat)
			//log.Printf("ParamsHuobi %+v", paramUserB)

			start := time.Now()

			p2pbinance.P2P3stepsTakerTaker(fiat, paramUserB)

			log.Println(fiat, "TakerTakerBinance", time.Since(start), "\n")

			time.Sleep(10 * time.Second)

			start_2 := time.Now()

			p2pbinance.P2P3stepsTakerMaker(fiat, paramUserB)

			log.Println(fiat, "TakerMakerBinance", time.Since(start_2), "\n")

			time.Sleep(10 * time.Second)

			if fiat != "AMD" && fiat != "AZN" {
				paranUserH := workinghuobi.GetParamHuobi(fiat)
				start_3 := time.Now()

				p2phuobi.P2P3stepsTakerTakerHuobi(fiat, paranUserH)

				log.Println(fiat, "TakerTakerHuobi", time.Since(start_3), "\n")

				time.Sleep(10 * time.Second)

				start_4 := time.Now()

				p2phuobi.P2P3stepsTakerMakerHuobi(fiat, paranUserH)

				log.Println(fiat, "TakerMakerHuobi", time.Since(start_4), "\n")

				time.Sleep(10 * time.Second)

				start_5 := time.Now()

				Interexchange.P2P2stepsBinanceHuobiTT(fiat, paramUserB, paranUserH)

				log.Println(fiat, "TakerMakerBinanceHuobi", time.Since(start_5), "\n")

				time.Sleep(10 * time.Second)

				start_6 := time.Now()

				Interexchange.P2P2stepsBinanceHuobiTM(fiat, paramUserB, paranUserH)

				log.Println(fiat, "TakerMakerBinanceHuobi", time.Since(start_6), "\n")

				time.Sleep(10 * time.Second)

			}
			//log.Println(fiat, "TOTAL TIME", time.Since(time_total), "\n")
		}
	}
}
