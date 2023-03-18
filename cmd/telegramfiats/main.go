package main

import (
	"fmt"
	"github.com/Zmey56/arbitrage/pkg/p2p2stepsoneexchange"
	"github.com/Zmey56/arbitrage/pkg/workingbinance"
	"github.com/Zmey56/arbitrage/pkg/workinghuobi"
	"github.com/Zmey56/arbitrage/pkg/workingokx"
	"log"
	"strconv"
	"time"
)

func main() {
	//
	fiats := []string{"AED", "EUR", "GEL", "KZT", "RUB", "TRY", "UAH", "USD"}
	count := 0
	tmp_multi := 0
	multi := []float32{0.1, 0.5, 1}

	for {
		for _, fiat := range fiats {
			time_total := time.Now()

			paramUserH := workinghuobi.GetParamHuobi(fiat)
			paramUserB := workingbinance.GetParam(fiat)
			paramUserO := workingokx.GetParamOKX(fiat)

			paramUserB.Border = 10
			paramUserH.Border = 10
			paramUserO.Border = 10

			if count%2 == 0 {
				paramUserB.PublisherType = "merchant"
				paramUserH.IsMerchant = "true"
				paramUserO.IsMerchant = "true"
			} else {
				paramUserB.PublisherType = "null"
				paramUserH.IsMerchant = "false"
				paramUserO.IsMerchant = "false"
			}

			if count < 3 {
				tmp_multi = count
			} else {
				tmp_multi = (count - 3) % 3
			}

			amount_B, _ := strconv.Atoi(paramUserB.TransAmount)

			paramUserB.TransAmount = fmt.Sprintf("%f", float32(amount_B)*multi[tmp_multi])
			paramUserH.Amount = fmt.Sprintf("%f", float32(amount_B)*multi[tmp_multi])
			paramUserO.Amount = fmt.Sprintf("%f", float32(amount_B)*multi[tmp_multi])

			log.Printf("ParamsBinance %+v", paramUserB)
			log.Printf("ParamsHuobi %+v", paramUserH)
			log.Printf("ParamsOKX %+v", paramUserO)

			start := time.Now()

			p2p2stepsoneexchange.P2P2stepsBinance(fiat, paramUserB)

			log.Println(fiat, "P2P2stepsBinance", time.Since(start), "\n")

			time.Sleep(10 * time.Second)

			start_2 := time.Now()

			p2p2stepsoneexchange.P2P2stepsHuobi(fiat, paramUserH)

			log.Println(fiat, "P2P2stepsHuobi", time.Since(start_2), "\n")

			time.Sleep(10 * time.Second)

			start_3 := time.Now()

			p2p2stepsoneexchange.P2P2stepsOKX(fiat, paramUserO)

			log.Println(fiat, "P2P2stepsOKX", time.Since(start_3), "\n")

			log.Println(fiat, "TOTAL TIME", time.Since(time_total), "COUNT:", count, "Amount", "\n")

		}
		count++
	}
}
