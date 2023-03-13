package getinfobinancehuobi

import (
	"encoding/json"
	"fmt"
	"github.com/Zmey56/arbitrage/pkg/getdata"
	"github.com/Zmey56/arbitrage/pkg/getdatahuobi"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

type ExchangeInfoBinance struct {
	Timezone   string `json:"timezone"`
	ServerTime int64  `json:"serverTime"`
	RateLimits []struct {
		RateLimitType string `json:"rateLimitType"`
		Interval      string `json:"interval"`
		IntervalNum   int    `json:"intervalNum"`
		Limit         int    `json:"limit"`
	} `json:"rateLimits"`
	ExchangeFilters []interface{} `json:"exchangeFilters"`
	Symbols         []struct {
		Symbol                     string   `json:"symbol"`
		Status                     string   `json:"status"`
		BaseAsset                  string   `json:"baseAsset"`
		BaseAssetPrecision         int      `json:"baseAssetPrecision"`
		QuoteAsset                 string   `json:"quoteAsset"`
		QuotePrecision             int      `json:"quotePrecision"`
		QuoteAssetPrecision        int      `json:"quoteAssetPrecision"`
		BaseCommissionPrecision    int      `json:"baseCommissionPrecision"`
		QuoteCommissionPrecision   int      `json:"quoteCommissionPrecision"`
		OrderTypes                 []string `json:"orderTypes"`
		IcebergAllowed             bool     `json:"icebergAllowed"`
		OcoAllowed                 bool     `json:"ocoAllowed"`
		QuoteOrderQtyMarketAllowed bool     `json:"quoteOrderQtyMarketAllowed"`
		AllowTrailingStop          bool     `json:"allowTrailingStop"`
		CancelReplaceAllowed       bool     `json:"cancelReplaceAllowed"`
		IsSpotTradingAllowed       bool     `json:"isSpotTradingAllowed"`
		IsMarginTradingAllowed     bool     `json:"isMarginTradingAllowed"`
		Filters                    []struct {
			FilterType            string `json:"filterType"`
			MinPrice              string `json:"minPrice,omitempty"`
			MaxPrice              string `json:"maxPrice,omitempty"`
			TickSize              string `json:"tickSize,omitempty"`
			MinQty                string `json:"minQty,omitempty"`
			MaxQty                string `json:"maxQty,omitempty"`
			StepSize              string `json:"stepSize,omitempty"`
			MinNotional           string `json:"minNotional,omitempty"`
			ApplyToMarket         bool   `json:"applyToMarket,omitempty"`
			AvgPriceMins          int    `json:"avgPriceMins,omitempty"`
			Limit                 int    `json:"limit,omitempty"`
			MinTrailingAboveDelta int    `json:"minTrailingAboveDelta,omitempty"`
			MaxTrailingAboveDelta int    `json:"maxTrailingAboveDelta,omitempty"`
			MinTrailingBelowDelta int    `json:"minTrailingBelowDelta,omitempty"`
			MaxTrailingBelowDelta int    `json:"maxTrailingBelowDelta,omitempty"`
			BidMultiplierUp       string `json:"bidMultiplierUp,omitempty"`
			BidMultiplierDown     string `json:"bidMultiplierDown,omitempty"`
			AskMultiplierUp       string `json:"askMultiplierUp,omitempty"`
			AskMultiplierDown     string `json:"askMultiplierDown,omitempty"`
			MaxNumOrders          int    `json:"maxNumOrders,omitempty"`
			MaxNumAlgoOrders      int    `json:"maxNumAlgoOrders,omitempty"`
		} `json:"filters"`
		Permissions                     []string `json:"permissions"`
		DefaultSelfTradePreventionMode  string   `json:"defaultSelfTradePreventionMode"`
		AllowedSelfTradePreventionModes []string `json:"allowedSelfTradePreventionModes"`
	} `json:"symbols"`
}

func GetListPairHBB(fiat string) {
	assetscoinpair := make(map[string][]string)

	//Get Assets from Binance
	assets := getdata.GetAssets(fiat)
	assets_symbol := make([]string, 0, len(assets))

	for k, _ := range assets {
		assets_symbol = append(assets_symbol, k)
	}

	//Get Coin from Huobi
	coinarr := getdatahuobi.GetCurrencyHuobi(fiat)

	url := "https://api.binance.com/api/v3/exchangeInfo"

	resp, err := http.Get(url)

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	eib := ExchangeInfoBinance{}

	//var dat map[string]interface{}
	//
	if err := json.Unmarshal(body, &eib); err != nil {
		panic(err)
	}

	for _, f := range coinarr {
		pair := []string{}
		for _, i := range eib.Symbols {
			for _, j := range assets_symbol {
				if strings.ToLower(i.BaseAsset) == strings.ToLower(f) && strings.ToLower(i.QuoteAsset) == strings.ToLower(j) {
					pair = append(pair, i.Symbol)
				}
				if strings.ToLower(i.QuoteAsset) == strings.ToLower(f) && strings.ToLower(i.BaseAsset) == strings.ToLower(j) {
					pair = append(pair, i.Symbol)
				}
			}

		}
		assetscoinpair[f] = pair
	}

	name_json := fmt.Sprintf("data/datahbb/%s/%s_pair.json", fiat, fiat)
	jsonStr, err := json.MarshalIndent(assetscoinpair, "", " ")
	if err != nil {
		log.Printf("Error: %s", err.Error())
	}

	_ = os.WriteFile(name_json, jsonStr, 0644)
}
