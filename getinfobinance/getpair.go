package getinfobinance

import (
	"encoding/json"
	"fmt"
	"os"
)

func GetPairFromJSON(fiat string) map[string][]string {
	pair := ""
	switch fiat {
	case "AED":
		pair = fmt.Sprintf("data/%s/%s_pair.json", fiat, fiat)
	case "AMD":
		pair = fmt.Sprintf("data/%s/%s_pair.json", fiat, fiat)
	case "AZN":
		pair = fmt.Sprintf("data/%s/%s_pair.json", fiat, fiat)
	case "EUR":
		pair = fmt.Sprintf("data/%s/%s_pair.json", fiat, fiat)
	case "GEL":
		pair = fmt.Sprintf("data/%s/%s_pair.json", fiat, fiat)
	case "KZT":
		pair = fmt.Sprintf("data/%s/%s_pair.json", fiat, fiat)
	case "RUB":
		pair = fmt.Sprintf("data/%s/%s_pair.json", fiat, fiat)
	case "TRY":
		pair = fmt.Sprintf("data/%s/%s_pair.json", fiat, fiat)
	case "UAH":
		pair = fmt.Sprintf("data/%s/%s_pair.json", fiat, fiat)
	case "UZS":
		pair = fmt.Sprintf("data/%s/%s_pair.json", fiat, fiat)
	default:
		fmt.Printf("For %v don't have para\n", fiat)
	}
	//switch fiat {
	//case "RUB":
	//	pair = "data/RUB/RUB_pair.json"
	//case "GEL":
	//	pair = "data/GEL/GEL_payment.json"
	//default:
	//	fmt.Printf("For %v don't have payments methods\n", fiat)
	//}
	jsonfile, err := os.ReadFile(pair)
	if err != nil {
		panic(err)
	}
	var result map[string][]string
	_ = json.Unmarshal(jsonfile, &result)

	return result
}
