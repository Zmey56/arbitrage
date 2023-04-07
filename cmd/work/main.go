package main

import (
	"github.com/Zmey56/arbitrage/pkg/p2p2stepsoneexchange"
	"github.com/Zmey56/arbitrage/pkg/p2pbinance"
	"github.com/Zmey56/arbitrage/pkg/workingbinance"
	"log"
	"time"
)

func main() {

	count := 0

	for {
		fiats := []string{"USD"}
		values := []string{"Payeer", "Advcash"}

		for _, fiat := range fiats {
			log.Println("FIAT", fiat)
			log.Println("Binance")

			paramUser := workingbinance.GetParam(fiat)

			paramUser.PublisherType = "null"
			//paramUser.TransAmount = "10000"
			//paramUser.PayTypes = []string{"TinkoffNew", "RosBankNew"}
			//log.Println("COUNT", count)
			//if count%2 == 0 {
			//	paramUser.PayTypes = []string{"RosBankNew"}
			//} else {
			//	paramUser.PayTypes = []string{"RUBfiatbalance", "Payeer", "Advcash"}
			//}

			paramUser.PayTypes = []string{values[count%len(values)]}

			log.Println(paramUser)

			log.Println("P2P3stepsTakerTaker")

			p2pbinance.P2P3stepsTakerTaker(fiat, paramUser)

			time.Sleep(1 * time.Minute)

			log.Println("P2P3stepsTakerMaker")

			p2pbinance.P2P3stepsTakerMaker(fiat, paramUser)

			time.Sleep(1 * time.Minute)

			log.Println("P2P2stepsBinance")

			p2p2stepsoneexchange.P2P2stepsBinance(fiat, paramUser)

			time.Sleep(1 * time.Minute)

		}
		count++
	}
}
