package getinfohuobi

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func GetPeymontMethodsHuobi(fiat string) []PaymentHuobi {

	url := "https://otc-cf.huobi.com/v1/data/config-list?type=currency,marketQuery,pay,allCountry,coin"

	//var payment []PayMethodHuobi
	response, err := http.Get(url)
	if err != nil {
		log.Printf("The HTTP request failed with error %s\n", err)
		//return payment, err
	}
	defer response.Body.Close()

	data, _ := io.ReadAll(response.Body)
	var ph PaymentsHuobi

	json.Unmarshal(data, &ph)
	//create map with id and payment
	idpayment := make(map[int]string)

	for _, val := range ph.Data.PayMethod {
		idpayment[val.PayMethodID] = val.Name
	}

	//for _, f := range fiat {
	allpayment := []PaymentHuobi{}
	cp := PaymentHuobi{}
	for _, c := range ph.Data.Currency {
		if c.NameShort == fiat {
			for _, j := range c.SupportPayments {
				tmp := &cp
				tmp.PayMethodId = j
				tmp.TradeMethodName = idpayment[j]
				allpayment = append(allpayment, *tmp)
			}
		}
	}
	return allpayment
}

func SavePaymentToJSONHuobi(p []PaymentHuobi, fiat string) {
	file, err := json.MarshalIndent(p, "", " ")
	if err != nil {
		panic(err)
	}
	path_json := fmt.Sprintf("data/datahuobi/%s/%s_payment.json", fiat, fiat)
	os.WriteFile(path_json, file, 0644)
}
