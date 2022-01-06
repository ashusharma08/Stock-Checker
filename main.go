package main

import (
	"fmt"
	"time"
)

const amazonURL = "https://www.amazon.in/Xbox-Series-X/dp/B08J7QX1N1"
const flipkartURL = "https://www.flipkart.com/microsoft-xbox-series-x-1024-gb/p/itm63ff9bd504f27"

func main() {
	client := GetNewClient()

	bot, err := getBot()
	if err != nil {
		return
	}
	vl := getVendorList(client)
	lastRun := time.Now()
	for true {
		available := ""
		for _, item := range vl {
			a, err := item.CheckForAvailability()
			if err != nil {
				sendMsg(err.Error(), bot)
				continue
			}
			available += a
		}
		if len(available) == 0 {
			if time.Since(lastRun.Add(time.Hour*1)) > 0 {
				sendMsg(fmt.Sprintf("bot is still up and stocks are still down. %s \n %s", amazonURL, flipkartURL), bot)
				lastRun = time.Now()
			}
			time.Sleep(1 * time.Minute)
		} else {
			sendMsg(available, bot)
		}
	}
}

func getVendorList(client *client) []Vendor {
	v := make([]Vendor, 0)
	v = append(v, GetAmazon(amazonURL, client.client))
	v = append(v, GetFlipkart(flipkartURL, client.client))
	return v
}
