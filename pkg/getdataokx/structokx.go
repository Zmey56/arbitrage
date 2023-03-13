package getdataokx

type OKXBuy struct {
	Code int `json:"code"`
	Data struct {
		Buy  []interface{} `json:"buy"`
		Sell []struct {
			AlreadyTraded             bool        `json:"alreadyTraded"`
			AvailableAmount           string      `json:"availableAmount"`
			BaseCurrency              string      `json:"baseCurrency"`
			Black                     bool        `json:"black"`
			CancelledOrderQuantity    int         `json:"cancelledOrderQuantity"`
			CompletedOrderQuantity    int         `json:"completedOrderQuantity"`
			CompletedRate             string      `json:"completedRate"`
			CreatorType               string      `json:"creatorType"`
			GuideUpgradeKyc           bool        `json:"guideUpgradeKyc"`
			ID                        string      `json:"id"`
			Intention                 bool        `json:"intention"`
			MaxCompletedOrderQuantity int         `json:"maxCompletedOrderQuantity"`
			MaxUserCreatedDate        int         `json:"maxUserCreatedDate"`
			MerchantID                string      `json:"merchantId"`
			MinCompletedOrderQuantity int         `json:"minCompletedOrderQuantity"`
			MinCompletionRate         string      `json:"minCompletionRate"`
			MinKycLevel               int         `json:"minKycLevel"`
			MinSellOrders             int         `json:"minSellOrders"`
			Mine                      bool        `json:"mine"`
			NickName                  string      `json:"nickName"`
			PaymentMethods            []string    `json:"paymentMethods"`
			Price                     string      `json:"price"`
			PublicUserID              string      `json:"publicUserId"`
			QuoteCurrency             string      `json:"quoteCurrency"`
			QuoteMaxAmountPerOrder    string      `json:"quoteMaxAmountPerOrder"`
			QuoteMinAmountPerOrder    string      `json:"quoteMinAmountPerOrder"`
			QuoteScale                int         `json:"quoteScale"`
			QuoteSymbol               string      `json:"quoteSymbol"`
			ReceivingAds              bool        `json:"receivingAds"`
			SafetyLimit               bool        `json:"safetyLimit"`
			Side                      string      `json:"side"`
			UserActiveStatusVo        interface{} `json:"userActiveStatusVo"`
			UserType                  string      `json:"userType"`
			VerificationType          int         `json:"verificationType"`
		} `json:"sell"`
		Total int `json:"total"`
	} `json:"data"`
	DetailMsg    string `json:"detailMsg"`
	ErrorCode    string `json:"error_code"`
	ErrorMessage string `json:"error_message"`
	Msg          string `json:"msg"`
	RequestID    string `json:"requestId"`
}

type OKXSell struct {
	Code int `json:"code"`
	Data struct {
		Buy []struct {
			AlreadyTraded             bool        `json:"alreadyTraded"`
			AvailableAmount           string      `json:"availableAmount"`
			BaseCurrency              string      `json:"baseCurrency"`
			Black                     bool        `json:"black"`
			CancelledOrderQuantity    int         `json:"cancelledOrderQuantity"`
			CompletedOrderQuantity    int         `json:"completedOrderQuantity"`
			CompletedRate             string      `json:"completedRate"`
			CreatorType               string      `json:"creatorType"`
			GuideUpgradeKyc           bool        `json:"guideUpgradeKyc"`
			ID                        string      `json:"id"`
			Intention                 bool        `json:"intention"`
			MaxCompletedOrderQuantity int         `json:"maxCompletedOrderQuantity"`
			MaxUserCreatedDate        int         `json:"maxUserCreatedDate"`
			MerchantID                string      `json:"merchantId"`
			MinCompletedOrderQuantity int         `json:"minCompletedOrderQuantity"`
			MinCompletionRate         string      `json:"minCompletionRate"`
			MinKycLevel               int         `json:"minKycLevel"`
			MinSellOrders             int         `json:"minSellOrders"`
			Mine                      bool        `json:"mine"`
			NickName                  string      `json:"nickName"`
			PaymentMethods            []string    `json:"paymentMethods"`
			Price                     string      `json:"price"`
			PublicUserID              string      `json:"publicUserId"`
			QuoteCurrency             string      `json:"quoteCurrency"`
			QuoteMaxAmountPerOrder    string      `json:"quoteMaxAmountPerOrder"`
			QuoteMinAmountPerOrder    string      `json:"quoteMinAmountPerOrder"`
			QuoteScale                int         `json:"quoteScale"`
			QuoteSymbol               string      `json:"quoteSymbol"`
			ReceivingAds              bool        `json:"receivingAds"`
			SafetyLimit               bool        `json:"safetyLimit"`
			Side                      string      `json:"side"`
			UserActiveStatusVo        interface{} `json:"userActiveStatusVo"`
			UserType                  string      `json:"userType"`
			VerificationType          int         `json:"verificationType"`
		} `json:"buy"`
		Sell  []interface{} `json:"sell"`
		Total int           `json:"total"`
	} `json:"data"`
	DetailMsg    string `json:"detailMsg"`
	ErrorCode    string `json:"error_code"`
	ErrorMessage string `json:"error_message"`
	Msg          string `json:"msg"`
	RequestID    string `json:"requestId"`
}

type RatePairOKX struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
	Data []struct {
		InstType  string `json:"instType"`
		InstID    string `json:"instId"`
		Last      string `json:"last"`
		LastSz    string `json:"lastSz"`
		AskPx     string `json:"askPx"`
		AskSz     string `json:"askSz"`
		BidPx     string `json:"bidPx"`
		BidSz     string `json:"bidSz"`
		Open24H   string `json:"open24h"`
		High24H   string `json:"high24h"`
		Low24H    string `json:"low24h"`
		VolCcy24H string `json:"volCcy24h"`
		Vol24H    string `json:"vol24h"`
		Ts        string `json:"ts"`
		SodUtc0   string `json:"sodUtc0"`
		SodUtc8   string `json:"sodUtc8"`
	} `json:"data"`
}
