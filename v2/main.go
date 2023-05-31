package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html"
	"log"
	"os"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
)

type Store struct {
	Address       string   `json:"address"`
	PhoneNumber   string   `json:"phoneNumber"`
	StoreName     string   `json:"storeName"`
	StoreFeatures []string `json:"storeFeatures"`
}

const path = "stores"

func main() {
	var linkVar string
	var headlessVar bool

	flag.StringVar(&linkVar, "link", "", "Starbucks Store Locator Link to scrape")
	flag.BoolVar(&headlessVar, "headless", true, "Choose whether to run this process using a headless browser; default is true")
	flag.Parse()

	if linkVar == "" {
		log.Fatal("No link included; Terminating early")
	}

	var page *rod.Page
	var stores []Store

	u := launcher.New().Headless(headlessVar).MustLaunch()
	page = rod.New().ControlURL(u).MustConnect().MustPage(linkVar)
	defer page.Close()

	log.Println("Starting scrape")

	locationList := page.MustElement("[data-e2e=locationList]")
	locations := locationList.MustElements("article")
	for _, location := range locations {
		infoButton := location.MustElement("[data-e2e=cardLink]")

		if err := infoButton.Click(proto.InputMouseButtonLeft, 1); err != nil {
			log.Fatal("Could not press a tag")
		}

		expandedLocationContent := page.MustElement("[data-e2e=expanded-location-content]")
		storeName := expandedLocationContent.MustElement("h2").MustText()
		address := expandedLocationContent.MustElement(".gridItem").MustText()
		phoneNumber := expandedLocationContent.MustElement("a").MustProperty("href").Str()
		f := expandedLocationContent.MustElement("[data-e2e=store-features]")
		features := f.MustElements("li")

		var storeFeatures []string

		for _, feature := range features {
			storeFeatures = append(storeFeatures, feature.MustText())
		}

		storeInfo := Store{
			Address:       html.UnescapeString(address),
			StoreName:     html.UnescapeString(storeName),
			PhoneNumber:   strings.TrimPrefix(phoneNumber, "tel:"),
			StoreFeatures: storeFeatures,
		}

		stores = append(stores, storeInfo)

		closeDetailsButton := page.MustElement("[data-e2e=overlay-close-button]")
		if err := closeDetailsButton.Click(proto.InputMouseButtonLeft, 1); err != nil {
			log.Fatal("Failed to close details modal")
		}
	}

	if err := createStoresDirectory(path); err != nil {
		log.Fatal("Something went really wrong")
	}

	fileName := fmt.Sprintf("%v/stores_%v.json", path, time.Now().Unix())
	file, _ := json.MarshalIndent(stores, "", " ")
	_ = os.WriteFile(fileName, file, 0644)

	log.Printf("Scrape finished: output saved to %v", fileName)
}

func createStoresDirectory(path string) error {
	log.Println("stores directory not found; creating it")
	const mode = 0755
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.Mkdir(path, mode)
		if err != nil {
			return err
		}
	}
	return nil
}
