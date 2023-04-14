// Function for parsing data from Binance about orders and get number pages

package getinfobinance

type AdvertiserAdv struct {
	Code          string      `json:"code"`
	Message       interface{} `json:"message"`
	MessageDetail interface{} `json:"messageDetail"`
	Data          []struct {
		Adv struct {
			AdvNo                 string      `json:"advNo"`
			Classify              string      `json:"classify"`
			TradeType             string      `json:"tradeType"`
			Asset                 string      `json:"asset"`
			FiatUnit              string      `json:"fiatUnit"`
			AdvStatus             interface{} `json:"advStatus"`
			PriceType             interface{} `json:"priceType"`
			PriceFloatingRatio    interface{} `json:"priceFloatingRatio"`
			RateFloatingRatio     interface{} `json:"rateFloatingRatio"`
			CurrencyRate          interface{} `json:"currencyRate"`
			Price                 float64     `json:",string"`
			InitAmount            interface{} `json:"initAmount"`
			SurplusAmount         string      `json:"surplusAmount"`
			AmountAfterEditing    interface{} `json:"amountAfterEditing"`
			MaxSingleTransAmount  string      `json:"maxSingleTransAmount"`
			MinSingleTransAmount  string      `json:"minSingleTransAmount"`
			BuyerKycLimit         interface{} `json:"buyerKycLimit"`
			BuyerRegDaysLimit     interface{} `json:"buyerRegDaysLimit"`
			BuyerBtcPositionLimit interface{} `json:"buyerBtcPositionLimit"`
			Remarks               interface{} `json:"remarks"`
			AutoReplyMsg          string      `json:"autoReplyMsg"`
			PayTimeLimit          interface{} `json:"payTimeLimit"`
			TradeMethods          []struct {
				PayID                interface{} `json:"payId"`
				PayMethodID          string      `json:"payMethodId"`
				PayType              interface{} `json:"payType"`
				PayAccount           interface{} `json:"payAccount"`
				PayBank              interface{} `json:"payBank"`
				PaySubBank           interface{} `json:"paySubBank"`
				Identifier           string      `json:"identifier"`
				IconURLColor         interface{} `json:"iconUrlColor"`
				TradeMethodName      string      `json:"tradeMethodName"`
				TradeMethodShortName interface{} `json:"tradeMethodShortName"`
				TradeMethodBgColor   string      `json:"tradeMethodBgColor"`
			} `json:"tradeMethods"`
			UserTradeCountFilterTime        interface{}   `json:"userTradeCountFilterTime"`
			UserBuyTradeCountMin            interface{}   `json:"userBuyTradeCountMin"`
			UserBuyTradeCountMax            interface{}   `json:"userBuyTradeCountMax"`
			UserSellTradeCountMin           interface{}   `json:"userSellTradeCountMin"`
			UserSellTradeCountMax           interface{}   `json:"userSellTradeCountMax"`
			UserAllTradeCountMin            interface{}   `json:"userAllTradeCountMin"`
			UserAllTradeCountMax            interface{}   `json:"userAllTradeCountMax"`
			UserTradeCompleteRateFilterTime interface{}   `json:"userTradeCompleteRateFilterTime"`
			UserTradeCompleteCountMin       interface{}   `json:"userTradeCompleteCountMin"`
			UserTradeCompleteRateMin        interface{}   `json:"userTradeCompleteRateMin"`
			UserTradeVolumeFilterTime       interface{}   `json:"userTradeVolumeFilterTime"`
			UserTradeType                   interface{}   `json:"userTradeType"`
			UserTradeVolumeMin              interface{}   `json:"userTradeVolumeMin"`
			UserTradeVolumeMax              interface{}   `json:"userTradeVolumeMax"`
			UserTradeVolumeAsset            interface{}   `json:"userTradeVolumeAsset"`
			CreateTime                      interface{}   `json:"createTime"`
			AdvUpdateTime                   interface{}   `json:"advUpdateTime"`
			FiatVo                          interface{}   `json:"fiatVo"`
			AssetVo                         interface{}   `json:"assetVo"`
			AdvVisibleRet                   interface{}   `json:"advVisibleRet"`
			AssetLogo                       interface{}   `json:"assetLogo"`
			AssetScale                      int           `json:"assetScale"`
			FiatScale                       int           `json:"fiatScale"`
			PriceScale                      int           `json:"priceScale"`
			FiatSymbol                      string        `json:"fiatSymbol"`
			IsTradable                      bool          `json:"isTradable"`
			DynamicMaxSingleTransAmount     string        `json:"dynamicMaxSingleTransAmount"`
			MinSingleTransQuantity          string        `json:"minSingleTransQuantity"`
			MaxSingleTransQuantity          string        `json:"maxSingleTransQuantity"`
			DynamicMaxSingleTransQuantity   string        `json:"dynamicMaxSingleTransQuantity"`
			TradableQuantity                string        `json:"tradableQuantity"`
			CommissionRate                  string        `json:"commissionRate"`
			TradeMethodCommissionRates      []interface{} `json:"tradeMethodCommissionRates"`
			LaunchCountry                   interface{}   `json:"launchCountry"`
			AbnormalStatusList              interface{}   `json:"abnormalStatusList"`
			CloseReason                     interface{}   `json:"closeReason"`
		} `json:"adv"`
		Advertiser struct {
			UserNo           string        `json:"userNo"`
			RealName         interface{}   `json:"realName"`
			NickName         string        `json:"nickName"`
			Margin           interface{}   `json:"margin"`
			MarginUnit       interface{}   `json:"marginUnit"`
			OrderCount       interface{}   `json:"orderCount"`
			MonthOrderCount  int           `json:"monthOrderCount"`
			MonthFinishRate  float64       `json:"monthFinishRate"`
			AdvConfirmTime   interface{}   `json:"advConfirmTime"`
			Email            interface{}   `json:"email"`
			RegistrationTime interface{}   `json:"registrationTime"`
			Mobile           interface{}   `json:"mobile"`
			UserType         string        `json:"userType"`
			TagIconUrls      []interface{} `json:"tagIconUrls"`
			UserGrade        int           `json:"userGrade"`
			UserIdentity     string        `json:"userIdentity"`
			ProMerchant      interface{}   `json:"proMerchant"`
			IsBlocked        interface{}   `json:"isBlocked"`
		} `json:"advertiser"`
	} `json:"data"`
	Total   int  `json:"total"`
	Success bool `json:"success"`
}

type Binance struct {
	Code          string      `json:"code"`
	Message       interface{} `json:"message"`
	MessageDetail interface{} `json:"messageDetail"`
	Data          []struct {
		Adv struct {
			AdvNo                 string      `json:"advNo"`
			Classify              string      `json:"classify"`
			TradeType             string      `json:"tradeType"`
			Asset                 string      `json:"asset"`
			FiatUnit              string      `json:"fiatUnit"`
			AdvStatus             interface{} `json:"advStatus"`
			PriceType             interface{} `json:"priceType"`
			PriceFloatingRatio    interface{} `json:"priceFloatingRatio"`
			RateFloatingRatio     interface{} `json:"rateFloatingRatio"`
			CurrencyRate          interface{} `json:"currencyRate"`
			Price                 float64     `json:",string"`
			InitAmount            interface{} `json:"initAmount"`
			SurplusAmount         string      `json:"surplusAmount"`
			AmountAfterEditing    interface{} `json:"amountAfterEditing"`
			MaxSingleTransAmount  string      `json:"maxSingleTransAmount"`
			MinSingleTransAmount  string      `json:"minSingleTransAmount"`
			BuyerKycLimit         interface{} `json:"buyerKycLimit"`
			BuyerRegDaysLimit     interface{} `json:"buyerRegDaysLimit"`
			BuyerBtcPositionLimit interface{} `json:"buyerBtcPositionLimit"`
			Remarks               interface{} `json:"remarks"`
			AutoReplyMsg          string      `json:"autoReplyMsg"`
			PayTimeLimit          interface{} `json:"payTimeLimit"`
			TradeMethods          []struct {
				PayID                interface{} `json:"payId"`
				PayMethodID          string      `json:"payMethodId"`
				PayType              interface{} `json:"payType"`
				PayAccount           interface{} `json:"payAccount"`
				PayBank              interface{} `json:"payBank"`
				PaySubBank           interface{} `json:"paySubBank"`
				Identifier           string      `json:"identifier"`
				IconURLColor         interface{} `json:"iconUrlColor"`
				TradeMethodName      string      `json:"tradeMethodName"`
				TradeMethodShortName interface{} `json:"tradeMethodShortName"`
				TradeMethodBgColor   string      `json:"tradeMethodBgColor"`
			} `json:"tradeMethods"`
			UserTradeCountFilterTime        interface{}   `json:"userTradeCountFilterTime"`
			UserBuyTradeCountMin            interface{}   `json:"userBuyTradeCountMin"`
			UserBuyTradeCountMax            interface{}   `json:"userBuyTradeCountMax"`
			UserSellTradeCountMin           interface{}   `json:"userSellTradeCountMin"`
			UserSellTradeCountMax           interface{}   `json:"userSellTradeCountMax"`
			UserAllTradeCountMin            interface{}   `json:"userAllTradeCountMin"`
			UserAllTradeCountMax            interface{}   `json:"userAllTradeCountMax"`
			UserTradeCompleteRateFilterTime interface{}   `json:"userTradeCompleteRateFilterTime"`
			UserTradeCompleteCountMin       interface{}   `json:"userTradeCompleteCountMin"`
			UserTradeCompleteRateMin        interface{}   `json:"userTradeCompleteRateMin"`
			UserTradeVolumeFilterTime       interface{}   `json:"userTradeVolumeFilterTime"`
			UserTradeType                   interface{}   `json:"userTradeType"`
			UserTradeVolumeMin              interface{}   `json:"userTradeVolumeMin"`
			UserTradeVolumeMax              interface{}   `json:"userTradeVolumeMax"`
			UserTradeVolumeAsset            interface{}   `json:"userTradeVolumeAsset"`
			CreateTime                      interface{}   `json:"createTime"`
			AdvUpdateTime                   interface{}   `json:"advUpdateTime"`
			FiatVo                          interface{}   `json:"fiatVo"`
			AssetVo                         interface{}   `json:"assetVo"`
			AdvVisibleRet                   interface{}   `json:"advVisibleRet"`
			AssetLogo                       interface{}   `json:"assetLogo"`
			AssetScale                      int           `json:"assetScale"`
			FiatScale                       int           `json:"fiatScale"`
			PriceScale                      int           `json:"priceScale"`
			FiatSymbol                      string        `json:"fiatSymbol"`
			IsTradable                      bool          `json:"isTradable"`
			DynamicMaxSingleTransAmount     string        `json:"dynamicMaxSingleTransAmount"`
			MinSingleTransQuantity          string        `json:"minSingleTransQuantity"`
			MaxSingleTransQuantity          string        `json:"maxSingleTransQuantity"`
			DynamicMaxSingleTransQuantity   string        `json:"dynamicMaxSingleTransQuantity"`
			TradableQuantity                string        `json:"tradableQuantity"`
			CommissionRate                  string        `json:"commissionRate"`
			TradeMethodCommissionRates      []interface{} `json:"tradeMethodCommissionRates"`
			LaunchCountry                   interface{}   `json:"launchCountry"`
			AbnormalStatusList              interface{}   `json:"abnormalStatusList"`
			CloseReason                     interface{}   `json:"closeReason"`
		} `json:"adv"`
		Advertiser struct {
			UserNo           string        `json:"userNo"`
			RealName         interface{}   `json:"realName"`
			NickName         string        `json:"nickName"`
			Margin           interface{}   `json:"margin"`
			MarginUnit       interface{}   `json:"marginUnit"`
			OrderCount       interface{}   `json:"orderCount"`
			MonthOrderCount  int           `json:"monthOrderCount"`
			MonthFinishRate  float64       `json:"monthFinishRate"`
			AdvConfirmTime   interface{}   `json:"advConfirmTime"`
			Email            interface{}   `json:"email"`
			RegistrationTime interface{}   `json:"registrationTime"`
			Mobile           interface{}   `json:"mobile"`
			UserType         string        `json:"userType"`
			TagIconUrls      []interface{} `json:"tagIconUrls"`
			UserGrade        int           `json:"userGrade"`
			UserIdentity     string        `json:"userIdentity"`
			ProMerchant      interface{}   `json:"proMerchant"`
			IsBlocked        interface{}   `json:"isBlocked"`
		} `json:"advertiser"`
	} `json:"data"`
	Total   int  `json:"total"`
	Success bool `json:"success"`
}

type PaymentBinance struct {
	Identifier           string
	PayAccount           string
	PayMethodId          int
	PayType              string
	TradeMethodName      string
	TradeMethodShortName string
}
