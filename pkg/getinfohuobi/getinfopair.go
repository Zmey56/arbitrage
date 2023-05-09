package getinfohuobi

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

// GetPairFromJSONPairPairPair return map with pair of map[cryptoPair]cryptoPair|cryptoPair
func GetPairFromJSONPairPairPair(asset string) map[string][]string {
	pair := ""
	switch asset {
	case "USDT":
		pair = fmt.Sprintf("data/datahuobi/%s/%s_pair_pair_pair.json", asset, asset)
	default:
		log.Printf("For %v don't have para\n", asset)
	}

	jsonfile, err := os.ReadFile(pair)
	if err != nil {
		panic(err)
	}
	var result map[string][]string
	_ = json.Unmarshal(jsonfile, &result)

	return result
}
