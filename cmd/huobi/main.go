package main

import (
	"github.com/Zmey56/arbitrage/pkg/Interexchange"
	"github.com/Zmey56/arbitrage/pkg/workingbinance"
	"github.com/Zmey56/arbitrage/pkg/workinghuobi"
)

func main() {
	fiat := "RUB"

	paramUserB := workingbinance.GetParam(fiat)
	paramUserH := workinghuobi.GetParamHuobi(fiat)

	Interexchange.P2P3stepsTTBBH(fiat, paramUserB, paramUserH)
	//Interexchange.P2P3stepsTMBHH(fiat, paramUserH, paramUserB)
	//Interexchange.P2P3stepsTMHBB(fiat, paramUserB, paramUserH)
	//Interexchange.P2P3stepsTMHHB(fiat, paramUserH, paramUserB)
	//Interexchange.P2P3stepsTTBBH(fiat, paramUserB, paramUserH)
	//Interexchange.P2P3stepsTTBHH(fiat, paramUserH, paramUserB)
	//Interexchange.P2P3stepsTTHBB(fiat, paramUserB, paramUserH)
	//Interexchange.P2P3stepsTTHHB(fiat, paramUserH, paramUserB)

	//workinghuobi.InputCommandLineHuobi("RUB")

	//log.Println(getdata.GetAssets("RUB"))

	//workinghuobi.GetParametrsHuobi()

	//paramUser := workinghuobi.GetParamHuobi("RUB")
	//
	//p2phuobi.P2P2stepsHuobi("RUB", paramUser)

	//for _, j := range crypto {
	//	p2phuobi.P2P3stepsTakerTakerHuobi(j)
	//}
	//pair := []string{"btcusdc",
	//	"btcusdt",
	//	"trxbtc",
	//	"ethbtc",
	//	"htbtc",
	//	"btcusdd",
	//	"eosbtc",
	//	"xrpbtc",
	//	"ltcbtc"}
	//
	//log.Println(getdatahuobi.GetRatePairHuobi(pair))

	//p2phuobi.P2P3stepsTakerTakerHuobi("RUB", paramUser)

	//p := "TRXUSDT"
	//a := "USDT"
	//var assetSell string
	//if strings.HasPrefix(p, a) {
	//	//transAmountSecond = (transAmountFirst * pair_rate[p])
	//	assetSell = p[len(a)-1:]
	//} else {
	//	//transAmountSecond = (transAmountFirst / pair_rate[p])
	//	assetSell = p[:(len(p) - len(a))]
	//}
	//log.Println(assetSell)
}
