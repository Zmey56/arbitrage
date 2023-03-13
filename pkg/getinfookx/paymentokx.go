package getinfookx

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func GetPeymontMethodsOKX(fiat string) []PaymentOKX {

	url := fmt.Sprintf("https://www.okx.com/v3/c2c/configs/receipt/templates?t=1676884354591&quoteCurrency=%s", fiat)

	response, err := http.Get(url)
	if err != nil {
		log.Printf("The HTTP request failed with error %s\n", err)
		//return payment, err
	}
	defer response.Body.Close()

	data, _ := io.ReadAll(response.Body)
	var po PaymentsOKX

	json.Unmarshal(data, &po)

	allpayment := []PaymentOKX{}
	cp := PaymentOKX{}
	for _, c := range po.Data {
		cp.TradeMethodName = c.PaymentMethod
		cp.TradeMethodShortName = c.PaymentMethodDescription
		allpayment = append(allpayment, cp)
	}

	savePaymentToJSONOKX(allpayment, fiat)

	return allpayment
}

func savePaymentToJSONOKX(p []PaymentOKX, fiat string) {
	file, err := json.MarshalIndent(p, "", " ")
	if err != nil {
		panic(err)
	}
	path_json := fmt.Sprintf("data/dataokx/%s/%s_payment.json", fiat, fiat)
	os.WriteFile(path_json, file, 0644)
}
