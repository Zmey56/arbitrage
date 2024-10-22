package downloaddata

import (
	"github.com/Zmey56/arbitrage/pkg/getdataokx"
	"github.com/Zmey56/arbitrage/pkg/result"
	"math"
)

func DownloadDataOKX(fiat string) {
	//start01 := time.Now()
	assetsO := getdataokx.GetCurrencyOKX(fiat)

	for _, asset := range assetsO {
		// Create channel
		buyCh := make(chan getdataokx.OKXBuy)
		sellCh := make(chan getdataokx.OKXSell)

		// Get first result
		resultBuy := getdataokx.GetDataP2POKXBuyVer2(asset, fiat, "BUY", 1)
		resultSell := getdataokx.GetDataP2POKXSellVer2(asset, fiat, "SELL", 1)

		// Calculating the number of pages to get data
		resultBuyTotal := float64(resultBuy.Data.Total)
		resultSellTotal := float64(resultSell.Data.Total)
		pageBuy := int(math.Ceil(resultBuyTotal / 10))
		pageSell := int(math.Ceil(resultSellTotal / 10))

		// Creating arrays to store results
		arrayBuy := make([]getdataokx.OKXBuy, pageBuy)
		arraySell := make([]getdataokx.OKXSell, pageSell)

		// Saving the first results to arrays
		if len(resultBuy.Data.Sell) > 0 {
			arrayBuy[0] = resultBuy
		}
		if len(resultSell.Data.Buy) > 0 {
			arraySell[0] = resultSell
		}

		// Launching goroutin to get results
		if pageBuy > 1 {
			go func() {
				for i := 2; i <= pageBuy; i++ {
					buyCh <- getdataokx.GetDataP2POKXBuyVer2(asset, fiat, "BUY", i)
				}
				close(buyCh)
			}()
		}

		if pageSell > 1 {
			go func() {
				for i := 2; i <= pageSell; i++ {
					sellCh <- getdataokx.GetDataP2POKXSellVer2(asset, fiat, "SELL", i)
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

	}
}
