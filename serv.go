package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/smtp"
	"os"
	"strconv"
	"strings"
	"time"
)

var ans string
var ConversionRate string
var XmrAmount string
var DollarAmount string
var WalletAdress string
var HashRate string

func printJstWalt(w http.ResponseWriter, r *http.Request) {
	getwalletinfo()
	pigstr := " cryptPiggy: " + XmrAmount + " Piggy: " + DollarAmount + " Crypt Rate: " + HashRate
	io.WriteString(w, pigstr)

}

//GET XMR TO USD CONVERSION VALUE
func getconversion() float64 {
	var num float64 = 0
	response, err := http.Get("https://min-api.cryptocompare.com/data/price?fsym=XMR&tsyms=USD")
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	} else {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Printf("%s", err)
			os.Exit(1)
		}
		var dat map[string]interface{}
		if err := json.Unmarshal(contents, &dat); err != nil {
			panic(err)
		}
		num = dat["USD"].(float64)
		return num
	}
	return num
}

//GET MINTER INFORMATION AND CALCULATE WALLET DOLLAR AMOUNT
func getwalletinfo() {
	searchterms := []string{"approx.speed", "Current balance"}
	dwarfpoolurl := "https://dwarfpool.com/xmr/address?wallet="
	dwarfpoolurl = dwarfpoolurl + WalletAdress
	num := getconversion()
	response, err := http.Get(dwarfpoolurl)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	} else {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Printf("%s", err)
			os.Exit(1)
		}
		data := strings.Split(string(contents), "\n")
		for iter, term := range searchterms {
			for i := 0; i < len(data); i++ {
				if strings.Contains(data[i], term) {
					line := strings.Split(data[i-1], " ")
					for _, val := range line {
						if strings.Contains(val, `badge-money">`) {
							ans := strings.TrimPrefix(val, `badge-money">`)
							i, err := strconv.ParseFloat(ans, 64)
							if err != nil {
								fmt.Println("Problem with float64 conversion")
							}
							if term == "approx.speed" {
								HashRate = ans
							} else if term == "Current balance" {
								ConversionRate = strconv.FormatFloat(num, 'f', 6, 64)
								XmrAmount = ans
								DollarAmount = strconv.FormatFloat((num * i), 'f', 6, 64)
							} else {
								fmt.Println("String Array not right at: ", iter)
							}
						}
					}
				}
			}
		}
	}
}

//SEND TROUBLE EMAIL
func emailalert() {
	fmt.Println("We have major problem")
	body := "Rig Seems To be having problems: " + HashRate
	from := os.Getenv("FROMEMAIL")
	pass := os.Getenv("PASSEMAIL")
	to := os.Getenv("TOEMAIL")
	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: Hello there\n\n" +
		body

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{to}, []byte(msg))

	if err != nil {
		fmt.Println("smtp error: ", err)
	}

}

func hashrwatcher() {
	c := time.Tick(60 * time.Second)
	for now := range c {
		fmt.Println("In hash Rate Watcher")
		getwalletinfo()
		time.Sleep(time.Second * 5)
		rate, err := strconv.ParseFloat(HashRate, 64)
		if err != nil {
			fmt.Println("Problem with string to float conversion at: ", now)
		}
		if rate <= 0.0 {
			fmt.Println("ALERT: EMAIL PROBLEM CURRENT RATE IS: ", rate)
			emailalert()
			time.Sleep(time.Hour * 4)
		}
	}
}

func main() {
	var servermode string
	if len(os.Args) != 3 {
		fmt.Println("Please supply 2 arguments")
		fmt.Println("./server wallet-address server-mode")
		fmt.Println("wallet-address is (your public XMR wallet address)")
		fmt.Println("server-mode is (stand-alone) or (web-enabled)")
		os.Exit(3)
	} else {
		WalletAdress = os.Args[1]
		servermode = os.Args[2]
	}
	fmt.Println("start")
	if servermode == "web-enabled" {
		go hashrwatcher()
		fmt.Println("web-enabled starting...")
		http.HandleFunc("/", printJstWalt)
		if err := http.ListenAndServe(":10037", nil); err != nil {
			panic(err)
		}
	} else if servermode == "stand-alone" {
		getwalletinfo()
		pigstr := " cryptPiggy: " + XmrAmount + " Piggy: " + DollarAmount + " Crypt Rate: " + HashRate
		fmt.Println(pigstr)
	} else {
		fmt.Println("You did not enter a correct server mode")
	}
}
