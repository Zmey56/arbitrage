package getdata

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

const rubpayment = "data/rub_payment.json"

type Payments []struct {
	Identifier           string `json:"Identifier"`
	PayAccount           string `json:"PayAccount"`
	PayMethodID          string `json:"PayMethodId"`
	PayType              string `json:"PayType"`
	TradeMethodName      string `json:"TradeMethodName"`
	TradeMethodShortName string `json:"TradeMethodShortName"`
}

func GetPaymentFromJSON(fiat string) Payments {
	payment := ""
	switch fiat {
	case "RUB":
		payment = rubpayment
	default:
		fmt.Printf("For %v don't have payments methods\n", fiat)
	}
	jsonfile, err := os.ReadFile(payment)
	if err != nil {
		panic(err)
	}
	allpayments := Payments{}
	_ = json.Unmarshal(jsonfile, &allpayments)

	return allpayments
}

func GetFilenameDate(name, ext string) string {
	// Use layout string for time format.
	const layout = "01-02-2006"
	// Place now in the string.
	t := time.Now()
	return name + t.Format(layout) + "." + ext
}
