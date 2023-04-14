package main

import (
	"fmt"
	"github.com/Zmey56/arbitrage/pkg/p2pbinance"
	"github.com/Zmey56/arbitrage/pkg/p2phuobi"
	"github.com/Zmey56/arbitrage/pkg/p2pokx"
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
	//fiats := []string{"EUR", "RUB", "UAH", "USD"}
	count := 0
	tmp_multi := 0
	multi := []float32{0.1, 0.5, 1}

	for {
		time_total := time.Now()
		for _, fiat := range fiats {

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

			//Binance

			start01 := time.Now()

			p2pbinance.P2P3stepsTakerMaker(fiat, paramUserB)

			log.Println(fiat, "P2P3stepsBinanceTM", time.Since(start01), "\n")

			time.Sleep(10 * time.Second)

			start02 := time.Now()

			p2pbinance.P2P3stepsTakerTaker(fiat, paramUserB)

			log.Println(fiat, "P2P2stepsBinanceTT", time.Since(start02), "\n")

			//Huobi

			start_21 := time.Now()

			p2phuobi.P2P3stepsTakerMakerHuobi(fiat, paramUserH)

			log.Println(fiat, "P2P3stepsHuobiTM", time.Since(start_21), "\n")

			time.Sleep(10 * time.Second)

			start_22 := time.Now()

			p2phuobi.P2P3stepsTakerTakerHuobi(fiat, paramUserH)

			log.Println(fiat, "P2P3stepsHuobiTT", time.Since(start_22), "\n")

			//OKX

			start_31 := time.Now()

			p2pokx.P2P3stepsTakerMakerOKX(fiat, paramUserO)

			log.Println(fiat, "P2P3stepsOKXTM", time.Since(start_31), "\n")

			time.Sleep(10 * time.Second)

			start_32 := time.Now()

			p2pokx.P2P3stepsTakerTakerOKX(fiat, paramUserO)

			log.Println(fiat, "P2P3stepsOKXTT", time.Since(start_32), "\n")

			<-time.After(10 * time.Second)

		}
		log.Println("TOTAL TIME", time.Since(time_total), "COUNT:", count, "\n")

		count++
	}
}
