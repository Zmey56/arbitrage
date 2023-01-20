package main

import (
	"github.com/Zmey56/arbitrage/pkg/p2phuobi"
	"github.com/Zmey56/arbitrage/pkg/workinghuobi"
	"log"
)

func main() {

	//paramUser := getinfohuobi.ParametersHuobi{
	//	CoinId:       "2",
	//	Currency:     "2",
	//	TradeType:    "sell",
	//	CurrPage:     "1",
	//	PayMethod:    "0",
	//	AcceptOrder:  "0",
	//	Country:      "",
	//	BlockType:    "general",
	//	Online:       "1",
	//	Range:        "0",
	//	Intrange:     "",
	//	Amount:       "",
	//	IsThumbsUp:   "false",
	//	IsMerchant:   "false",
	//	IsTraded:     "false",
	//	OnlyTradable: "false",
	//	IsFollowed:   "false",
	//}

	//getinfohuobi.GetDataP2PHuobi(paramUser)

	fiats := []string{"AED", "ARS", "EUR", "GEL", "KZT", "RUB", "TRY", "UAH", "UZS", "VND"}

	//getinfobinance.GetPeymontMethodsBinance("VND")

	for _, i := range fiats {
		paranUserH := workinghuobi.GetParamHuobi(i)
		log.Println("FIAT", i)
		p2phuobi.P2P3stepsTakerTakerHuobi(i, paranUserH)
		p2phuobi.P2P3stepsTakerMakerHuobi(i, paranUserH)
		//	workinghuobi.InputCommandLineHuobi(i)
	}

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
