package getdataokx

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func GetDataP2POKXSellVer2(asset, fiat, tradeType string, page int) OKXSell {

	params := url.Values{}
	params.Set("side", "buy")
	params.Set("paymentMethod", "all")

	params.Set("userType", "all")
	params.Set("hideOverseasVerificationAds", "false")
	params.Set("sortType", "price_asc")
	params.Set("urlId", "0")
	params.Set("limit", "10000")
	params.Set("cryptoCurrency", strings.ToLower(asset))
	params.Set("fiatCurrency", strings.ToLower(fiat))
	params.Set("currentPage", "1") //;
	params.Set("numberPerPage", "10")

	resultokxsell := OKXSell{}

	url_sell := ("https://www.okx.com/v3/c2c/tradingOrders/getMarketplaceAdsPrelogin" + "?" + params.Encode())
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
		resultokxsell, err = requestOrdersP2POKXSellVer2(url_sell)
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

	if len(resultokxsell.Data.Buy) > 0 {
		//log.Println("Size for sell before return", len(resultokxsell.Data.Buy), "FIAT ", fiat, " Currency ", currency)
		return resultokxsell
	} else {
		return OKXSell{}
	}
}

func requestOrdersP2POKXSellVer2(j string) (OKXSell, error) {

	for {
		backoff := 5 * time.Second
		req, err := http.NewRequest("GET", j, nil)
		if err != nil {
			fmt.Println("Error creating HTTP request:", err)
			return OKXSell{}, err
		}
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36")
		req.Header.Set("Accept-Language", "en-US,en;q=0.5")

		// Send the HTTP request using the default HTTP client
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Println("Error sending HTTP request:", err)
			return OKXSell{}, err
		}
		//log.Println(resp.Body)

		// Make sure the response body is closed when we're done with it
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusTooManyRequests {
			return parsingJsonOKXSellVer2(resp.Body), nil
		}
		log.Println("Too many requests, backing off for", backoff)
		time.Sleep(backoff)
	}

}

func parsingJsonOKXSellVer2(r io.Reader) OKXSell {
	var result OKXSell

	body, _ := io.ReadAll(r)
	err := json.Unmarshal([]byte(body), &result)

	if err != nil {
		log.Println("Error unmarshal json URL Huobi:", err, string(body))
	}

	//log.Println("Size for sell", len(result.Data.Buy))
	//for i, j := range result.Data.Buy {
	//	log.Println(i, " - ", j.Price, j.NickName)
	//}

	return result
}
