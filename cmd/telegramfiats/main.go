package main

import (
	"fmt"
	"github.com/Zmey56/arbitrage/pkg/p2p2stepsamongexchange"
	"github.com/Zmey56/arbitrage/pkg/p2p2stepsoneexchange"
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
	//fiats := []string{"AED", "EUR", "GEL", "KZT", "RUB", "TRY", "UAH", "USD"}
	fiats := []string{"EUR", "RUB", "UAH", "USD"}
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

			start := time.Now()

			//Binance

			p2p2stepsoneexchange.P2P2stepsBinance(fiat, paramUserB)

			log.Println(fiat, "P2P2stepsBinance", time.Since(start), "\n")

			start01 := time.Now()

			p2pbinance.P2P3stepsTakerMaker(fiat, paramUserB)

			log.Println(fiat, "P2P3stepsBinanceTM", time.Since(start01), "\n")

			start02 := time.Now()

			p2pbinance.P2P3stepsTakerTaker(fiat, paramUserB)

			log.Println(fiat, "P2P2stepsBinanceTT", time.Since(start02), "\n")

			//between Binance and Huobi

			start_4 := time.Now()

			p2p2stepsamongexchange.P2P2stepsBinanceHuobiTM(fiat, paramUserB)

			log.Println(fiat, "P2P2stepsBinanceHuobiTM", time.Since(start_4), "\n")

			start_5 := time.Now()

			p2p2stepsamongexchange.P2P2stepsBinanceHuobiTT(fiat, paramUserB)

			log.Println(fiat, "P2P2stepsBinanceHuobiTT", time.Since(start_5), "\n")

			start_7 := time.Now()

			p2p2stepsamongexchange.P2P2stepsHuobiBinanceTM(fiat, paramUserH)

			log.Println(fiat, "P2P2stepsHuobiBinanceTM", time.Since(start_7), "\n")

			start_8 := time.Now()

			p2p2stepsamongexchange.P2P2stepsHuobiBinanceTT(fiat, paramUserH)

			log.Println(fiat, "P2P2stepsHuobiBinanceTT", time.Since(start_8), "\n")

			start_2 := time.Now()

			//Huobi

			p2p2stepsoneexchange.P2P2stepsHuobi(fiat, paramUserH)

			log.Println(fiat, "P2P2stepsHuobi", time.Since(start_2), "\n")

			start_21 := time.Now()

			p2phuobi.P2P3stepsTakerMakerHuobi(fiat, paramUserH)

			log.Println(fiat, "P2P3stepsHuobiTM", time.Since(start_21), "\n")

			start_22 := time.Now()

			p2phuobi.P2P3stepsTakerTakerHuobi(fiat, paramUserH)

			log.Println(fiat, "P2P3stepsHuobiTT", time.Since(start_22), "\n")

			//between Huobi and OKX

			start_9 := time.Now()

			p2p2stepsamongexchange.P2P2stepsHuobiOKXTM(fiat, paramUserH)

			log.Println(fiat, "P2P2stepsHuobiOKXTM", time.Since(start_9), "\n")

			start_91 := time.Now()

			p2p2stepsamongexchange.P2P2stepsHuobiOKXTT(fiat, paramUserH)

			log.Println(fiat, "P2P2stepsHuobiOKXTT", time.Since(start_91), "\n")

			start_12 := time.Now()

			p2p2stepsamongexchange.P2P2stepsOKXHuobiTM(fiat, paramUserO)

			log.Println(fiat, "P2P2stepsOKXHuobiTM", time.Since(start_12), "\n")

			start_13 := time.Now()

			p2p2stepsamongexchange.P2P2stepsOKXHuobiTT(fiat, paramUserO)

			log.Println(fiat, "P2P2stepsOKXHuobiTT", time.Since(start_13), "\n")

			//OKX

			start_3 := time.Now()

			p2p2stepsoneexchange.P2P2stepsOKX(fiat, paramUserO)

			log.Println(fiat, "P2P2stepsOKX", time.Since(start_3), "\n")

			start_31 := time.Now()

			p2pokx.P2P3stepsTakerMakerOKX(fiat, paramUserO)

			log.Println(fiat, "P2P3stepsOKXTM", time.Since(start_31), "\n")

			start_32 := time.Now()

			p2pokx.P2P3stepsTakerTakerOKX(fiat, paramUserO)

			log.Println(fiat, "P2P3stepsOKXTT", time.Since(start_32), "\n")

			//between OKX and Binance

			start_6 := time.Now()

			p2p2stepsamongexchange.P2P2stepsBinanceOKXTM(fiat, paramUserB)

			log.Println(fiat, "P2P2stepsBinanceOKXTM", time.Since(start_6), "\n")

			start_61 := time.Now()

			p2p2stepsamongexchange.P2P2stepsBinanceOKXTT(fiat, paramUserB)

			log.Println(fiat, "P2P2stepsBinanceOKXTT", time.Since(start_61), "\n")

			start_10 := time.Now()

			p2p2stepsamongexchange.P2P2stepsOKXBinanceTM(fiat, paramUserO)

			log.Println(fiat, "P2P2stepsOKXBinanceTM", time.Since(start_10), "\n")

			start_11 := time.Now()

			p2p2stepsamongexchange.P2P2stepsOKXBinanceTT(fiat, paramUserO)

			log.Println(fiat, "P2P2stepsOKXBinanceTT", time.Since(start_11), "\n")

			<-time.After(1 * time.Minute)

		}
		log.Println("TOTAL TIME", time.Since(time_total), "COUNT:", count, "\n")

		count++
	}
}
