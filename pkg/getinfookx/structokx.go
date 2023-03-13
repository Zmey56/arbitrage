package getinfookx

type PairOKX struct {
	Code string `json:"code"`
	Data []struct {
		CoinName string `json:"coinName"`
		InstID   string `json:"instId"`
		IsoFlag  string `json:"isoFlag"`
		Lever    string `json:"lever,omitempty"`
		ListTime string `json:"listTime"`
		NewLabel string `json:"newLabel"`
		NewTag   string `json:"newTag"`
	} `json:"data"`
	Msg string `json:"msg"`
}

type PaymentsOKX struct {
	Code int `json:"code"`
	Data []struct {
		FieldJSON []struct {
			EnableCopy        bool   `json:"enableCopy"`
			FieldFormat       int    `json:"fieldFormat"`
			FieldKey          string `json:"fieldKey"`
			FieldName         string `json:"fieldName"`
			FieldPlaceholder  string `json:"fieldPlaceholder"`
			FieldType         string `json:"fieldType"`
			FieldValidation   int    `json:"fieldValidation"`
			IsActive          int    `json:"isActive"`
			IsDesensitization int    `json:"isDesensitization"`
			IsRequired        int    `json:"isRequired"`
			Sort              int    `json:"sort"`
			ValidationFormat  string `json:"validationFormat"`
		} `json:"fieldJson"`
		MostUsed                 bool   `json:"mostUsed"`
		PaymentMethod            string `json:"paymentMethod"`
		PaymentMethodDescription string `json:"paymentMethodDescription"`
	} `json:"data"`
	DetailMsg    string `json:"detailMsg"`
	ErrorCode    string `json:"error_code"`
	ErrorMessage string `json:"error_message"`
	Msg          string `json:"msg"`
	RequestID    string `json:"requestId"`
}

type PaymentOKX struct {
	Identifier           string
	PayAccount           string
	PayMethodId          int
	PayType              string
	TradeMethodName      string
	TradeMethodShortName string
}

//type CoinOKX struct {
//	Code int `json:"code"`
//	Data []struct {
//		BaseColorfulIconURL string `json:"baseColorfulIconUrl"`
//		BaseCurrency        string `json:"baseCurrency"`
//		BaseCurrencyID      int    `json:"baseCurrencyId"`
//		BaseCurrencyName    string `json:"baseCurrencyName"`
//		BaseIconURL         string `json:"baseIconUrl"`
//		BaseName            string `json:"baseName"`
//		BaseSymbol          string `json:"baseSymbol"`
//	} `json:"data"`
//	DetailMsg    string `json:"detailMsg"`
//	ErrorCode    string `json:"error_code"`
//	ErrorMessage string `json:"error_message"`
//	Msg          string `json:"msg"`
//	RequestID    string `json:"requestId"`
//}

type CoinOKX struct {
	Code int `json:"code"`
	Data []struct {
		BaseColorfulIconURL string   `json:"baseColorfulIconUrl"`
		BaseCurrency        string   `json:"baseCurrency"`
		BaseCurrencyID      int      `json:"baseCurrencyId"`
		BaseCurrencyName    string   `json:"baseCurrencyName"`
		BaseIconURL         string   `json:"baseIconUrl"`
		BaseName            string   `json:"baseName"`
		BaseScale           int      `json:"baseScale"`
		BaseSymbol          string   `json:"baseSymbol"`
		BizTypes            []string `json:"bizTypes"`
		CanBeDisplayed      bool     `json:"canBeDisplayed"`
		CanBeSearched       bool     `json:"canBeSearched"`
		CountryIcon         string   `json:"countryIcon"`
		Gateway             struct {
			CanDeposit        bool `json:"canDeposit"`
			DepositTotalRange struct {
				Buy  string `json:"buy"`
				Sell string `json:"sell"`
			} `json:"depositTotalRange"`
		} `json:"gateway"`
		P2P struct {
			BaseDeposit               int  `json:"baseDeposit"`
			BlockTrade                bool `json:"blockTrade"`
			CanPlaceTrade             bool `json:"canPlaceTrade"`
			CanPlaceTradingOrder      bool `json:"canPlaceTradingOrder"`
			DefaultSelect             bool `json:"defaultSelect"`
			FlashTradeBetter          bool `json:"flashTradeBetter"`
			FlashTradeOrderTotalRange struct {
				Buy  string `json:"buy"`
				Sell string `json:"sell"`
			} `json:"flashTradeOrderTotalRange"`
			OrderTotalRange struct {
				Buy  string `json:"buy"`
				Sell string `json:"sell"`
			} `json:"orderTotalRange"`
			PlatformCommissionRate  int     `json:"platformCommissionRate"`
			PriceDeviationRateLimit float64 `json:"priceDeviationRateLimit"`
			QuotePriceIncrement     float64 `json:"quotePriceIncrement"`
			QuotePriceScale         int     `json:"quotePriceScale"`
			Swap                    bool    `json:"swap"`
			TradingOrderTotalRange  struct {
				Buy  string `json:"buy"`
				Sell string `json:"sell"`
			} `json:"tradingOrderTotalRange"`
			UsdStandard bool `json:"usdStandard"`
		} `json:"p2p"`
		QuoteCurrency       string  `json:"quoteCurrency"`
		QuoteName           string  `json:"quoteName"`
		QuotePriceIncrement float64 `json:"quotePriceIncrement"`
		QuotePriceScale     int     `json:"quotePriceScale"`
		QuoteScale          int     `json:"quoteScale"`
		QuoteSymbol         string  `json:"quoteSymbol"`
		Route               struct {
			Buy  string `json:"buy"`
			Sell string `json:"sell"`
		} `json:"route"`
	} `json:"data"`
	DetailMsg    string `json:"detailMsg"`
	ErrorCode    string `json:"error_code"`
	ErrorMessage string `json:"error_message"`
	Msg          string `json:"msg"`
	RequestID    string `json:"requestId"`
}

type ParametersOKX struct {
	CoinId       string
	Currency     string
	TradeType    string
	CurrPage     string
	PayMethod    string
	AcceptOrder  string
	Country      string
	BlockType    string
	Online       string
	Range        string
	Intrange     string
	Amount       string
	IsThumbsUp   string
	IsMerchant   string
	IsTraded     string
	OnlyTradable string
	IsFollowed   string
	PercentUser  float64
	Border       int
}
