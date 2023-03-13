// choose available payment methods
package workingbinance

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func InputCommandLineBinance(fiat string) ParametersBinance {
	file_path := "cmd/enterparam/parambinance/" + fiat + ".json"
	paramUser := ParametersBinance{}
	n := bufio.NewReader(os.Stdout)
	fmt.Println("Do you want to enter new params? If you want - enter Yes")
	readagree, err := n.ReadString('\n')
	if err != nil {
		log.Println("User try to enter wong value:", err)
	}
	if strings.ToLower(strings.TrimRight(readagree, "\n")) == "yes" {

		//selected user
		var userpayment []string
		for {
			fmt.Println("Enter the balance (integer only) of the money(asset):")
			readbalance, err := n.ReadString('\n')
			if readbalance == "\n" {
				log.Println("You don't enter value")
				paramUser.TransAmount = ""
				break
			}
			if err != nil {
				log.Printf("Problem with enter the balance and error: %v", err)
			}
			readbalance = strings.TrimSpace(readbalance)
			var digitCheck = regexp.MustCompile(`^[0-9]+$`)
			if digitCheck.MatchString(readbalance) {
				paramUser.TransAmount = readbalance
				break
			} else {
				fmt.Println("You entered the wrong value", readbalance)
			}
		}

		//MERCHANT
		fmt.Println("Do you want to choose only the merchant? If you want - enter Yes")

		readmerchant, err := n.ReadString('\n')
		if err != nil {
			log.Println("User try to enter wrong value:", err)
		}
		if strings.ToLower(strings.TrimRight(readmerchant, "\n")) == "yes" {
			paramUser.PublisherType = "merchant"
		}

		//BORDER
		fmt.Println("What is the minimum limit for the number of ads to set?")

		readborder, _ := n.ReadString('\n')
		readborder = strings.TrimSpace(readborder)
		for {
			if readborder != "" {
				var f int
				f, err := strconv.Atoi(readborder)
				if err != nil {
					log.Println("User try to enter wrong value of percent to track:", err)
				} else {
					log.Println(f)
					paramUser.Border = f
					break
				}
			} else {
				paramUser.Border = 0
				break
			}
		}

		//PERCENT
		fmt.Println("What percentage will you track?(0.01 - 100)")

		readpercent, _ := n.ReadString('\n')
		readpercent = strings.TrimSpace(readpercent)
		for {
			if readpercent != "" {
				var f float64
				f, err = strconv.ParseFloat(readpercent, 64)
				if err != nil {
					log.Println("User try to enter wrong value of percent to track:", err)
				} else {
					log.Println(f)
					paramUser.PercentUser = f
					break
				}
			} else {
				paramUser.PercentUser = 0
				break
			}
		}
		if strings.ToLower(strings.TrimRight(readmerchant, "\n")) == "yes" {
			paramUser.PublisherType = "merchant"
		}

		gpm := GetPaymentFromJSONBinance(fiat)

		var availablepayment []string

		for i, _ := range gpm {
			fmt.Println(i+1, " - ", gpm[i].TradeMethodName)
			availablepayment = append(availablepayment, gpm[i].Identifier)
		}

		fmt.Println("Available means of payment, choose number(If you want to finished enter \"7777\"):")
		for {
			readnumber, _ := n.ReadString('\n')
			readnumber = strings.TrimSpace(readnumber)
			number, err := strconv.Atoi(readnumber)
			if err != nil {
				log.Println(err)
			}

			if number == 7777 {
				break
			} else if number >= len(availablepayment) || number < 1 {
				fmt.Println("You have selected the wrong number")
				continue
			} else {
				fmt.Println(availablepayment[number-1])
				userpayment = append(userpayment, availablepayment[number-1])
			}
		}
		paramUser.PayTypes = unique(userpayment)
		file, err := json.MarshalIndent(paramUser, "", " ")
		if err != nil {
			log.Println("Can't marshalIndent paramUser", err)
		}
		_ = os.WriteFile(file_path, file, 0644)
	} else {
		file, err := os.ReadFile(file_path)
		if err != nil {
			log.Println("Can't read file with parameters", err)
		}
		_ = json.Unmarshal([]byte(file), &paramUser)
	}

	return paramUser
}

func unique(arr []string) []string {
	occurred := map[string]bool{}
	result := []string{}
	for e := range arr {

		// check if already the mapped
		// variable is set to true or not
		if occurred[arr[e]] != true {
			occurred[arr[e]] = true

			// Append to result slice.
			result = append(result, arr[e])
		}
	}

	return result
}

// file exist?
func Exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

// "AED", "AMD", "AZN", "EUR", "GEL", "KZT", "RUB", "TRY", "UAH", "UZS"
func GetPaymentFromJSONBinance(fiat string) PaymentsBinance {
	payment := ""
	switch fiat {
	case "AED":
		payment = fmt.Sprintf("data/databinance/%s/%s_payment.json", fiat, fiat)
	case "AMD":
		payment = fmt.Sprintf("data/databinance/%s/%s_payment.json", fiat, fiat)
	case "AZN":
		payment = fmt.Sprintf("data/databinance/%s/%s_payment.json", fiat, fiat)
	case "ARS":
		payment = fmt.Sprintf("data/databinance/%s/%s_payment.json", fiat, fiat)
	case "EUR":
		payment = fmt.Sprintf("data/databinance/%s/%s_payment.json", fiat, fiat)
	case "GEL":
		payment = fmt.Sprintf("data/databinance/%s/%s_payment.json", fiat, fiat)
	case "KZT":
		payment = fmt.Sprintf("data/databinance/%s/%s_payment.json", fiat, fiat)
	case "RUB":
		payment = fmt.Sprintf("data/databinance/%s/%s_payment.json", fiat, fiat)
	case "TRY":
		payment = fmt.Sprintf("data/databinance/%s/%s_payment.json", fiat, fiat)
	case "UAH":
		payment = fmt.Sprintf("data/databinance/%s/%s_payment.json", fiat, fiat)
	case "USD":
		payment = fmt.Sprintf("data/databinance/%s/%s_payment.json", fiat, fiat)
	case "UZS":
		payment = fmt.Sprintf("data/databinance/%s/%s_payment.json", fiat, fiat)
	case "VND":
		payment = fmt.Sprintf("data/databinance/%s/%s_payment.json", fiat, fiat)
	default:
		fmt.Printf("For %v don't have payments methods\n", fiat)
	}
	jsonfile, err := os.ReadFile(payment)
	if err != nil {
		panic(err)
	}
	allpayments := PaymentsBinance{}
	_ = json.Unmarshal(jsonfile, &allpayments)

	return allpayments
}
