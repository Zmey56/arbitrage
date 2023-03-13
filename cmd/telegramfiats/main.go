package main

import (
	"fmt"
	"github.com/Zmey56/arbitrage/pkg/p2pbinance"
	"github.com/Zmey56/arbitrage/pkg/p2phuobi"
	"github.com/Zmey56/arbitrage/pkg/p2pinterexchange"
	"github.com/Zmey56/arbitrage/pkg/workingbinance"
	"github.com/Zmey56/arbitrage/pkg/workinghuobi"
	"log"
	"strconv"
	"time"
)

func main() {
	//
	//fiats := []string{"AED", "AMD", "AZN", "ARS", "EUR", "GEL", "KZT", "RUB", "TRY", "UAH", "UZS"}
	//fiats := []string{"USD"}
	fiats := []string{"EUR", "RUB", "UAH", "USD"}
	count := 0
	tmp_multi := 0
	multi := []float32{0.1, 0.5, 1}

	for {
		for _, fiat := range fiats {
			time_total := time.Now()

			paramUserH := workinghuobi.GetParamHuobi(fiat)
			paramUserB := workingbinance.GetParam(fiat)

			paramUserB.Border = 10
			paramUserH.Border = 10

			if count%2 == 0 {
				paramUserB.PublisherType = "merchant"
				paramUserH.IsMerchant = "true"
			} else {
				paramUserB.PublisherType = "null"
				paramUserH.IsMerchant = "false"
			}

			if count < 3 {
				tmp_multi = count
			} else {
				tmp_multi = (count - 3) % 3
			}

			amount_B, _ := strconv.Atoi(paramUserB.TransAmount)

			paramUserB.TransAmount = fmt.Sprintf("%f", float32(amount_B)*multi[tmp_multi])
			paramUserH.Amount = fmt.Sprintf("%f", float32(amount_B)*multi[tmp_multi])

			log.Printf("ParamsBinance %+v", paramUserB)
			log.Printf("ParamsHuobi %+v", paramUserH)

			start := time.Now()

			log.Println("P2P2stepsBinance")

			p2pbinance.P2P2stepsBinance(fiat, paramUserB)

			log.Println("P2P3stepsTakerTaker")

			p2pbinance.P2P3stepsTakerTaker(fiat, paramUserB)

			log.Println(fiat, "P2P3stepsTakerTaker", time.Since(start), "\n")

			time.Sleep(10 * time.Second)

			start_2 := time.Now()

			log.Println("P2P3stepsTakerMaker")

			p2pbinance.P2P3stepsTakerMaker(fiat, paramUserB)

			log.Println(fiat, "P2P3stepsTakerMaker", time.Since(start_2), "\n")

			time.Sleep(10 * time.Second)

			start_3 := time.Now()

			log.Println("P2P2stepsHuobi")

			p2phuobi.P2P2stepsHuobi(fiat, paramUserH)

			p2phuobi.P2P3stepsTakerTakerHuobi(fiat, paramUserH)

			log.Println(fiat, "P2P3stepsTakerTakerHuobi", time.Since(start_3), "\n")

			time.Sleep(10 * time.Second)
			//
			start_4 := time.Now()

			p2phuobi.P2P3stepsTakerMakerHuobi(fiat, paramUserH)

			log.Println(fiat, "P2P3stepsTakerMakerHuobi", time.Since(start_4), "\n")

			time.Sleep(10 * time.Second)

			start_5 := time.Now()

			log.Println("P2P2stepsBinanceHuobiTT")

			p2pinterexchange.P2P2stepsBinanceHuobiTT(fiat, paramUserB)

			p2pinterexchange.P2P3stepsTTBBH(fiat, paramUserB, paramUserH)

			log.Println(fiat, "P2P3stepsTTBBH", time.Since(start_5), "\n")

			time.Sleep(10 * time.Second)

			start_6 := time.Now()

			log.Println("P2P2stepsBinanceHuobiTM")

			p2pinterexchange.P2P2stepsBinanceHuobiTM(fiat, paramUserB)

			p2pinterexchange.P2P2stepsBinanceHuobiTM(fiat, paramUserB)

			log.Println(fiat, "P2P2stepsBinanceHuobiTM", time.Since(start_6), "\n")

			time.Sleep(10 * time.Second)

			start_7 := time.Now()

			log.Println("P2P2stepsHuobiBinanceTT")

			p2pinterexchange.P2P2stepsHuobiBinanceTT(fiat, paramUserH)

			p2pinterexchange.P2P3stepsTMBBH(fiat, paramUserB, paramUserH)

			log.Println(fiat, "P2P3stepsTMBBH", time.Since(start_7), "\n")

			time.Sleep(10 * time.Second)

			start_8 := time.Now()

			log.Println("P2P2stepsHuobiBinanceTM")

			p2pinterexchange.P2P2stepsHuobiBinanceTM(fiat, paramUserH)

			p2pinterexchange.P2P3stepsTMBHH(fiat, paramUserH, paramUserB)

			log.Println(fiat, "P2P3stepsTMBHH", time.Since(start_8), "\n")

			time.Sleep(10 * time.Second)

			start_9 := time.Now()

			p2pinterexchange.P2P3stepsTTBHH(fiat, paramUserH, paramUserB)

			log.Println(fiat, "P2P3stepsTTBHH", time.Since(start_9), "\n")

			time.Sleep(10 * time.Second)

			start_10 := time.Now()

			p2pinterexchange.P2P3stepsTMHBB(fiat, paramUserB, paramUserH)

			log.Println(fiat, "P2P3stepsTMHBB", time.Since(start_10), "\n")

			time.Sleep(10 * time.Second)

			start_11 := time.Now()

			p2pinterexchange.P2P3stepsTTHBB(fiat, paramUserB, paramUserH)

			log.Println(fiat, "P2P3stepsTTHBB", time.Since(start_11), "\n")

			time.Sleep(10 * time.Second)

			start_12 := time.Now()

			p2pinterexchange.P2P3stepsTTHHB(fiat, paramUserH, paramUserB)

			log.Println(fiat, "P2P3stepsTTHHB", time.Since(start_12), "\n")

			time.Sleep(10 * time.Second)

			start_13 := time.Now()

			p2pinterexchange.P2P3stepsTMHHB(fiat, paramUserH, paramUserB)

			log.Println(fiat, "P2P3stepsTMHHB", time.Since(start_13), "\n")

			time.Sleep(10 * time.Second)

			count++

			log.Println(fiat, "TOTAL TIME", time.Since(time_total), "COUNT:", count, "\n")
		}
	}
}
