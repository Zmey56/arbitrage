// Function for parsing data from Binance about orders and get number pages

package getinfobinance

import (
	"encoding/json"
	"io"
	"strconv"
)

type Adv struct {
	AdvNo                           string
	Classify                        string
	TradeType                       string
	Asset                           string
	FiatUnit                        string
	AdvStatus                       string
	PriceType                       string
	PriceFloatingRatio              string
	RateFloatingRatio               string
	CurrencyRate                    string
	Price                           string
	InitAmount                      string
	SurplusAmount                   string
	AmountAfterEditing              string
	MaxSingleTransAmount            string
	MinSingleTransAmount            string
	BuyerKycLimit                   string
	BuyerRegDaysLimit               string
	BuyerBtcPositionLimit           string
	Remarks                         string
	AutoReplyMsg                    string
	PayTimeLimit                    string
	TradeMethods                    TradeMethods
	UserTradeCountFilterTime        string
	UserBuyTradeCountMin            string
	UserBuyTradeCountMax            string
	UserSellTradeCountMin           string
	UserSellTradeCountMax           string
	UserAllTradeCountMin            string
	UserAllTradeCountMax            string
	UserTradeCompleteRateFilterTime string
	UserTradeCompleteCountMin       string
	UserTradeCompleteRateMin        string
	UserTradeVolumeFilterTime       string
	UserTradeType                   string
	UserTradeVolumeMin              string
	UserTradeVolumeMax              string
	UserTradeVolumeAsset            string
	CreateTime                      string
	AdvUpdateTime                   string
	FiatVo                          string
	AssetVo                         string
	AdvVisibleRet                   string
	AssetLogo                       string
	AssetScale                      string
	FiatScale                       string
	PriceScale                      string
	FiatSymbol                      string
	IsTradable                      string
	DynamicMaxSingleTransAmount     string
	MinSingleTransQuantity          string
	MaxSingleTransQuantity          string
	DynamicMaxSingleTransQuantity   string
	TradableQuantity                string
	CommissionRate                  string
	TradeMethodCommissionRates      string
	LaunchCountry                   string
	AbnormalStatusList              string
	CloseReason                     string
}

type TradeMethods struct {
	PayId                string
	PayMethodId          string
	PayType              string
	PayAccount           string
	PayBank              string
	PaySubBank           string
	Identifier           string
	IconUrlColor         string
	TradeMethodName      string
	TradeMethodShortName string
	TradeMethodBgColor   string
}

type Advertiser struct {
	UserNo           string
	RealName         string
	NickName         string
	Margin           string
	MarginUnit       string
	OrderCount       string
	MonthOrderCount  float64
	MonthFinishRate  float64
	AdvConfirmTime   string
	Email            string
	RegistrationTime string
	Mobile           string
	UserType         string
	TagIconUrls      []string
	UserGrade        string
	UserIdentity     string
	ProMerchant      string
	IsBlocked        string
}

// All orders and Advertiser
// type AllAdv []Adv
type AllAdvertiser []Advertiser

type AdvertiserAdv struct {
	Advs        Adv
	Advertisers Advertiser
}

func ParsingJson(r io.Reader, tradeType string, transAmount float64, finish bool) (AdvertiserAdv, int, bool) {
	//alladv := Adv{}
	alladvertiser := AllAdvertiser{}
	advertiseradv := AdvertiserAdv{}
	//var transAmountTrue []bool
	//foundenough := false //this var will be true than found user with enougth money
	count := 0
	indexresult := 0

	var result map[string]any

	body, _ := io.ReadAll(r)
	json.Unmarshal([]byte(body), &result)

	//numbers of rows
	numberRows := int(result["total"].(float64))

	for _, value := range result["data"].([]interface{}) {
		for j, v := range value.(map[string]any) {
			adv := Adv{}
			advertiser := Advertiser{}
			if j == "adv" {
				count++
				for q, s := range v.(map[string]any) {
					switch q {
					case "advNo":
						adv.AdvNo = s.(string)
					case "classify":
						adv.Classify = s.(string)
					case "price":
						adv.Price = s.(string)
					case "surplusAmount":
						//finish and need sell
						adv.SurplusAmount = s.(string)
					case "tradableQuantity":
						adv.TradableQuantity = s.(string)
					case "maxSingleTransAmount":
						adv.MaxSingleTransAmount = s.(string)
					case "minSingleTransAmount":
						adv.MinSingleTransAmount = s.(string)
					case "maxSingleTransQuantity":
						adv.MinSingleTransQuantity = s.(string)
					case "minSingleTransQuantity":
						adv.MinSingleTransQuantity = s.(string)
					case "fiatUnit":
						adv.FiatUnit = s.(string)
					case "asset":
						adv.Asset = s.(string)
					case "commissionRate":
						adv.CommissionRate = s.(string)
					case "tradeMethods":
						for _, d := range s.([]interface{}) {
							for m, t := range d.(map[string]any) {
								switch m {
								case "identifier":
									adv.TradeMethods.Identifier = t.(string)
								case "tradeMethodName":
									if t != nil {
										adv.TradeMethods.TradeMethodName = t.(string)
									}
								case "tradeMethodShortName":
									if t != nil {
										adv.TradeMethods.TradeMethodShortName = t.(string)
									}
								default:
									continue
								}
							}
						}
					default:
						continue
					}
				}
				if tradeType == "Buy" {
					maxSingleTransAmount, _ := strconv.ParseFloat(adv.MaxSingleTransAmount, 64)
					minSingleTransAmount := 0.0
					if transAmount > 0 {
						minSingleTransAmount, _ = strconv.ParseFloat(adv.MinSingleTransAmount, 64)
					}
					if transAmount >= minSingleTransAmount && transAmount <= maxSingleTransAmount && !finish {
						finish = true
						indexresult = count
						advertiseradv.Advs = adv
					} else {
						continue
					}
				} else if tradeType == "Sell" {
					if transAmount > 0 {
						surplusAmount, _ := strconv.ParseFloat(adv.SurplusAmount, 64)
						if transAmount <= surplusAmount && !finish {
							finish = true
							indexresult = count
							advertiseradv.Advs = adv
						}
					}
				}
			} else if j == "advertiser" {

				for m, k := range v.(map[string]any) {
					switch m {
					case "nickName": //+
						advertiser.NickName = k.(string)
					case "monthOrderCount":
						advertiser.MonthOrderCount = k.(float64)
					case "monthFinishRate":
						advertiser.MonthFinishRate = k.(float64)
					case "userNo":
						advertiser.UserNo = k.(string)
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

	if indexresult != 0 {
		advertiseradv.Advertisers = alladvertiser[indexresult-1]
	}

	return advertiseradv, numberPages, finish
}
