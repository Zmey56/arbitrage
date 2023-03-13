package getdataokx

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func GetCurrencyOKX(fiat string) []string {
	pathcurrency := fmt.Sprintf("data/dataokx/%s/%s_pair.json", fiat, fiat)
	file, _ := os.Open(pathcurrency)
	defer file.Close()
	decoder := json.NewDecoder(file)
	var data map[string][]string
	decoder.Decode(&data)
	currency := []string{}
	for i, _ := range data {
		currency = append(currency, strings.ToUpper(i))
	}
	return currency
}
