package getinfohuobi

type PairHuobi struct {
	Status string `json:"status"`
	Data   []struct {
		BaseCurrency                    string  `json:"base-currency"`
		QuoteCurrency                   string  `json:"quote-currency"`
		PricePrecision                  int     `json:"price-precision"`
		AmountPrecision                 int     `json:"amount-precision"`
		SymbolPartition                 string  `json:"symbol-partition"`
		Symbol                          string  `json:"symbol"`
		State                           string  `json:"state"`
		ValuePrecision                  int     `json:"value-precision"`
		MinOrderAmt                     int     `json:"min-order-amt"`
		MaxOrderAmt                     int     `json:"max-order-amt"`
		MinOrderValue                   float64 `json:"min-order-value"`
		LimitOrderMinOrderAmt           int     `json:"limit-order-min-order-amt"`
		LimitOrderMaxOrderAmt           int     `json:"limit-order-max-order-amt"`
		LimitOrderMaxBuyAmt             int     `json:"limit-order-max-buy-amt"`
		LimitOrderMaxSellAmt            int     `json:"limit-order-max-sell-amt"`
		BuyLimitMustLessThan            float64 `json:"buy-limit-must-less-than"`
		SellLimitMustGreaterThan        float64 `json:"sell-limit-must-greater-than"`
		SellMarketMinOrderAmt           int     `json:"sell-market-min-order-amt"`
		SellMarketMaxOrderAmt           int     `json:"sell-market-max-order-amt"`
		BuyMarketMaxOrderValue          int     `json:"buy-market-max-order-value"`
		MarketSellOrderRateMustLessThan float64 `json:"market-sell-order-rate-must-less-than"`
		MarketBuyOrderRateMustLessThan  float64 `json:"market-buy-order-rate-must-less-than"`
		APITrading                      string  `json:"api-trading"`
		Tags                            string  `json:"tags,omitempty"`
		MaxOrderValue                   int     `json:"max-order-value,omitempty"`
		Underlying                      string  `json:"underlying,omitempty"`
		MgmtFeeRate                     float64 `json:"mgmt-fee-rate,omitempty"`
		ChargeTime                      string  `json:"charge-time,omitempty"`
		RebalTime                       string  `json:"rebal-time,omitempty"`
		RebalThreshold                  int     `json:"rebal-threshold,omitempty"`
		InitNav                         int     `json:"init-nav,omitempty"`
		LeverageRatio                   int     `json:"leverage-ratio,omitempty"`
		SuperMarginLeverageRatio        int     `json:"super-margin-leverage-ratio,omitempty"`
		FundingLeverageRatio            int     `json:"funding-leverage-ratio,omitempty"`
	} `json:"data"`
}

type CryptoFiat struct {
	Code int64 `json:"code"`
	Data []struct {
		CryptoAsset struct {
			AssetType interface{} `json:"assetType"`
			CoinID    int64       `json:"coinId"`
			Name      string      `json:"name"`
		} `json:"cryptoAsset"`
		QuoteAsset []struct {
			AssetType int64  `json:"assetType"`
			CoinID    int64  `json:"coinId"`
			Name      string `json:"name"`
		} `json:"quoteAsset"`
		Side int64 `json:"side"`
	} `json:"data"`
	Extend  interface{} `json:"extend"`
	Message string      `json:"message"`
	Success bool        `json:"success"`
}

type PaymentsHuobi struct {
	Success bool   `json:"success"`
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Country []struct {
			CountryID    int         `json:"countryId"`
			Code         int         `json:"code"`
			Name         string      `json:"name"`
			CurrencyID   int         `json:"currencyId"`
			CurrencyName string      `json:"currencyName"`
			IsForceAuth  int         `json:"isForceAuth"`
			AppShort     interface{} `json:"appShort"`
			NationalFlag string      `json:"nationalFlag"`
		} `json:"country"`
		MarketQuery []struct {
			CoinID   string `json:"coinId"`
			CoinName string `json:"coinName"`
			Query    struct {
				Type      string `json:"type"`
				Field     string `json:"field"`
				Name      string `json:"name"`
				Condition []struct {
					Key   string `json:"key"`
					Value string `json:"value"`
				} `json:"condition"`
			} `json:"query"`
		} `json:"marketQuery"`
		PayMethod []struct {
			PayMethodID  int         `json:"payMethodId"`
			Name         string      `json:"name"`
			DefaultName  interface{} `json:"defaultName"`
			Template     int         `json:"template"`
			BankImage    string      `json:"bankImage"`
			BankImageWeb interface{} `json:"bankImageWeb"`
			BankType     int         `json:"bankType"`
			Color        string      `json:"color"`
		} `json:"payMethod"`
		Currency []struct {
			CurrencyID         int         `json:"currencyId"`
			Name               string      `json:"name"`
			NameShort          string      `json:"nameShort"`
			Symbol             string      `json:"symbol"`
			PurchaseSellFlag   bool        `json:"purchaseSellFlag"`
			PurchaseBuyFlag    bool        `json:"purchaseBuyFlag"`
			AppShowP2P         interface{} `json:"appShowP2p"`
			ShowPtoP           interface{} `json:"showPtoP"`
			SupportPayments    []int       `json:"supportPayments"`
			ShowOrderPremature int         `json:"showOrderPremature"`
			WhetherGib         int         `json:"whetherGib"`
			Entrance           int         `json:"entrance"`
		} `json:"currency"`
		Coin []struct {
			CoinID             int         `json:"coinId"`
			CoinCode           string      `json:"coinCode"`
			ShortName          string      `json:"shortName"`
			Name               string      `json:"name"`
			IsTransfer         bool        `json:"isTransfer"`
			IsRecharge         bool        `json:"isRecharge"`
			IsTrade            bool        `json:"isTrade"`
			AppLogo            string      `json:"appLogo"`
			WebLogo            interface{} `json:"webLogo"`
			Confirmations      int         `json:"confirmations"`
			MinRecharge        int         `json:"minRecharge"`
			IsRemote           interface{} `json:"isRemote"`
			IsShow             bool        `json:"isShow"`
			IsLiteShow         bool        `json:"isLiteShow"`
			LiteLogo           interface{} `json:"liteLogo"`
			IsLiteCurrencySell bool        `json:"isLiteCurrencySell"`
			IsLiteCurrencyBuy  bool        `json:"isLiteCurrencyBuy"`
			IsLiteCoinSell     bool        `json:"isLiteCoinSell"`
			IsLiteCoinBuy      bool        `json:"isLiteCoinBuy"`
			TradePrecision     string      `json:"tradePrecision"`
			ShowPrecision      string      `json:"showPrecision"`
			FullName           string      `json:"fullName"`
			CoinType           int         `json:"coinType"`
			EnFullName         string      `json:"enFullName"`
			IsShowAssetDetail  bool        `json:"isShowAssetDetail"`
		} `json:"coin"`
	} `json:"data"`
	Extend interface{} `json:"extend"`
}

type PaymentHuobi struct {
	Identifier           string
	PayAccount           string
	PayMethodId          int
	PayType              string
	TradeMethodName      string
	TradeMethodShortName string
}

type ParametersHuobi struct {
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
}
