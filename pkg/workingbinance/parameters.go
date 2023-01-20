package workingbinance

type ParametersBinance struct {
	PayTypes      []string `json:"payTypes"`
	TransAmount   string   `json:"transAmount"`
	PublisherType string   `json:"publisher_type"`
	PercentUser   float64  `json:"percentUser"`
}

type PaymentsBinance []struct {
	Identifier           string `json:"Identifier"`
	PayAccount           string `json:"PayAccount"`
	PayMethodID          string `json:"PayMethodId"`
	PayType              string `json:"PayType"`
	TradeMethodName      string `json:"TradeMethodName"`
	TradeMethodShortName string `json:"TradeMethodShortName"`
}
