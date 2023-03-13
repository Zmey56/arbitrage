package workingokx

import (
	"encoding/json"
	"fmt"
	"github.com/Zmey56/arbitrage/pkg/getinfookx"
	"log"
	"os"
)

func GetParamOKX(fiat string) getinfookx.ParametersOKX {
	_, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	paramUser := getinfookx.ParametersOKX{}
	file_path := ""
	switch fiat {
	case "AED":
		file_path = fmt.Sprintf("cmd/enterparam/paramokx/%s.json", fiat)
	case "AMD":
		file_path = fmt.Sprintf("cmd/enterparam/paramokx/%s.json", fiat)
	case "ARS":
		file_path = fmt.Sprintf("cmd/enterparam/paramokx/%s.json", fiat)
	case "AZN":
		file_path = fmt.Sprintf("cmd/enterparam/paramokx/%s.json", fiat)
	case "EUR":
		file_path = fmt.Sprintf("cmd/enterparam/paramokx/%s.json", fiat)
	case "GEL":
		file_path = fmt.Sprintf("cmd/enterparam/paramokx/%s.json", fiat)
	case "KZT":
		file_path = fmt.Sprintf("cmd/enterparam/paramokx/%s.json", fiat)
	case "RUB":
		file_path = fmt.Sprintf("cmd/enterparam/paramokx/%s.json", fiat)
	case "TRY":
		file_path = fmt.Sprintf("cmd/enterparam/paramokx/%s.json", fiat)
	case "UAH":
		file_path = fmt.Sprintf("cmd/enterparam/paramokx/%s.json", fiat)
	case "UZS":
		file_path = fmt.Sprintf("cmd/enterparam/paramokx/%s.json", fiat)
	case "VND":
		file_path = fmt.Sprintf("cmd/enterparam/paramokx/%s.json", fiat)
	default:
		log.Printf("For %v don't have parametr\n", fiat)
	}
	file, err := os.ReadFile(file_path)
	if err != nil {
		log.Println("Can't read file with parameters", err)
	}
	_ = json.Unmarshal([]byte(file), &paramUser)

	return paramUser
}
