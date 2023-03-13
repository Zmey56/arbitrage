package result

import "time"

type ResultP2P2steps struct {
	Asset             string
	FiatUnit          string
	MarketOne         string
	MarketTwo         string
	DataTime          time.Time
	Merchant          bool
	PriceB            float64
	PriceS            float64
	PriceBSecond      float64
	PriceSSecond      float64
	MeanPriceB        float64
	MeanPriceS        float64
	SDPriceB          float64
	SDPriceS          float64
	MeanWeighB        float64
	MeanWeighS        float64
	GiantPriceB       float64
	GiantVolB         float64
	GiantPriceS       float64
	GiantVolS         float64
	AdvToalBuy        int
	AdvToalSell       int
	DeltaBuySell      float64
	DeltaFirstSecondB float64
	DeltaFirstSecondS float64
	DeltaMean         float64
	DeltaMeanWeight   float64
	DeltaGiantPriceB  float64
	DeltaGiantPriceS  float64
	DeltaGiant        float64
	DeltaADV          float64
	User              struct {
		FirstUser  string
		SecondUser string
	}
	PaymentBuy  []string
	PaymentSell []string
	Amount      float64
}

type ResultP2P struct {
	Amount string
	User   struct {
		FirstUser  string
		SecondUser string
		ThirdUser  string
	}
	Market struct {
		First  string
		Second string
		Third  string
	}
	Merchant struct {
		FirstMerch  bool
		SecondMerch bool
		ThirdMerch  bool
	}
	Profit          bool
	DataTime        time.Time
	Fiat            string
	AssetsBuy       string
	PriceAssetsBuy  float64
	PaymentBuy      []string
	LinkAssetsBuy   string
	Pair            string
	PricePair       float64
	LinkMarket      string
	AssetsSell      string
	PriceAssetsSell float64
	PaymentSell     []string
	LinkAssetsSell  string
	ProfitValue     float64
	ProfitPercet    float64
	TotalAdvBuy     int
	TotalAdvSell    int
	AdvNoBuy        string
	AdvNoSell       string
}

//type ResultP2P2Steps struct {
//	Amount string
//	User   struct {
//		FirstUser  string
//		SecondUser string
//	}
//	Market struct {
//		First  string
//		Second string
//	}
//	Merchant struct {
//		FirstMerch  bool
//		SecondMerch bool
//	}
//	Profit          bool
//	DataTime        time.Time
//	Fiat            string
//	AssetsBuy       string
//	PriceAssetsBuy  float64
//	PaymentBuy      []string
//	LinkAssetsBuy   string
//	AssetsSell      string
//	PriceAssetsSell float64
//	PaymentSell     []string
//	LinkAssetsSell  string
//	ProfitValue     float64
//	ProfitPercet    float64
//	TotalAdvBuy     int
//	TotalAdvSell    int
//	AdvNoBuy        string
//	AdvNoSell       string
//}
