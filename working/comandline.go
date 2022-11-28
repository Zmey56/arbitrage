// choose available payment methods
package working

import (
	"bufio"
	"fmt"
	"github.com/Zmey56/arbitrage/getdata"
	"log"
	"os"
	"strconv"
	"strings"
)

func InputCommandLine(fiat string) ([]string, float64) {

	//selected user
	var userpayment []string
	n := bufio.NewReader(os.Stdout)
	fmt.Println("Enter the balance of the money(asset):")
	readbalance, _ := n.ReadString('\n')
	readbalance = strings.TrimSpace(readbalance)
	transAmount, err := strconv.ParseFloat(readbalance, 64)
	if err != nil {
		log.Println(err)
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
	userpayment = unique(userpayment)
	return userpayment, transAmount
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
