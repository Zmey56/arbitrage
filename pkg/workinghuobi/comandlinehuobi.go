package workinghuobi

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/Zmey56/arbitrage/pkg/getinfohuobi"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func InputCommandLineHuobi(fiat string) {
	file_path := "cmd/enterparam/paramhuobi/" + fiat + ".json"
	paramUser := getinfohuobi.ParametersHuobi{}
	n := bufio.NewReader(os.Stdout)
	fmt.Println("Do you want to enter new params for Huobi? If you want - enter Yes")
	readagree, err := n.ReadString('\n')
	if err != nil {
		log.Println("User try to enter wrong value:", err)
	}
	if strings.ToLower(strings.TrimRight(readagree, "\n")) == "yes" {
		paramUser.CoinId = strconv.Itoa(GetCoinIDHuobo(fiat)[fiat])
		//selected user
		var userpayment []string
		for {
			fmt.Println("Enter the balance (integer only) of the money(asset):")
			readbalance, err := n.ReadString('\n')
			if readbalance == "\n" {
				log.Println("You don't enter value")
				paramUser.Amount = ""
				break
			}
			if err != nil {
				log.Printf("Problem with enter the balance and error: %v", err)
			}
			readbalance = strings.TrimSpace(readbalance)
			var digitCheck = regexp.MustCompile(`^[0-9]+$`)
			if digitCheck.MatchString(readbalance) {
				paramUser.Amount = readbalance
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
			paramUser.IsMerchant = "true"
		}
		//BORDER
		fmt.Println("What is the minimum limit for the number of ads to set?")

		readborder, _ := n.ReadString('\n')
		readborder = strings.TrimSpace(readborder)
		for {
			if readborder != "" {
				var f int
				f, err = strconv.Atoi(readborder)
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

		gpm := getinfohuobi.GetPeymontMethodsHuobi(fiat)

		var availablepayment []string
		var choosePay []string
		mappayid := make(map[string]int)

		for i, j := range gpm {
			fmt.Println(i+1, " - ", j.TradeMethodName)
			availablepayment = append(availablepayment, strconv.Itoa(j.PayMethodId))
			mappayid[j.TradeMethodName] = j.PayMethodId
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
		log.Println(userpayment)
		for _, i := range userpayment {
			log.Println(i)
			choosePay = append(choosePay, i)
		}
		log.Println(choosePay)
		paramUser.PayMethod = unique_huobi(choosePay)
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
}

func unique_huobi(arr []string) string {
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

	if len(result) > 1 {
		result_str := strings.Join(result, ",")
		return result_str
	} else {
		return result[0]
	}
}

func GetCoinIDHuobo(fiat string) map[string]int {
	file_path := fmt.Sprintf("data/datahuobi/%s/%s_coinId.json", fiat, fiat)
	file, _ := os.Open(file_path)
	defer file.Close()
	decoder := json.NewDecoder(file)
	var data map[string]int
	decoder.Decode(&data)
	return (data)
}
