package getdata

import (
	"encoding/json"
	"fmt"
	"os"
)

const rubfiatpair = "data/RUB_pair_2022_11_28.json"

func GetPairFromJSON(fiat string) map[string][]string {
	pair := ""
	switch fiat {
	case "RUB":
		pair = rubfiatpair
	default:
		fmt.Printf("For %v don't have payments methods\n", fiat)
	}
	jsonfile, err := os.ReadFile(pair)
	if err != nil {
		panic(err)
	}
	var result map[string][]string
	_ = json.Unmarshal(jsonfile, &result)

	return result
}
