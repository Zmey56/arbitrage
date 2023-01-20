package result

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

const chatID = -1001592565485

func SendTextToTelegramChat(chatId int, text string) (string, error) {

	log.Printf("Sending %s to chat_id: %d", text, chatId)
	var telegramApi string = "https://api.telegram.org/bot" + os.Getenv("TELEGRAM_BOT_TOKEN") + "/sendMessage"
	log.Println(telegramApi)

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
	log.Printf("Body of Telegram Response: %s", bodyString)

	return bodyString, nil
}

// https://p2p.binance.com/en/trade/sell/BNB?fiat=RUB&payment=ALL
func FormatMessageAndSend(r ResultP2P, fs, ss string) {
	text := fmt.Sprintf(
		chooseFlag(r.Fiat)+"<b><u>%s</u></b>\n"+
			"\n"+
			"Data and Time  %s \n"+
			"\n"+
			"<i>FIRST STEP - %s</i> - %s\n"+
			"Assets to buy  %s by price %f\n"+
			"Total number of ads for buy: %v\n"+
			"Payment(s): %s \n"+
			"%s\n"+
			"\n"+
			"<i>SECOND STEP - %s</i>\n"+
			"Pair %s by price %f\n"+
			"%s\n"+
			"\n"+
			"<i>THIRD STEP - %s</i> %s\n"+
			"Assets to sell  %s by price %f \n"+
			"Total number of ads for sale %v\n"+
			"Payment(s): %s \n"+
			"%s\n"+
			"\n"+
			"Your profit is <b>%.2f</b> (%.2f %%)",
		r.Fiat,
		r.DataTime.Format("2006/01/02 15:04:05"),
		r.Market.First, fs,
		r.AssetsBuy, r.PriceAssetsBuy, r.TotalAdvBuy,
		strings.Join(r.PaymentBuy, ", "),
		r.LinkAssetsBuy,
		r.Market.Second,
		r.Pair, r.PricePair,
		r.LinkMarket,
		r.Market.Third,
		ss,
		r.AssetsSell, r.PriceAssetsSell, r.TotalAdvSell,
		strings.Join(r.PaymentSell, ", "),
		r.LinkAssetsSell,
		r.ProfitValue, r.ProfitPercet)
	fmt.Println(text)

	SendTextToTelegramChat(chatID, text)
}

func FormatMessageAndSend2steps(r ResultP2P, fs, ss string) {
	text := fmt.Sprintf(
		chooseFlag(r.Fiat)+"<b><u>%s</u></b>\n"+
			"\n"+
			"Data and Time  %s \n"+
			"\n"+
			"<i>FIRST STEP - %s</i> - %s\n"+
			"Assets to buy  %s by price %f\n"+
			"Total number of ads for buy: %v\n"+
			"Payment(s): %s \n"+
			"%s\n"+
			"\n"+
			"<i>SECOND STEP - %s</i> %s\n"+
			"Assets to sell  %s by price %f \n"+
			"Total number of ads for sale %v\n"+
			"Payment(s): %s \n"+
			"%s\n"+
			"\n"+
			"Your profit is <b>%.2f</b> (%.2f %%)",
		r.Fiat,
		r.DataTime.Format("2006/01/02 15:04:05"),
		r.Market.First, fs,
		r.AssetsBuy, r.PriceAssetsBuy, r.TotalAdvBuy,
		strings.Join(r.PaymentBuy, ", "),
		r.LinkAssetsBuy,
		r.Market.Second,
		ss,
		r.AssetsSell, r.PriceAssetsSell, r.TotalAdvSell,
		strings.Join(r.PaymentSell, ", "),
		r.LinkAssetsSell,
		r.ProfitValue, r.ProfitPercet)
	fmt.Println(text)

	SendTextToTelegramChat(chatID, text)
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
	case "UZS":
		return "🇺🇿"
	case "VND":
		return "🇻🇳"
	default:
		return ""
	}
}
