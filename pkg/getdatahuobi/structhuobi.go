package getdatahuobi

type Huobi struct {
	TimeData   int64  `json:"time_data"`
	Code       int    `json:"code"`
	Message    string `json:"message"`
	TotalCount int    `json:"totalCount"`
	PageSize   int    `json:"pageSize"`
	TotalPage  int    `json:"totalPage"`
	CurrPage   int    `json:"currPage"`
	Data       []struct {
		ID            int         `json:"id"`
		UID           int         `json:"uid"`
		UserName      string      `json:"userName"`
		MerchantLevel int         `json:"merchantLevel"`
		MerchantTags  interface{} `json:"merchantTags"`
		CoinID        int         `json:"coinId"`
		Currency      int         `json:"currency"`
		TradeType     int         `json:"tradeType"`
		BlockType     int         `json:"blockType"`
		PayMethod     string      `json:"payMethod"`
		PayMethods    []struct {
			PayMethodID int         `json:"payMethodId"`
			Name        string      `json:"name"`
			Color       string      `json:"color"`
			IsRecommend interface{} `json:"isRecommend"`
		} `json:"payMethods"`
		PayTerm           int         `json:"payTerm"`
		PayName           string      `json:"payName"`
		MinTradeLimit     string      `json:"minTradeLimit"`
		MaxTradeLimit     string      `json:"maxTradeLimit"`
		Price             string      `json:"price"`
		TradeCount        string      `json:"tradeCount"`
		IsOnline          bool        `json:"isOnline"`
		IsFollowed        bool        `json:"isFollowed"`
		TradeMonthTimes   int         `json:"tradeMonthTimes"`
		OrderCompleteRate string      `json:"orderCompleteRate"`
		TakerAcceptOrder  int         `json:"takerAcceptOrder"`
		TakerAcceptAmount string      `json:"takerAcceptAmount"`
		TakerLimit        int         `json:"takerLimit"`
		GmtSort           int64       `json:"gmtSort"`
		IsCopyBlock       bool        `json:"isCopyBlock"`
		ThumbUp           int         `json:"thumbUp"`
		SeaViewRoom       interface{} `json:"seaViewRoom"`
	} `json:"data"`
	Success bool `json:"success"`
}

type CryptoHuobi struct {
	Success bool   `json:"success"`
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    []struct {
		Side       int `json:"side"`
		QuoteAsset []struct {
			Name      string `json:"name"`
			CoinID    int    `json:"coinId"`
			AssetType int    `json:"assetType"`
		} `json:"quoteAsset"`
		CryptoAsset struct {
			Name      string      `json:"name"`
			CoinID    int         `json:"coinId"`
			AssetType interface{} `json:"assetType"`
		} `json:"cryptoAsset"`
	} `json:"data"`
	Extend interface{} `json:"extend"`
}

type PayMethodHuobi struct {
	PayMethodID  int         `json:"payMethodId"`
	Name         string      `json:"name"`
	DefaultName  interface{} `json:"defaultName"`
	Template     int         `json:"template"`
	BankImage    string      `json:"bankImage"`
	BankImageWeb interface{} `json:"bankImageWeb"`
	BankType     int         `json:"bankType"`
	Color        string      `json:"color"`
}

type HuobiRate struct {
	Сh     string `json: json:"сh"`
	Status string `json:"status"`
	Ts     int    `json:"ts"`
	Tick   struct {
		Id      int       `json:"id"`
		Version int       `json:"version"`
		Open    float64   `json:"open"`
		Close   float64   `json:"close"`
		Low     float64   `json:"low"`
		High    float64   `json:"high"`
		Amount  float64   `json:"amount"`
		Vol     float64   `json:"vol"`
		Count   int       `json:"count"`
		Bid     []float64 `json:"bid"`
		Ask     []float64 `json:"ask"`
	}
}
