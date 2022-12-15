package telegramm

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

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

//Profit          bool
//DataTime        time.Time
//Fiat            string
//AssetsBuy       string
//PriceAssetsBuy  float64
//PaymentBuy      []string
//LinkAssetsBuy   string
//Pair            string
//PricePair       float64
//LinkMarket      string
//AssetsSell      string
//PriceAssetsSell float64
//PaymentSell     []string
//LinkAssetsSell  string
//ProfitValue     float64
//ProfitPercet    string

//func FormatMessage(m working.ResultP2P) {
//	Fiat            string
//	//AssetsBuy       string
//	//PriceAssetsBuy  float64
//	//PaymentBuy      []string
//	//LinkAssetsBuy   string
//	//Pair            string
//	//PricePair       float64
//	//LinkMarket      string
//	//AssetsSell      string
//	//PriceAssetsSell float64
//	//PaymentSell     []string
//	//LinkAssetsSell  string
//	//ProfitValue     float64
//	//ProfitPercet    string
//}
