package main

import (
	"encoding/json"
	"fmt"
	"github.com/Zmey56/arbitrage/getdata"
	"github.com/Zmey56/arbitrage/interact"
	"github.com/Zmey56/arbitrage/working"
	"log"
	"os"
	"time"
)

const (
	fiat       = "RUB" //chose user
	tradeTypeB = "Buy"
	tradeTypeS = "Sell"
	asset      = "USDT"
)

func main() {
	//getNewFiat("UAH")
	//getNewFiat("RUB")
	//getNewFiat("GEL")
	//getNewFiat("KZT")
	//getNewFiat("AED")
	//getNewFiat("AMD")
	//getNewFiat("AZN")
	//getNewFiat("EUR")
	//getNewFiat("UZS")
	//getNewFiat("TRY")

	//fiats := []string{"AED", "AMD", "AZN", "EUR", "GEL", "KZT", "RUB", "TRY", "UAH", "UZS"}
	//fiats := []string{"RUB"}
	//
	//for _, f := range fiats {
	//	fmt.Println("FIAT:", f)
	//	interact.InputCommandLine(f)
	//}

	//fmt.Println(tmp)
	paramUser := interact.InputCommandLine(fiat)
	//
	for {
		//paramUser := interact.Parameters{}
		//for _, fiat := range fiats {
		//paramUser = getParam(fiat)

		start := time.Now()

		working.P2P3stepsTakerTaker(fiat, paramUser)

		log.Println(fiat, time.Since(start), "\n")

		time.Sleep(60 * time.Second)

		//}
	}

}

func getNewFiat(fiat string) {
	tmp := getdata.GetPeymontMethods(fiat)
	getdata.SavePaymentToJSON(tmp, fiat) //need to correctthis function
	a := getdata.GetAssets(fiat)
	getdata.GetListSymbols(a, fiat)
}

// temporary function
func getParam(fiat string) interact.Parameters {
	paramUser := interact.Parameters{}
	file_path := ""
	switch fiat {
	case "AED":
		file_path = fmt.Sprintf("data/paramUser%s.json", fiat)
	case "AMD":
		file_path = fmt.Sprintf("data/paramUser%s.json", fiat)
	case "AZN":
		file_path = fmt.Sprintf("data/paramUser%s.json", fiat)
	case "EUR":
		file_path = fmt.Sprintf("data/paramUser%s.json", fiat)
	case "GEL":
		file_path = fmt.Sprintf("data/paramUser%s.json", fiat)
	case "KZT":
		file_path = fmt.Sprintf("data/paramUser%s.json", fiat)
	case "RUB":
		file_path = fmt.Sprintf("data/paramUser%s.json", fiat)
	case "TRY":
		file_path = fmt.Sprintf("data/paramUser%s.json", fiat)
	case "UAH":
		file_path = fmt.Sprintf("data/paramUser%s.json", fiat)
	case "UZS":
		file_path = fmt.Sprintf("data/paramUser%s.json", fiat)
	default:
		fmt.Printf("For %v don't have parametr\n", fiat)
	}
	file, err := os.ReadFile(file_path)
	if err != nil {
		log.Println("Can't read file with parameters", err)
	}
	_ = json.Unmarshal([]byte(file), &paramUser)

	return paramUser
}
