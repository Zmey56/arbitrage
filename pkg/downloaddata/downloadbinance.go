package downloaddata

import (
	"github.com/Zmey56/arbitrage/pkg/getdata"
	"github.com/Zmey56/arbitrage/pkg/getinfobinance"
	"github.com/Zmey56/arbitrage/pkg/result"
	"log"
	"math"
	"time"
)

func DownloadDataBinance(fiat string) {
	start01 := time.Now()
	assetsB := getdata.GetAssetsLocalBinance(fiat)

	for _, asset := range assetsB {
		start02 := time.Now()
		// Create channel
		buyCh := make(chan getinfobinance.Binance)
		sellCh := make(chan getinfobinance.Binance)

		// Get first result
		resultBuy := getdata.GetDataP2PBinanceVer2(asset, fiat, "BUY", 1)
		resultSell := getdata.GetDataP2PBinanceVer2(asset, fiat, "SELL", 1)

		// Calculating the number of pages to get data
		resultBuyTotal := float64(resultBuy.Total)
		resultSellTotal := float64(resultSell.Total)
		pageBuy := int(math.Ceil(resultBuyTotal / 10))
		pageSell := int(math.Ceil(resultSellTotal / 10))

		// Creating arrays to store results
		arrayBuy := make([]getinfobinance.Binance, pageBuy)
		arraySell := make([]getinfobinance.Binance, pageSell)

		// Saving the first results to arrays
		if len(resultBuy.Data) > 0 {
			arrayBuy[0] = resultBuy
		}
		if len(resultSell.Data) > 0 {
			arraySell[0] = resultSell
		}

		// Launching goroutin to get results
		if pageBuy > 1 {
			go func() {
				for i := 2; i <= pageBuy; i++ {
					buyCh <- getdata.GetDataP2PBinanceVer2(asset, fiat, "BUY", i)
				}
				close(buyCh)
			}()
		}

		if pageSell > 1 {
			go func() {
				for i := 2; i <= pageSell; i++ {
					sellCh <- getdata.GetDataP2PBinanceVer2(asset, fiat, "SELL", i)
				}
				close(sellCh)
			}()
		}

		// Getting results from channels and storing them in arrays
		if pageBuy > 1 {
			for i := 1; i < pageBuy; i++ {
				arrayBuy[i] = <-buyCh
			}
		}

		if pageSell > 1 {
			for i := 1; i < pageSell; i++ {
				arraySell[i] = <-sellCh
			}
		}

		result.SaveAllData("BUY", fiat, asset, arrayBuy)
		result.SaveAllData("SELL", fiat, asset, arraySell)

		log.Println("TIME", time.Since(start02), "\n")
	}

	log.Println("TIME", time.Since(start01), "\n")
	//time.Sleep(1 * time.Minute)
}
