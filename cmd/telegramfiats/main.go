package main

import (
	"fmt"
	"github.com/Zmey56/arbitrage/pkg/downloaddata"
	"github.com/Zmey56/arbitrage/pkg/getinfohuobi"
	"github.com/Zmey56/arbitrage/pkg/getinfookx"
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
	fiats := []string{"AED", "EUR", "GEL", "KZT", "RUB", "TRY", "UAH", "USD"}
	//fiats := []string{"EUR", "RUB", "UAH", "USD"}
	//fiats := []string{"RUB"}
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

			//Binance

			downloaddata.DownloadDataBinance(fiat)

			time.Sleep(1 * time.Second)

			p2pbinance.P2P3stepsTakerMaker(fiat, paramUserB)

			time.Sleep(1 * time.Second)

			p2pbinance.P2P3stepsTakerTaker(fiat, paramUserB)

			time.Sleep(1 * time.Second)

			p2p2stepsoneexchange.P2P2stepsBinance(fiat, paramUserB)

			time.Sleep(1 * time.Second)

			rubAmongExchange2stepsBinance(fiat, paramUserB)

			//Huobi

			downloaddata.DownloadDataHuobi(fiat)

			time.Sleep(1 * time.Second)

			p2phuobi.P2P3stepsTakerMakerHuobi(fiat, paramUserH)

			time.Sleep(1 * time.Second)

			p2phuobi.P2P3stepsTakerTakerHuobi(fiat, paramUserH)

			time.Sleep(1 * time.Second)

			p2p2stepsoneexchange.P2P2stepsHuobi(fiat, paramUserH)

			time.Sleep(1 * time.Second)

			rubAmongExchange2stepsHuobi(fiat, paramUserH)

			//OKX

			downloaddata.DownloadDataOKX(fiat)

			time.Sleep(1 * time.Second)

			p2pokx.P2P3stepsTakerMakerOKX(fiat, paramUserO)

			time.Sleep(1 * time.Second)

			p2pokx.P2P3stepsTakerTakerOKX(fiat, paramUserO)

			time.Sleep(1 * time.Second)

			p2p2stepsoneexchange.P2P2stepsOKX(fiat, paramUserO)

			time.Sleep(1 * time.Second)

			rubAmongExchange2stepsOKX(fiat, paramUserO)

			<-time.After(1 * time.Second)

		}
		log.Println("TOTAL TIME", time.Since(time_total), "COUNT:", count, "\n")

		count++
	}
}

func rubAmongExchange2stepsBinance(fiat string, paramUserB workingbinance.ParametersBinance) {
	p2p2stepsamongexchange.P2P2stepsBinanceHuobiTM(fiat, paramUserB)

	p2p2stepsamongexchange.P2P2stepsBinanceHuobiTT(fiat, paramUserB)

	p2p2stepsamongexchange.P2P2stepsBinanceOKXTM(fiat, paramUserB)

	p2p2stepsamongexchange.P2P2stepsBinanceOKXTT(fiat, paramUserB)
}

func rubAmongExchange2stepsHuobi(fiat string, paramUserH getinfohuobi.ParametersHuobi) {
	p2p2stepsamongexchange.P2P2stepsHuobiBinanceTM(fiat, paramUserH)

	p2p2stepsamongexchange.P2P2stepsHuobiBinanceTT(fiat, paramUserH)

	p2p2stepsamongexchange.P2P2stepsHuobiOKXTM(fiat, paramUserH)

	p2p2stepsamongexchange.P2P2stepsHuobiOKXTT(fiat, paramUserH)
}

func rubAmongExchange2stepsOKX(fiat string, paramUserO getinfookx.ParametersOKX) {

	p2p2stepsamongexchange.P2P2stepsOKXBinanceTM(fiat, paramUserO)

	p2p2stepsamongexchange.P2P2stepsOKXBinanceTT(fiat, paramUserO)

	p2p2stepsamongexchange.P2P2stepsOKXHuobiTM(fiat, paramUserO)

	p2p2stepsamongexchange.P2P2stepsOKXHuobiTT(fiat, paramUserO)
}
