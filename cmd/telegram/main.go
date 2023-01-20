package main

import (
	"github.com/Zmey56/arbitrage/pkg/workingbinance"
	"log"
)

func main() {
	fiats := []string{"AED", "ARS", "EUR", "GEL", "KZT", "RUB", "TRY", "UAH", "UZS", "VND"}

	for _, j := range fiats {
		log.Println(j)
		workingbinance.InputCommandLineBinance(j)
	}

	//	http.DefaultTransport.(*http.Transport).MaxIdleConnsPerHost = 100
	//	paramUser := getParam("RUB")
	//	paramUser_2 := getParam("RUB_2")
	//	for {
	//
	//		start := time.Now()
	//
	//		p2pbinance.P2P3stepsTakerTaker(fiat, paramUser)
	//
	//		log.Println(fiat, time.Since(start), "\n")
	//
	//		time.Sleep(20 * time.Second)
	//
	//		start_2 := time.Now()
	//
	//		p2pbinance.P2P3stepsTakerMaker(fiat, paramUser_2)
	//
	//		log.Println(fiat, time.Since(start_2), "\n")
	//
	//		time.Sleep(20 * time.Second)
	//
	//	}
	//
}

//
//func getNewFiat(fiat string) {
//	tmp := getdata2.GetPeymontMethods(fiat)
//	getdata2.SavePaymentToJSON(tmp, fiat) //need to correctthis function
//	a := getdata2.GetAssets(fiat)
//	getdata2.GetListSymbols(a, fiat)
//}
//
//// temporary function
//func getParam(fiat string) interact.Parameters {
//	paramUser := interact.Parameters{}
//	file_path := ""
//	switch fiat {
//	case "AED":
//		file_path = fmt.Sprintf("data/paramUser%s.json", fiat)
//	case "AMD":
//		file_path = fmt.Sprintf("data/paramUser%s.json", fiat)
//	case "AZN":
//		file_path = fmt.Sprintf("data/paramUser%s.json", fiat)
//	case "EUR":
//		file_path = fmt.Sprintf("data/paramUser%s.json", fiat)
//	case "GEL":
//		file_path = fmt.Sprintf("data/paramUser%s.json", fiat)
//	case "KZT":
//		file_path = fmt.Sprintf("data/paramUser%s.json", fiat)
//	case "RUB":
//		file_path = fmt.Sprintf("data/paramUser%s.json", fiat)
//	case "RUB_2":
//		file_path = fmt.Sprintf("data/paramUser%s.json", fiat)
//	case "TRY":
//		file_path = fmt.Sprintf("data/paramUser%s.json", fiat)
//	case "UAH":
//		file_path = fmt.Sprintf("data/paramUser%s.json", fiat)
//	case "UZS":
//		file_path = fmt.Sprintf("data/paramUser%s.json", fiat)
//	default:
//		fmt.Printf("For %v don't have parametr\n", fiat)
//	}
//	file, err := os.ReadFile(file_path)
//	if err != nil {
//		log.Println("Can't read file with parameters", err)
//	}
//	_ = json.Unmarshal([]byte(file), &paramUser)
//
//	return paramUser
//}
