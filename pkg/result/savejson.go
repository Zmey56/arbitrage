package result

import (
	"encoding/json"
	"fmt"
	"github.com/Zmey56/arbitrage/pkg/getdatahuobi"
	"github.com/Zmey56/arbitrage/pkg/getdataokx"
	"github.com/Zmey56/arbitrage/pkg/getinfobinance"
	"log"
	"os"
	"strings"
	"time"
)

// SaveAllData save all data to file
func SaveAllData(tradeType, fiat, asset, timePath string, data interface{}) {
	//const layout = "2006-01-02_15-04"
	resultPath := ""
	//t := time.Now()
	//timePath := t.Format(layout)
	nameFile := fmt.Sprintf("%s_%s_%s", strings.ToLower(fiat), strings.ToLower(asset), timePath)

	switch data.(type) {
	case []getinfobinance.Binance:
		if tradeType == "BUY" {
			resultPath = fmt.Sprintf("jsonresult/Binance/buy_%s.json", nameFile)
		} else {
			resultPath = fmt.Sprintf("jsonresult/Binance/sell_%s.json", nameFile)
		}
	case []getdatahuobi.Huobi:
		log.Println("HUOBI")
		if tradeType == "BUY" {
			resultPath = fmt.Sprintf("jsonresult/Huobi/buy_%s.json", nameFile)
		} else {
			resultPath = fmt.Sprintf("jsonresult/Huobi/sell_%s.json", nameFile)
		}
	case []getdataokx.OKXBuy:
		resultPath = fmt.Sprintf("jsonresult/OKX/buy_%s.json", nameFile)
	case []getdataokx.OKXSell:
		resultPath = fmt.Sprintf("jsonresult/OKX/sell_%s.json", nameFile)
	default:
		log.Println("I don't now this format")
	}

	//log.Println("resultPath", resultPath)
	file, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		panic(err)
	}
	os.WriteFile(resultPath, file, 0644)
}

func SaveResultJsonFile(fiat string, pr ResultP2P, how string) {
	current := time.Now()
	current.String()
	var path_save string
	path_save = fmt.Sprintf("jsonresult/%s/%s_%s_%s.json", pr.Market.First, fiat, current.Format("2006_01_02"), how)

	tmp_result := []ResultP2P{}
	if exists(path_save) {
		jsonfile, err := os.ReadFile(path_save)
		if err != nil {
			panic(err)
		}

		_ = json.Unmarshal(jsonfile, &tmp_result)
	}
	tmp_result = append(tmp_result, pr)
	f, err := os.OpenFile(path_save,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	file, _ := json.MarshalIndent(tmp_result, "", " ")
	os.WriteFile(path_save, file, 0666)
	defer f.Close()
}

func SaveResultJsonFile2steps(fiat string, pr ResultP2P2steps, how string) {
	current := time.Now()
	current.String()
	var path_save string
	path_save = fmt.Sprintf("jsonresult/%s/%s_%s_%s.json", pr.MarketOne, fiat, current.Format("2006_01_02"), how)

	tmp_result := []ResultP2P2steps{}
	if exists(path_save) {
		jsonfile, err := os.ReadFile(path_save)
		if err != nil {
			panic(err)
		}

		_ = json.Unmarshal(jsonfile, &tmp_result)
	}
	tmp_result = append(tmp_result, pr)
	f, err := os.OpenFile(path_save,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	file, _ := json.MarshalIndent(tmp_result, "", " ")
	os.WriteFile(path_save, file, 0666)
	defer f.Close()
}

func exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}
