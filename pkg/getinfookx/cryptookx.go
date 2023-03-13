// get coin and pair. return json with coin and pair

package getinfookx

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func GetCoinOKX(fiat string) {
	url := fmt.Sprintf("https://www.okx.com/v3/c2c/currency/pairs?type=2&quote=%s", fiat)

	resp, err := http.Get(url)
	if err != nil {
		log.Println("trable with get response from OKX", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Can't read body", err)
	}

	var co CoinOKX

	json.Unmarshal(body, &co)

	allcoin := make([]string, len(co.Data))
	for i, j := range co.Data {
		allcoin[i] = j.BaseCurrency
	}

	//create all pair
	pairAll := GetPair()

	finalpair := make(map[string][]string)

	for i := 0; i < len(allcoin); i++ {
		for j := i + 1; j < len(allcoin); j++ {
			tmp := fmt.Sprintf("%s-%s", allcoin[i], allcoin[j])
			tmp_rev := fmt.Sprintf("%s-%s", allcoin[j], allcoin[i])
			if findSlice(tmp, pairAll) {
				log.Println("TMP", tmp)
				finalpair[allcoin[i]] = append(finalpair[allcoin[i]], tmp)
				finalpair[allcoin[j]] = append(finalpair[allcoin[j]], tmp)
			} else if findSlice(tmp_rev, pairAll) {
				log.Println("TMP_REV", tmp_rev)
				finalpair[allcoin[i]] = append(finalpair[allcoin[i]], tmp_rev)
				finalpair[allcoin[j]] = append(finalpair[allcoin[j]], tmp_rev)
			} else {
				continue
			}
		}
	}
	log.Println("finalpair", finalpair)

	name_json := fmt.Sprintf("data/dataokx/%s/%s_pair.json", fiat, fiat)
	log.Println("name_json", name_json)
	jsonStr, err := json.MarshalIndent(finalpair, "", " ")
	if err != nil {
		log.Printf("Error: %s", err.Error())
	}

	_ = os.WriteFile(name_json, jsonStr, 0644)
}

func GetPair() []string {
	url := fmt.Sprintf("https://www.okx.com/priapi/v5/public/simpleProduct?instType=SPOT")

	resp, err := http.Get(url)
	if err != nil {
		log.Println("trable with get response from OKX", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Can't read body", err)
	}

	var pair PairOKX

	json.Unmarshal(body, &pair)

	pairArray := make([]string, len(pair.Data))

	for i, j := range pair.Data {
		pairArray[i] = j.InstID
	}
	return pairArray
}

func findSlice(pair string, allPair []string) bool {
	for _, value := range allPair {
		if value == pair {
			fmt.Println("Value found in slice!")
			return true
		}
	}
	return false
}
