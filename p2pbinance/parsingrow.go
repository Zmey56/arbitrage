// Function for parsing data from Binance and get number pages

package p2pbinance

import (
	"encoding/json"
	"fmt"
	"io"
)

type Trader struct {
	advNo                           string
	classify                        string
	tradeType                       string
	asset                           string
	fiatUnit                        string
	advStatus                       string
	priceType                       string
	priceFloatingRatio              string
	rateFloatingRatio               string
	currencyRate                    string
	price                           string
	initAmount                      string
	surplusAmount                   string
	amountAfterEditing              string
	maxSingleTransAmount            string
	minSingleTransAmount            string
	buyerKycLimit                   string
	buyerRegDaysLimit               string
	buyerBtcPositionLimit           string
	remarks                         string
	autoReplyMsg                    string
	payTimeLimit                    string
	tradeMethods                    TradeMethods
	userTradeCountFilterTime        string
	userBuyTradeCountMin            string
	userBuyTradeCountMax            string
	userSellTradeCountMin           string
	userSellTradeCountMax           string
	userAllTradeCountMin            string
	userAllTradeCountMax            string
	userTradeCompleteRateFilterTime string
	userTradeCompleteCountMin       string
	userTradeCompleteRateMin        string
	userTradeVolumeFilterTime       string
	userTradeType                   string
	userTradeVolumeMin              string
	userTradeVolumeMax              string
	userTradeVolumeAsset            string
	createTime                      string
	advUpdateTime                   string
	fiatVo                          string
	assetVo                         string
	advVisibleRet                   string
	assetLogo                       string
	assetScale                      string
	fiatScale                       string
	priceScale                      string
	fiatSymbol                      string
	isTradable                      string
	dynamicMaxSingleTransAmount     string
	minSingleTransQuantity          string
	maxSingleTransQuantity          string
	dynamicMaxSingleTransQuantity   string
	tradableQuantity                string
	commissionRate                  string
	tradeMethodCommissionRates      string
	launchCountry                   string
	abnormalStatusList              string
	closeReason                     string
}

type TradeMethods struct {
	payId                string
	payMethodId          string
	payType              string
	payAccount           string
	payBank              string
	paySubBank           string
	identifier           string
	iconUrlColor         string
	tradeMethodName      string
	tradeMethodShortName string
	tradeMethodBgColor   string
}

type Advertiser struct {
	userNo           string
	realName         string
	nickName         string
	margin           string
	marginUnit       string
	orderCount       string
	monthOrderCount  float64
	monthFinishRate  float64
	advConfirmTime   string
	email            string
	registrationTime string
	mobile           string
	userType         string
	tagIconUrls      []string
	userGrade        string
	userIdentity     string
	proMerchant      string
	isBlocked        string
}

// All orders and Advertiser
type AllAdv []Trader
type AllAdvertiser []Advertiser

type AdvertiserAdv struct {
	Advs        AllAdv
	Advertisers AllAdvertiser
}

func ParsingJson(r io.Reader) (AdvertiserAdv, int) {
	alladv := AllAdv{}
	alladvertiser := AllAdvertiser{}
	advertiseradv := AdvertiserAdv{}

	var result map[string]any

	body, _ := io.ReadAll(r)
	json.Unmarshal([]byte(body), &result)

	//numbers of rows
	numberRows := int(result["total"].(float64))

	for _, value := range result["data"].([]interface{}) {

		for j, v := range value.(map[string]any) {
			adv := Trader{}
			advertiser := Advertiser{}
			if j == "adv" {
				fmt.Println("adv")
				for q, s := range v.(map[string]any) {
					switch q {
					case "advNo":
						adv.advNo = s.(string)
						fmt.Println("advNo - ", s.(string))
					case "classify":
						adv.classify = s.(string)
					case "price":
						adv.price = s.(string)
					case "tradableQuantity":
						adv.tradableQuantity = s.(string)
					case "maxSingleTransAmount":
						adv.maxSingleTransAmount = s.(string)
					case "maxSingleTransQuantity":
						adv.minSingleTransQuantity = s.(string)
					case "minSingleTransQuantity":
						adv.minSingleTransQuantity = s.(string)
					case "fiatUnit":
						adv.fiatUnit = s.(string)
					case "asset":
						adv.asset = s.(string)
					case "tradeMethods":
						for _, d := range s.([]interface{}) {
							for m, t := range d.(map[string]any) {
								switch m {
								case "identifier":
									adv.tradeMethods.identifier = t.(string)
								case "tradeMethodName":
									adv.tradeMethods.tradeMethodName = t.(string)
								//case "tradeMethodShortName":
								//	adv.tradeMethods.tradeMethodShortName = t.(string)
								default:
									continue
								}
							}
						}
					default:
						continue
					}
				}
				alladv = append(alladv, adv)
			} else if j == "advertiser" {
				fmt.Println("advertiser")

				for m, k := range v.(map[string]any) {
					switch m {
					case "nickName": //+
						advertiser.nickName = k.(string)
						fmt.Println("nickName - ", k.(string))
					case "monthOrderCount":
						advertiser.monthOrderCount = k.(float64)
					case "monthFinishRate":
						advertiser.monthFinishRate = k.(float64)
					case "userNo":
						advertiser.userNo = k.(string)
					default:
						continue
					}
				}
				alladvertiser = append(alladvertiser, advertiser)
			} else {
				continue
			}
		}
	}
	numberPages := numberRows / 10

	advertiseradv.Advertisers = alladvertiser
	advertiseradv.Advs = alladv

	return advertiseradv, numberPages
}

func NextPages() {

}

//if namberPages > 1 {
//	//fmt.Println(string(jsonData))
//	for i := 2; i <= namberPages; i++ {
//		var m map[string]interface{}
//		err := json.Unmarshal(jsonData, &m)
//		if err != nil {
//			log.Println("Error", err)
//		}
//		m["page"] = i
//		newData, err := json.Marshal(m)
//		newJsonData := []byte(string(newData))
//		request, error := http.NewRequest("POST", httpposturl, bytes.NewBuffer(newJsonData))
//		request.Header.Set("Content-Type", "application/json; charset=UTF-8")
//
//		client := &http.Client{}
//		response, error := client.Do(request)
//		if error != nil {
//			panic(error)
//		}
//		defer response.Body.Close()
//
//		fmt.Println("response Status:", response.Status)
//		fmt.Println("response Headers:", response.Header)
//		body, _ := io.ReadAll(response.Body)
//		fmt.Println("response Body:", string(body))
//		if error != nil {
//			log.Println("Error", err)
//		}
//	}
//}
