package p2p2local

import (
	"encoding/json"
	"github.com/Zmey56/arbitrage/pkg/getdatahuobi"
	"github.com/Zmey56/arbitrage/pkg/getdataokx"
	"github.com/Zmey56/arbitrage/pkg/getinfobinance"
	"github.com/Zmey56/arbitrage/pkg/getlocaldata"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

func getDataBinance(paramUser getlocaldata.ParameterForWorking) {
	//get last modified file
	dir := "jsonresult/Binance"
	searchTerm := "sell_RUB_BNB"
	nameFile, err := getlocaldata.FindLastModifiedFile(dir, searchTerm)
	if err != nil {
		log.Println("Error find and get last modified file")
	}

	file, err := os.Open(nameFile)
	if err != nil {
		log.Println("Error opening JSON file", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		log.Println("Error reading JSON file:", err)
	}

	var orders []getinfobinance.Binance

	err = json.Unmarshal(data, &orders)

	if err != nil {
		log.Println("Error unmarshal JSON data:", err)
	}

	if strings.HasPrefix(searchTerm, "buy") {
		//sort orders by price for buy
		sort.Slice(orders, func(i, j int) bool {
			return orders[i].Data[0].Adv.Price < orders[j].Data[0].Adv.Price &&
				orders[i].Data[len(orders[i].Data)-1].Adv.Price < orders[j].Data[len(orders[i].Data)-1].Adv.Price
		})
	}

	if strings.HasPrefix(searchTerm, "sell") {
		//sort orders by price for sell
		sort.Slice(orders, func(i, j int) bool {
			return orders[i].Data[0].Adv.Price > orders[j].Data[0].Adv.Price &&
				orders[i].Data[len(orders[i].Data)-1].Adv.Price > orders[j].Data[len(orders[i].Data)-1].Adv.Price
		})
	}

	for _, val := range orders {
		for _, subval := range val.Data {
			minTransAmount, _ := strconv.ParseFloat(subval.Adv.MinSingleTransAmount, 64)
			maxTransAmount, _ := strconv.ParseFloat(subval.Adv.MaxSingleTransAmount, 64)
			if minTransAmount > paramUser.Amount || paramUser.Amount > maxTransAmount ||
				subval.Advertiser.MonthOrderCount < paramUser.MonthOrderCount ||
				subval.Advertiser.MonthFinishRate < paramUser.MonthFinishRate {
				continue
			}

			tmpPayment := make([]string, len(paramUser.PayTypes))
			copy(tmpPayment, paramUser.PayTypes)

			if len(tmpPayment) > 0 {
				for _, payment := range subval.Adv.TradeMethods {
					tmpPayment = append(tmpPayment, payment.Identifier)
				}
				if !hasCommonElement(paramUser.PayTypes, tmpPayment) {
					continue
				}
			}

			log.Println(subval.Adv.Price)
			log.Println(subval.Adv.MinSingleTransAmount, " - ", subval.Adv.MaxSingleTransAmount)
			log.Println(subval.Advertiser.MonthOrderCount, " - ", subval.Advertiser.MonthFinishRate*100, "%")
			log.Println(subval.Advertiser.UserType)
			log.Println(tmpPayment)
			log.Println("---------------------------------------------------------------")
		}
	}
}

func getDataHuobi(paramUser getlocaldata.ParameterForWorking) {
	//get last modified file
	dir := "jsonresult/Huobi"
	searchTerm := "sell_RUB_BTC"
	nameFile, err := getlocaldata.FindLastModifiedFile(dir, searchTerm)
	if err != nil {
		log.Println("Error find and get last modified file")
	}

	file, err := os.Open(nameFile)
	if err != nil {
		log.Println("Error opening JSON file", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		log.Println("Error reading JSON file:", err)
	}

	var orders []getdatahuobi.Huobi

	err = json.Unmarshal(data, &orders)

	if err != nil {
		log.Println("Error unmarshal JSON data:", err)
	}

	if strings.HasPrefix(searchTerm, "buy") {
		//sort orders by price for buy
		sort.Slice(orders, func(i, j int) bool {
			return orders[i].Data[0].Price < orders[j].Data[0].Price &&
				orders[i].Data[len(orders[i].Data)-1].Price < orders[j].Data[len(orders[i].Data)-1].Price
		})
	}

	if strings.HasPrefix(searchTerm, "sell") {
		//sort orders by price for sell
		sort.Slice(orders, func(i, j int) bool {
			return orders[i].Data[0].Price > orders[j].Data[0].Price &&
				orders[i].Data[len(orders[i].Data)-1].Price > orders[j].Data[len(orders[i].Data)-1].Price
		})
	}

	// TO DO merchant in Huobi equels 2 not merchant - 1
	for _, val := range orders {
		for _, subval := range val.Data {
			minTransAmount, _ := strconv.ParseFloat(subval.MinTradeLimit, 64)
			maxTransAmount, _ := strconv.ParseFloat(subval.MaxTradeLimit, 64)
			orderCompleteRate, _ := strconv.ParseFloat(subval.OrderCompleteRate, 64)
			if minTransAmount > paramUser.Amount || paramUser.Amount > maxTransAmount ||
				subval.TradeMonthTimes < paramUser.MonthOrderCount ||
				orderCompleteRate < paramUser.MonthFinishRate*100 {
				continue
			}

			tmpPayment := make([]string, len(paramUser.PayTypes))
			copy(tmpPayment, paramUser.PayTypes)

			if len(tmpPayment) > 0 {
				for _, payment := range subval.PayMethods {
					tmpPayment = append(tmpPayment, payment.Name)
				}
				if !hasCommonElement(paramUser.PayTypes, tmpPayment) {
					continue
				}
			}

			log.Println(subval.Price)
			log.Println(subval.MinTradeLimit, " - ", subval.MaxTradeLimit)
			log.Println(subval.TradeMonthTimes, " - ", subval.OrderCompleteRate, "%")
			log.Println(subval.MerchantLevel)
			log.Println(tmpPayment)
			log.Println("---------------------------------------------------------------")
		}
	}
}

func getDataOKX(paramUser getlocaldata.ParameterForWorking) {
	//get last modified file
	dir := "jsonresult/OKX"
	searchTerm := "buy_RUB_BTC"
	nameFile, err := getlocaldata.FindLastModifiedFile(dir, searchTerm)
	if err != nil {
		log.Println("Error find and get last modified file")
	}
	log.Println(nameFile)

	file, err := os.Open(nameFile)
	if err != nil {
		log.Println("Error opening JSON file", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		log.Println("Error reading JSON file:", err)
	}

	var ordersSell []getdataokx.OKXSell
	var ordersBuy []getdataokx.OKXBuy

	if strings.HasPrefix(searchTerm, "buy") {
		log.Println("BUY")
		err = json.Unmarshal(data, &ordersBuy)
		if err != nil {
			log.Println("Error unmarshal JSON data:", err)
		}

		//sort orders by price for buy
		sort.Slice(ordersBuy, func(i, j int) bool {
			return ordersBuy[i].Data.Sell[0].Price < ordersBuy[j].Data.Sell[0].Price &&
				ordersBuy[i].Data.Sell[len(ordersBuy[i].Data.Sell)-1].Price <
					ordersBuy[j].Data.Sell[len(ordersBuy[i].Data.Sell)-1].Price
		})

		for _, val := range ordersBuy {
			for _, subval := range val.Data.Sell {
				minTransAmount, _ := strconv.ParseFloat(subval.QuoteMinAmountPerOrder, 64)
				maxTransAmount, _ := strconv.ParseFloat(subval.QuoteMaxAmountPerOrder, 64)
				completedRate, _ := strconv.ParseFloat(subval.CompletedRate, 64)
				if minTransAmount > paramUser.Amount || paramUser.Amount > maxTransAmount ||
					subval.CompletedOrderQuantity < paramUser.MonthOrderCount ||
					completedRate < paramUser.MonthFinishRate*100 {
					continue
				}

				//var tmpPayment []string

				if len(paramUser.PayTypes) > 0 {
					if !hasCommonElement(paramUser.PayTypes, subval.PaymentMethods) {
						continue
					}
				}

				log.Println(subval.Price)
				log.Println(minTransAmount, " - ", maxTransAmount)
				log.Println(subval.CompletedOrderQuantity, " - ", completedRate, "%")
				log.Println(subval.MerchantID)
				log.Println(subval.PaymentMethods)
				log.Println("---------------------------------------------------------------")
			}
		}
	}

	if strings.HasPrefix(searchTerm, "sell") {
		log.Println("SELL")
		err = json.Unmarshal(data, &ordersSell)
		if err != nil {
			log.Println("Error unmarshal JSON data:", err)
		}

		//sort orders by price for sell
		sort.Slice(ordersSell, func(i, j int) bool {
			priceOne, _ := strconv.ParseFloat(ordersSell[i].Data.Buy[0].Price, 64)
			priceTwo, _ := strconv.ParseFloat(ordersSell[j].Data.Buy[0].Price, 64)
			priceThree, _ := strconv.ParseFloat(ordersSell[i].Data.Buy[len(ordersSell[i].Data.Buy)-1].Price, 64)
			priceFour, _ := strconv.ParseFloat(ordersSell[j].Data.Buy[len(ordersSell[i].Data.Buy)-1].Price, 64)

			return priceOne > priceTwo && priceThree > priceFour
		})

		for _, val := range ordersSell {
			for _, subval := range val.Data.Buy {
				minTransAmount, _ := strconv.ParseFloat(subval.QuoteMinAmountPerOrder, 64)
				maxTransAmount, _ := strconv.ParseFloat(subval.QuoteMaxAmountPerOrder, 64)
				completedRate, _ := strconv.ParseFloat(subval.CompletedRate, 64)
				if minTransAmount > paramUser.Amount || paramUser.Amount > maxTransAmount ||
					subval.CompletedOrderQuantity < paramUser.MonthOrderCount ||
					completedRate < paramUser.MonthFinishRate*100 {
					continue
				}

				if len(paramUser.PayTypes) > 0 {
					if !hasCommonElement(paramUser.PayTypes, subval.PaymentMethods) {
						continue
					}
				}

				log.Println(subval.Price)
				log.Println(minTransAmount, " - ", maxTransAmount)
				log.Println(subval.CompletedOrderQuantity, " - ", completedRate, "%")
				log.Println(subval.MerchantID)
				log.Println(subval.PaymentMethods)
				log.Println("---------------------------------------------------------------")
			}
		}
	}

}

func hasCommonElement(arr1, arr2 []string) bool {
	for _, a := range arr1 {
		for _, b := range arr2 {
			if a == b {
				return true
			}
		}
	}
	return false
}
