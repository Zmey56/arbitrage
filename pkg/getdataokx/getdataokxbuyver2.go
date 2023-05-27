package getdataokx

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func GetDataP2POKXBuyVer2(asset, fiat, tradeType string, page int) OKXBuy {

	params := url.Values{}
	params.Set("side", "sell")
	params.Set("paymentMethod", "all")

	params.Set("userType", "all")
	params.Set("hideOverseasVerificationAds", "false")
	params.Set("sortType", "price_asc")
	params.Set("urlId", "0")
	params.Set("limit", "10000")
	params.Set("cryptoCurrency", strings.ToLower(asset))
	params.Set("fiatCurrency", strings.ToLower(fiat))
	params.Set("currentPage", strconv.Itoa(page)) //;
	params.Set("numberPerPage", "10")

	resultokxbuy := OKXBuy{}

	url_buy := ("https://www.okx.com/v3/c2c/tradingOrders/getMarketplaceAdsPrelogin" + "?" + params.Encode())
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
		resultokxbuy, err = requestOrdersP2POKXBuyVer2(url_buy)
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

	if len(resultokxbuy.Data.Sell) > 0 {
		return resultokxbuy
	} else {
		return OKXBuy{}
	}
}

func requestOrdersP2POKXBuyVer2(j string) (OKXBuy, error) {

	for {
		backoff := 5 * time.Second

		req, err := http.NewRequest("GET", j, nil)
		if err != nil {
			fmt.Println("Error creating HTTP request:", err)
			return OKXBuy{}, err
		}
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36")
		req.Header.Set("Accept-Language", "en-US,en;q=0.5")

		// Send the HTTP request using the default HTTP client
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Println("Error sending HTTP request:", err)
			return OKXBuy{}, err
		}

		// Make sure the response body is closed when we're done with it
		defer resp.Body.Close()

		dateStr := resp.Header.Get("Date")
		if dateStr == "" {
			log.Println("Data header not found")
			return OKXBuy{}, err
		}
		layout := http.TimeFormat
		date, err := time.Parse(layout, dateStr)
		if err != nil {
			log.Println("Invalid Date header")
			return OKXBuy{}, err
		}
		unixTime := date.Unix()

		if resp.StatusCode != http.StatusTooManyRequests {
			return parsingJsonOKXBuyVer2(resp.Body, unixTime), nil
		}
		log.Println("Too many requests, backing off for", backoff)
		time.Sleep(backoff)
	}

}

func parsingJsonOKXBuyVer2(r io.Reader, ut int64) OKXBuy {
	var result OKXBuy

	body, _ := io.ReadAll(r)
	err := json.Unmarshal([]byte(body), &result)
	if err != nil {
		log.Println("Error unmarshal json URL Huobi:", err, string(body))
	}

	result.TimeData = ut

	return result
}
