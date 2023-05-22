package downloaddata

import (
	"github.com/Zmey56/arbitrage/pkg/getdatahuobi"
	"github.com/Zmey56/arbitrage/pkg/result"
	"github.com/Zmey56/arbitrage/pkg/workinghuobi"
	"math"
	"strconv"
)

func DownloadDataHuobi(fiat string) {
	assetsH := getdatahuobi.GetCurrencyHuobi(fiat)
	coinidmap := workinghuobi.GetCoinIDHuobo(fiat)

	for _, asset := range assetsH {
		// Create channel
		buyCh := make(chan getdatahuobi.Huobi)
		sellCh := make(chan getdatahuobi.Huobi)

		// Get first result
		resultBuy := getdatahuobi.GetDataP2PHuobiVer2(coinidmap[asset], coinidmap[fiat], "BUY", 1)
		resultSell := getdatahuobi.GetDataP2PHuobiVer2(coinidmap[asset], coinidmap[fiat], "SELL", 1)
		timePathBuy := strconv.FormatInt(resultBuy.TimeData, 10)
		timePathSell := strconv.FormatInt(resultSell.TimeData, 10)

		// Calculating the number of pages to get data
		resultBuyTotal := float64(resultBuy.TotalCount)
		resultSellTotal := float64(resultSell.TotalCount)
		pageBuy := int(math.Ceil(resultBuyTotal / 10))
		pageSell := int(math.Ceil(resultSellTotal / 10))

		// Creating arrays to store results
		arrayBuy := make([]getdatahuobi.Huobi, pageBuy)
		arraySell := make([]getdatahuobi.Huobi, pageSell)

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
					buyCh <- getdatahuobi.GetDataP2PHuobiVer2(coinidmap[asset], coinidmap[fiat], "BUY", i)
				}
				close(buyCh)
			}()
		}

		if pageSell > 1 {
			go func() {
				for i := 2; i <= pageSell; i++ {
					sellCh <- getdatahuobi.GetDataP2PHuobiVer2(coinidmap[asset], coinidmap[fiat], "SELL", i)
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

		result.SaveAllData("BUY", fiat, asset, timePathBuy, arrayBuy)
		result.SaveAllData("SELL", fiat, asset, timePathSell, arraySell)

	}

}
