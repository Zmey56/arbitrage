package workingokx

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/Zmey56/arbitrage/pkg/getinfookx"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func InputCommandLineOKX(fiat string) {
	file_path := "cmd/enterparam/paramokx/" + fiat + ".json"
	paramUser := getinfookx.ParametersOKX{}
	n := bufio.NewReader(os.Stdout)
	fmt.Println("Do you want to enter new params for OKX? If you want - enter Yes")
	readagree, err := n.ReadString('\n')
	if err != nil {
		log.Println("User try to enter wrong value:", err)
	}
	if strings.ToLower(strings.TrimRight(readagree, "\n")) == "yes" {
		paramUser.CoinId = fiat
		//selected user
		var userpayment = "all"
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

		for {
			readpercent, _ := n.ReadString('\n')
			readpercent = strings.TrimSpace(readpercent)
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

		gpm := getinfookx.GetPeymontMethodsOKX(fiat)

		var availablepayment []string
		mappayid := make(map[string]int)

		for i, j := range gpm {
			fmt.Println(i+1, " - ", j.TradeMethodName)
			availablepayment = append(availablepayment, j.TradeMethodName)
			mappayid[j.TradeMethodName] = j.PayMethodId
		}
		//log.Println(availablepayment)

		fmt.Println("Available means of payment, choose number:")
		for {
			readnumber, _ := n.ReadString('\n')
			readnumber = strings.TrimSpace(readnumber)
			number, err := strconv.Atoi(readnumber)
			if err != nil {
				log.Println(err)
			}

			if number >= len(availablepayment) || number < 1 {
				fmt.Println("You have selected the wrong number")
				continue
			} else {
				fmt.Println(availablepayment[number-1])
				userpayment = strings.ReplaceAll(availablepayment[number-1], " ", "+")
				break
			}
		}
		//log.Println(userpayment)
		paramUser.PayMethod = userpayment
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
