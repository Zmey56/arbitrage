package result

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

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
