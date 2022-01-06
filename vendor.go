package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Vendor interface {
	CheckForAvailability() (string, error)
}

type Flipkart struct {
	client *http.Client
	req    *http.Request
	url    string
}

func GetFlipkart(url string, c *http.Client) *Flipkart {
	return &Flipkart{
		url:    url,
		req:    getRequest(url),
		client: c,
	}
}

func getRequest(url string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.90 Safari/537.36")
	return req
}

func (f *Flipkart) CheckForAvailability() (string, error) {
	res, err := f.client.Do(f.req)
	if err != nil {
		return "", fmt.Errorf("FLIPKART error in do %#v", err.Error())
	}
	if res != nil && res.Body != nil {
		defer res.Body.Close() // close it on defer
	}
	if res.StatusCode == http.StatusOK {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return "", fmt.Errorf("FLIPKART error decoding respnse %#v", err)
		}

		doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
		if err != nil {
			return "", fmt.Errorf("FLIPKART error parsing html %#v", err.Error())
		}
		notifyel := doc.Find("button._2uS5ZX")
		buynowbuttonel := doc.Find("button.ihZ75k")
		if (notifyel.Nodes != nil && len(notifyel.Nodes) > 0) && (buynowbuttonel.Nodes == nil) {
			//do this chutiyapa cuz html is not rendered properly
			button := notifyel.Nodes[0].FirstChild.Data
			if strings.Contains(button, "NOTIFY ME") {
				fmt.Printf("FLIPKART: time: %s still not in stock %#v", time.Now().Local(), button)
				fmt.Println()
				return "", nil
			}
		} else if notifyel.Nodes == nil && (buynowbuttonel.Nodes != nil && len(buynowbuttonel.Nodes) > 0) {
			button := buynowbuttonel.Nodes[0].LastChild.Data
			if strings.Contains(button, "BUY NOW") {
				fmt.Printf("FLIPKART: time: %s In stock %#v", time.Now().Local(), button)
				fmt.Println()
				return fmt.Sprintf("FLIPKART: In stock XBOX at: %s", flipkartURL), nil
			}
		}
	}
	return "", nil
}

type Amazon struct {
	client *http.Client
	url    string
	req    *http.Request
}

func GetAmazon(url string, c *http.Client) *Amazon {
	return &Amazon{
		url:    url,
		req:    getRequest(url),
		client: c,
	}
}

func (a *Amazon) CheckForAvailability() (string, error) {
	res, err := a.client.Do(a.req)
	if err != nil {
		return "", fmt.Errorf("AMAZON error in do %#v", err.Error())
	}
	if res != nil && res.Body != nil {
		defer res.Body.Close() // close it on defer
	}
	if res.StatusCode == http.StatusOK {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return "", fmt.Errorf("AMAZON error decoding respnse %#v", err.Error())
		}
		doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
		if err != nil {
			return "", fmt.Errorf("AMAZON error parsing html %#v", err.Error())
		}
		el := doc.Find("div#availability")
		if el != nil && el.Nodes != nil && len(el.Nodes) > 0 {
			//do this chutiyapa cuz html is not rendered properly
			spanEl := el.Nodes[0].FirstChild.NextSibling.FirstChild.Data
			if strings.Contains(spanEl, "Currently unavailable.") {
				fmt.Printf("AMAZON: time: %s still not in stock %#v", time.Now().Local(), spanEl)
				fmt.Println()
				return "", nil
			} else {
				//send notification
				fmt.Printf("AMAZON: time: %s In stock %#v", time.Now().Local(), spanEl)
				fmt.Println()
				return fmt.Sprintf("AMAZON: In stock XBOX at: %s", amazonURL), nil
			}
		}
	}
	return "", nil
}
