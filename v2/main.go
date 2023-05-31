package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

type Store struct {
	Address     string
	PhoneNumber string
	StoreName   string
}

const path = "stores"

func main() {
	var linkVar string

	flag.StringVar(&linkVar, "link", "", "Starbucks Store Locator Link to scrape")
	flag.Parse()

	if linkVar == "" {
		log.Fatal("No link included; Terminating early")
	}

	var stores []Store

	// u := launcher.New().Headless(false).MustLaunch()

	// page := rod.New().ControlURL(u).MustConnect().MustPage(linkVar)
	page := rod.New().MustConnect().MustPage(linkVar)
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

		storeInfo := Store{
			Address:     address,
			StoreName:   storeName,
			PhoneNumber: strings.TrimPrefix(phoneNumber, "tel:"),
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

	// time.Sleep(1 * time.Hour)
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
