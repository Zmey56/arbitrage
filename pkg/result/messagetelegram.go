package result

import (
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const chatID = -1001592565485

//const chatID = -993812776

const TELEGRAM_BOT_TOKEN = "5763797414:AAHJ8exgiqxHuW44SyEr15fKsWKPixNofVg"

//const TELEGRAM_BOT_TOKEN = "6072584429:AAGkRNgzQSZNJ9VMpqkol-r1H2D7jCZNuVA"

func SendTextToTelegramChat(chatId int, text string) (string, error) {

	//log.Printf("Sending %s to chat_id: %d", text, chatId)
	//var telegramApi string = "https://api.telegram.org/bot" + os.Getenv("TELEGRAM_BOT_TOKEN") + "/sendMessage"
	var telegramApi string = "https://api.telegram.org/bot" + TELEGRAM_BOT_TOKEN + "/sendMessage"
	//log.Println(telegramApi)

	response, err := http.PostForm(
		telegramApi,
		url.Values{
			"chat_id":                  {strconv.Itoa(chatId)},
			"parse_mode":               {"HTML"},
			"text":                     {text},
			"disable_web_page_preview": {"true"},
		})

	if err != nil {
		log.Printf("error when posting text to the chat: %s", err.Error())
		return "", err
	}
	defer response.Body.Close()

	var bodyBytes, errRead = io.ReadAll(response.Body)
	if errRead != nil {
		log.Printf("error in parsing telegram answer %s", errRead.Error())
		return "", err
	}
	bodyString := string(bodyBytes)
	//log.Printf("Body of Telegram Response: %s", bodyString)

	return bodyString, nil
}

// https://p2p.binance.com/en/trade/sell/BNB?fiat=RUB&payment=ALL
func FormatMessageAndSend(r ResultP2P) {
	tmpAmount, _ := strconv.ParseFloat(r.Amount, 64)
	text := fmt.Sprintf(
		chooseFlag(r.Fiat)+"<b><u>%s</u></b>\n"+
			"\n"+
			"Data and Time  %s \n"+
			"\n"+
			"Your profit is <b>%.2f</b> (%.2f %%) from <b>%s</b> %s\n"+
			"\n"+
			"<i>FIRST STEP - <b>%s</b></i> - You are %s\n"+
			"Assets to buy  <b>%s</b> by price <b>%s</b>\n"+
			"%s\n"+
			"<u>Payment(s):</u> %s \n"+
			"\n"+
			"%s\n"+
			"\n"+
			"<i>SECOND STEP - <b>%s</b></i>\n"+
			"Pair <b>%s</b> by price <b>%s</b>\n"+
			"\n"+
			"%s\n"+
			"\n"+
			"<i>THIRD STEP - <b>%s</b></i> You are %s\n"+
			"Assets to sell <b>%s</b> by price <b>%s</b> \n"+
			"%s\n"+
			"<u>Payment(s):</u> %s \n"+
			"\n"+
			"%s\n",
		r.Fiat,
		r.DataTime.Format("2006/01/02 15:04:05"),
		r.ProfitValue, r.ProfitPercet, formattedNum(tmpAmount), r.Fiat,
		r.Market.First, r.User.FirstUser,
		r.AssetsBuy, formattedNum(r.PriceAssetsBuy), isMerchantVerifiedFirst(r),
		strings.Join(r.PaymentBuy, ", "),
		r.LinkAssetsBuy,
		r.Market.Second,
		strings.ToUpper(r.Pair), formattedNum(r.PricePair),
		r.LinkMarket,
		r.Market.Third,
		r.User.ThirdUser,
		strings.ToUpper(r.AssetsSell), formattedNum(r.PriceAssetsSell), isMerchantVerifiedThird(r),
		strings.Join(r.PaymentSell, ", "),
		r.LinkAssetsSell)

	fmt.Println(text)

	SendTextToTelegramChat(chatID, text)
}

func FormatMessageAndSend2steps(r ResultP2P) {
	tmpAmount, _ := strconv.ParseFloat(r.Amount, 64)
	text := fmt.Sprintf(
		chooseFlag(r.Fiat)+"<b><u>%s</u></b>\n"+
			"\n"+
			"Data and Time  %s \n"+
			"\n"+
			"Your profit is <b>%.2f</b> (%.2f %%) from <b>%s</b> %s\n"+
			"\n"+
			"<i>FIRST STEP - <b>%s</b></i> - You are %s\n"+
			"Assets to buy  <b>%s</b> by price <b>%s</b>\n"+
			"%s\n"+
			"<u>Payment(s):</u> %s \n"+
			"\n"+
			"%s\n"+
			"\n"+
			"<i>SECOND STEP - <b>%s</b></i> You are %s\n"+
			"Assets to sell <b>%s</b> by price <b>%s</b> \n"+
			"%s\n"+
			"<u>Payment(s):</u> %s \n"+
			"\n"+
			"%s\n",
		r.Fiat,
		r.DataTime.Format("2006/01/02 15:04:05"),
		r.ProfitValue, r.ProfitPercet, formattedNum(tmpAmount), r.Fiat,
		r.Market.First, r.User.FirstUser,
		r.AssetsBuy, formattedNum(r.PriceAssetsBuy), isMerchantVerifiedFirst(r),
		strings.Join(r.PaymentBuy, ", "),
		r.LinkAssetsBuy,
		r.Market.Third,
		r.User.ThirdUser,
		strings.ToUpper(r.AssetsSell), formattedNum(r.PriceAssetsSell), isMerchantVerifiedThird(r),
		strings.Join(r.PaymentSell, ", "),
		r.LinkAssetsSell)

	fmt.Println(text)

	SendTextToTelegramChat(chatID, text)
}

func isMerchantVerifiedFirst2steps(r ResultP2P2steps) string {
	if r.Merchant == true {
		return fmt.Sprintf("The lowest buying among <b><u>%v</u></b> and highest selling price"+
			" among <b><u>%v</u></b> <b>Verified</b> advertisers. Delta %.2f %%", r.AdvToalBuy, r.AdvToalSell,
			r.DeltaADV)
	} else {
		return fmt.Sprintf("The lowest buying among <b><u>%v</u></b> and highest selling price"+
			" among <b><u>%v</u></b> <b>ALL</b> advertisers. Delta %.2f %%", r.AdvToalBuy, r.AdvToalSell,
			r.DeltaADV)
	}
}

func isMerchantVerifiedFirst(r ResultP2P) string {
	if r.Merchant.FirstMerch == true && r.User.FirstUser == "Taker" {
		return fmt.Sprintf("The lowest buying price among <b><u>%v</u></b> <b>Verified</b> advertisers", r.TotalAdvBuy)
	} else if r.Merchant.FirstMerch == false && r.User.FirstUser == "Taker" {
		return fmt.Sprintf("The lowest buying price among <b><u>%v</u></b> <b>All</b> advertisers", r.TotalAdvBuy)
	} else {
		return ""
	}
}

func isMerchantVerifiedThird(r ResultP2P) string {
	if r.Merchant.ThirdMerch == true && r.User.ThirdUser == "Taker" {
		return fmt.Sprintf("The highest selling price among <b><u>%v</u></b> <b>Verified</b> advertisers", r.TotalAdvSell)
	} else if r.Merchant.ThirdMerch == false && r.User.ThirdUser == "Taker" {
		return fmt.Sprintf("The highest selling price among <b><u>%v</u></b> <b>All</b> advertisers", r.TotalAdvSell)
	} else if r.Merchant.ThirdMerch == true && r.User.ThirdUser == "Maker" {
		return fmt.Sprintf("The lowest buying price among <b><u>%v</u></b> <b>Verified</b> advertisers below which the ad is placed", r.TotalAdvSell)
	} else if r.Merchant.ThirdMerch == false && r.User.ThirdUser == "Maker" {
		return fmt.Sprintf("The lowest buying price among <b><u>%v</u></b> <b>All</b> advertisers below which the ad is placed", r.TotalAdvSell)
	} else {
		return ""
	}
}

func chooseFlag(fiat string) string {
	switch fiat {
	case "AED":
		return "🇦🇪"
	case "AMD":
		return "🇦🇲"
	case "AZN":
		return "🇦🇿"
	case "ARS":
		return "🇦🇷"
	case "EUR":
		return "🇪🇺"
	case "GEL":
		return "🇬🇪"
	case "KZT":
		return "🇰🇿"
	case "RUB":
		return "🇷🇺"
	case "TRY":
		return "🇹🇷"
	case "UAH":
		return "🇺🇦"
	case "USD":
		return "🇺🇸"
	case "UZS":
		return "🇺🇿"
	case "VND":
		return "🇻🇳"
	default:
		return ""
	}
}

func formattedNum(num float64) string {
	if num > 1 {
		num = math.Round(num*10000) / 10000
	} else {
		num = math.Round(num*1000000) / 1000000
	}
	formattedNum := strconv.FormatFloat(num, 'f', -1, 64)
	parts := strings.Split(formattedNum, ".")
	if len(parts) > 1 {
		decimalPart := parts[1]
		i := 0
		for i < len(decimalPart) && decimalPart[i] == '0' {
			i++
		}
		decimalPart = decimalPart[:i] + strings.TrimRightFunc(decimalPart[i:], func(r rune) bool {
			return r == '0'
		})
		formattedNum = parts[0] + "." + decimalPart
	}

	return formattedNum
}

//func FormatMessageAndSend2steps(r ResultP2P2steps) {
//	text := fmt.Sprintf(
//		chooseFlag(r.FiatUnit)+"<b><u>%s</u></b> - <b><u>%s</u></b>\n"+
//			"\n"+
//			"Market one  <b>%s</b> (%s) and two <b>%s</b> (%s)\n"+
//			"\n"+
//			"Data and Time  %s \n"+
//			"\n"+
//			"Price BUY <b>%s</b> and SELL <b>%s</b> (Delta %.2f %%)\n"+ //+
//			"\n"+
//			"<i>Second price BUY %s(%.2f %%) and SELL %s(%.2f %%)</i>\n"+ //+
//			"\n"+
//			"Mean price BUY <b>%s</b> and SELL <b>%s</b> (Delta %.2f %%)\n"+
//			"\n"+
//			"<i>Standard deviation price BUY <b>%s</b> (%.2f %%) and SELL <b>%s</b></i> (%.2f %%) (Delta: %.2f %%)\n"+
//			"\n"+
//			"Weight mean price BUY <b>%s</b> and SELL <b>%s</b> (Delta %.2f %%) \n"+
//			"\n"+
//			"Weighted Average <b>%s</b> SD: <b>%s</b> (Delta %.2f %%)\n"+
//			"\n"+
//			"Resistance BUY <b>%s</b>(%.2f %%) and SELL <b>%s</b>(%.2f %%) (Delta %.2f %%)\n"+
//			"\n"+
//			"<i>Volume resistance BUY <b>%s</b> and SELL <b>%s</b></i>\n"+
//			"\n"+
//			"%s\n"+
//			"\n"+
//			"Payment(s) BUY: %s \n"+
//			"\n"+
//			"Payment(s) SELL: %s \n"+
//			"\n"+
//			"Amount: %s %s",
//		r.FiatUnit, r.Asset,
//		r.MarketOne, r.User.FirstUser, r.MarketTwo, r.User.SecondUser,
//		r.DataTime.Format("2006/01/02 15:04:05"),
//		formattedNum(r.PriceB), formattedNum(r.PriceS), r.DeltaBuySell,
//		formattedNum(r.PriceBSecond), r.DeltaFirstSecondB, formattedNum(r.PriceSSecond), r.DeltaFirstSecondS,
//		formattedNum(r.MeanPriceB), formattedNum(r.MeanPriceS), r.DeltaMean,
//		formattedNum(r.SDPriceB), r.DeltaSDb, formattedNum(r.SDPriceS), r.DeltaSDs, r.DeltaSD,
//		formattedNum(r.MeanWeighB), formattedNum(r.MeanWeighS), r.DeltaMeanWeight,
//		formattedNum(r.MeanWeight), formattedNum(r.MeanWeightSD), r.DeltaWSD,
//		formattedNum(r.GiantPriceB), r.DeltaGiantPriceB, formattedNum(r.GiantPriceS), r.DeltaGiantPriceS, r.DeltaGiant,
//		formattedNum(r.GiantVolB), formattedNum(r.GiantVolS),
//		isMerchantVerifiedFirst2steps(r),
//		strings.Join(r.PaymentBuy, ", "),
//		strings.Join(r.PaymentSell, ", "),
//		formattedNum(r.Amount), r.FiatUnit)
//
//	log.Println("RESULT", text)
//
//	SendTextToTelegramChat(chatID, text)
//}

//func isMerchantVerifiedFirst2steps(r ResultP2P2Steps) string {
//	if r.Merchant.FirstMerch == true && r.User.FirstUser == "Taker" {
//		return fmt.Sprintf("The lowest buying price among <b><u>%v</u></b> <b>Verified</b> advertisers", r.TotalAdvBuy)
//	} else if r.Merchant.FirstMerch == false && r.User.FirstUser == "Taker" {
//		return fmt.Sprintf("The lowest buying price among <b><u>%v</u></b> <b>All</b> advertisers", r.TotalAdvBuy)
//	} else {
//		return ""
//	}
//}
//
//func isMerchantVerifiedSecond2steps(r ResultP2P2Steps) string {
//	if r.Merchant.SecondMerch == true && r.User.SecondUser == "Taker" {
//		return fmt.Sprintf("The highest selling price among <b><u>%v</u></b> <b>Verified</b> advertisers", r.TotalAdvSell)
//	} else if r.Merchant.SecondMerch == false && r.User.SecondUser == "Taker" {
//		return fmt.Sprintf("The highest selling price among <b><u>%v</u></b> <b>All</b> advertisers", r.TotalAdvSell)
//	} else if r.Merchant.SecondMerch == true && r.User.SecondUser == "Maker" {
//		return fmt.Sprintf("The lowest buying price among <b><u>%v</u></b> <b>Verified</b> advertisers below which the ad is placed", r.TotalAdvSell)
//	} else if r.Merchant.SecondMerch == false && r.User.SecondUser == "Maker" {
//		return fmt.Sprintf("The lowest buying price among <b><u>%v</u></b> <b>All</b> advertisers below which the ad is placed", r.TotalAdvSell)
//	} else {
//		return ""
//	}
//}
//
//func FormatMessageAndSend2steps(r ResultP2P2steps) {
//	text := fmt.Sprintf(
//		chooseFlag(r.Fiat)+"<b><u>%s</u></b>\n"+
//			"\n"+
//			"Data and Time  %s \n"+
//			"Your profit is <b>%.2f</b> (%.2f %%)\n"+ //+
//			"\n"+
//			"<i>FIRST STEP - %s</i> - %s\n"+ //+
//			"\n"+
//			"Assets to buy  %s by price %f \n"+
//			"\n"+
//			"%s\n"+
//			"\n"+
//			"Payment(s): %s \n"+
//			"\n"+
//			"%s\n"+
//			"\n"+
//			"<i>SECOND STEP - %s</i> %s\n"+
//			"\n"+
//			"Assets to sell  %s by price %f \n"+
//			"\n"+
//			"%s\n"+
//			"\n"+
//			"Payment(s): %s \n"+
//			"\n"+
//			"%s\n"+
//			"\n"+
//			"\n"+
//			"Your profit is <b>%.2f</b> (%.2f %%)",
//		r.FiatUnit,
//		r.DataTime.Format("2006/01/02 15:04:05"),
//		r.ProfitValue, r.ProfitPercet,
//		r.Market.First, r.User.FirstUser,
//		r.AssetsBuy, r.PriceAssetsBuy, isMerchantVerifiedFirst2steps(r),
//		strings.Join(r.PaymentBuy, ", "),
//		r.LinkAssetsBuy,
//		r.Market.Second, r.User.SecondUser,
//		r.AssetsSell, r.PriceAssetsSell, isMerchantVerifiedSecond2steps(r),
//		strings.Join(r.PaymentSell, ", "),
//		r.LinkAssetsSell)
//	//fmt.Println(text)
//
//	SendTextToTelegramChat(chatID, text)
//}
