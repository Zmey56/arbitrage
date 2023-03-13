package getdataokx

import (
	"encoding/json"
	"fmt"
	"github.com/Zmey56/arbitrage/pkg/getinfookx"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func GetDataP2POKXBuy(fiat, currency string, pu getinfookx.ParametersOKX) OKXBuy {

	params := url.Values{}
	params.Set("side", "sell")
	params.Set("paymentMethod", pu.PayMethod)
	if pu.IsMerchant != "" {
		params.Set("userType", "certified")
	} else {
		params.Set("userType", "all")
	}
	params.Set("sortType", "price_asc")
	params.Set("cryptoCurrency", strings.ToLower(currency))
	params.Set("fiatCurrency", strings.ToLower(fiat))
	params.Set("quoteMinAmountPerOrder", pu.Amount)
	params.Set("currentPage", "1") //;
	params.Set("numberPerPage", "10")

	resultokxbuy := OKXBuy{}

	url_buy := ("https://www.okx.com/v3/c2c/tradingOrders/getMarketplaceAdsPrelogin" + "?" + params.Encode())
	//log.Println("URL", url_buy)
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
		resultokxbuy, err = requestOrdersP2POKXBuy(url_buy)
		//log.Println("resulthOKX Buy", len(resultokxbuy.Data.Sell), "FIAT ", fiat, " Currency ", currency)
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
		//log.Println("Size for buy before return", len(resultokxbuy.Data.Sell), " for fiat", fiat, "and coin", currency)
		return resultokxbuy
	} else {
		return OKXBuy{}
	}
}

func requestOrdersP2POKXBuy(j string) (OKXBuy, error) {

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
		//log.Println("RESP BODY", resp.Body)

		// Make sure the response body is closed when we're done with it
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusTooManyRequests {
			return parsingJsonOKXBuy(resp.Body), nil
		}
		log.Println("Too many requests, backing off for", backoff)
		time.Sleep(backoff)
	}

}

func parsingJsonOKXBuy(r io.Reader) OKXBuy {
	var result OKXBuy

	body, _ := io.ReadAll(r)
	err := json.Unmarshal([]byte(body), &result)

	if err != nil {
		log.Println("Error unmarshal json URL Huobi:", err, string(body))
	}

	//log.Println("Size for buy", len(result.Data.Sell))
	//for i, j := range result.Data.Sell {
	//	log.Println(i, " - ", j.Price, j.NickName)
	//}

	return result
}
