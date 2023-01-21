package getdatahuobi

import (
	"encoding/json"
	"fmt"
	"github.com/Zmey56/arbitrage/pkg/getinfohuobi"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func GetDataP2PHuobi(fiat, currency int, tradeType string, paramUser getinfohuobi.ParametersHuobi) Huobi {

	params := url.Values{}
	params.Set("coinId", strconv.Itoa(fiat))
	params.Set("currency", strconv.Itoa(currency)) //find this
	params.Set("tradeType", tradeType)
	params.Set("currPage", "1")
	if len(paramUser.PayMethod) > 0 {
		params.Set("payMethod", paramUser.PayMethod)
	} else {
		params.Set("payMethod", "0")
	}
	// +
	if tradeType == "buy" {
		params.Set("acceptOrder", "0")
	} else {
		params.Set("acceptOrder", "-1")
	}
	params.Set("country", "")
	params.Set("blockType", "general")
	params.Set("online", "1")
	params.Set("range", "0")
	if len(paramUser.Amount) > 0 {
		params.Set("amount", paramUser.Amount)
	} else {
		params.Set("amount", "")
	}
	params.Set("isThumbsUp", "false")
	if len(paramUser.IsMerchant) > 0 {
		params.Set("isMerchant", paramUser.IsMerchant)
	} else {
		params.Set("isMerchant", "true")
	}
	params.Set("isTraded", "false")
	params.Set("onlyTradable", "false")
	params.Set("isFollowed", "false")
	resulthuobi := Huobi{}
	//log.Println("PARAMS HUOBI", params)

	url := ("https://otc-api.huobi.com/v1/data/trade-market" + "?" + params.Encode())
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
		resulthuobi, err = requestOrdersP2PHuobi(url)
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

func requestOrdersP2PHuobi(j string) (Huobi, error) {

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
			return ParsingJsonHuobi(response.Body), nil
		}
		fmt.Println("Too many requests, backing off for", backoff)
		time.Sleep(backoff)
	}

}

func ParsingJsonHuobi(r io.Reader) Huobi {
	var result Huobi

	body, _ := io.ReadAll(r)
	err := json.Unmarshal([]byte(body), &result)

	if err != nil {
		log.Println("Error unmarshal json URL Huobi:", err, string(body))
	}

	return result
}
