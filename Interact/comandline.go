// choose available payment methods
package Interact

import (
	"bufio"
	"fmt"
	"github.com/Zmey56/arbitrage/getdata"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Parameters struct {
	PayTypes      []string `json:"payTypes"`
	TransAmount   string   `json:"transAmount"`
	PublisherType string   `json:"publisher_type"`
}

func InputCommandLine(fiat string) Parameters {
	paramUser := Parameters{}

	//selected user
	var userpayment []string
	n := bufio.NewReader(os.Stdout)
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

	fmt.Println("Do you want to choose only the merchant? If you want - enter Yes")

	readmerchant, err := n.ReadString('\n')
	if err != nil {
		log.Println("User try to enter wong value:", err)
	}
	if strings.ToLower(readmerchant) == "yes" {
		paramUser.PublisherType = "merchant"
	}

	gpm := getdata.GetPaymentFromJSON(fiat)

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
