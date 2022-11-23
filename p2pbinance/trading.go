package p2pbinance

import "github.com/Zmey56/arbitrage/getinfobinance"

func P2PTrading() (asset, fiat string, transAmount float64, payTypes []string) {
	assetget := getinfobinance.GetAssets(fiat)

	symbolasset := make([]string, len(assetget))

	//get only symbol without text
	for k, _ := range assetget {
		symbolasset = append(symbolasset, k)
	}

	//create pair

}
