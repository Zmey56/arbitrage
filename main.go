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
	paramUser := getParam("RUB")
	paramUser_2 := getParam("RUB_2")
	for {

		start := time.Now()

		working.P2P3stepsTakerTaker(fiat, paramUser)

		log.Println(fiat, time.Since(start), "\n")

		time.Sleep(20 * time.Second)

		start_2 := time.Now()

		working.P2P3stepsTakerMaker(fiat, paramUser_2)

		log.Println(fiat, time.Since(start_2), "\n")

		time.Sleep(20 * time.Second)

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
	case "RUB_2":
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
