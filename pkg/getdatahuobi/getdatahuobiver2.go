package getdatahuobi

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func GetDataP2PHuobiVer2(asset, fiat int, tradeType string, page int) Huobi {

	params := url.Values{}
	params.Set("coinId", strconv.Itoa(asset))  //+
	params.Set("currency", strconv.Itoa(fiat)) //+
	if tradeType == "SELL" {
		params.Set("tradeType", "buy")  //+
		params.Set("acceptOrder", "-1") //+
	} else {
		params.Set("tradeType", "sell") //+
		params.Set("acceptOrder", "0")  //+
	}
	params.Set("tradeType", tradeType)
	params.Set("currPage", strconv.Itoa(page)) //+
	params.Set("payMethod", "0")               //+

	params.Set("country", "")          //+
	params.Set("blockType", "general") //+
	params.Set("online", "1")          //+
	params.Set("range", "0")           //+

	params.Set("amount", "") //+

	params.Set("onlyTradable", "false") //+
	params.Set("isFollowed", "false")   //+

	resulthuobi := Huobi{}

	url := ("https://otc-api.trygofast.com/v1/data/trade-market" + "?" + params.Encode())
	//log.Println("URL", url)

	var err error

	for {
		defer func() {
			if r := recover(); r != nil {
				if r == "connection reset by peer" {
					log.Println("An error occured 'connection reset by peer', reconecting...")
					time.Sleep(time.Second * 1)
				} else {
					// Handling other errors
					log.Println("An error occured:", r)
					time.Sleep(time.Second * 1)
				}
			}
		}()
		resulthuobi, err = requestOrdersP2PHuobiVer2(url)
		//log.Println("resulthuobi", resulthuobi)
		if err != nil {
			if err.Error() == "connection reset by peer" {
				// reconecting
				panic("connection reset by peer")
			} else {
				log.Println("Error:", err)
			}
		} else {
			break
		}
	}

	if len(resulthuobi.Data) > 0 {
		return resulthuobi
	} else {
		return Huobi{}
	}
}

func requestOrdersP2PHuobiVer2(j string) (Huobi, error) {

	for {
		backoff := 5 * time.Second
		response, err := http.Get(j)
		if err != nil {
			return Huobi{}, err
		}

		defer response.Body.Close()

		if err != nil {
			return Huobi{}, err
		}
		if response.StatusCode != http.StatusTooManyRequests {
			return parsingJsonHuobiVer2(response.Body), nil
		}
		log.Println("Too many requests, backing off for", backoff)
		time.Sleep(backoff)
	}

}

func parsingJsonHuobiVer2(r io.Reader) Huobi {
	var result Huobi

	body, _ := io.ReadAll(r)
	err := json.Unmarshal([]byte(body), &result)

	if err != nil {
		log.Println("Error unmarshal json URL Huobi:", err, string(body))
	}

	return result
}
