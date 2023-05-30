package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/go-rod/rod"
)

type Store struct {
	Address           string `json:"address"`
	StoreLocationLink string `json:"storeLocationLink"`
	PhoneNumber       string `json:"phoneNumber"`
}

func getStores(link string) map[string]Store {
	page := rod.New().MustConnect().MustPage(link)
	defer page.MustClose()

	locationList := page.MustElement(".base___3sH_T")
	locations := locationList.MustElements(".base___3LiS9")

	stores := make(map[string]Store)
	for _, location := range locations[:30] {
		address := location.MustElement("[data-e2e=address]").MustText()
		storelink := location.MustElement("[data-e2e=cardLink]").MustProperty("href")

		stores[address] = Store{
			StoreLocationLink: storelink.Str(),
			Address:           address,
		}

	}
	return stores
}

func main() {
	var linkVar string

	flag.StringVar(&linkVar, "link", "link", "link used to scrape for stores")
	flag.Parse()

	if linkVar == "" {
		log.Fatal("Could not start scrape without link")
	}
	browser := rod.New().MustConnect()
	defer browser.MustClose()

	pool := rod.NewPagePool(4)

	create := func() *rod.Page {
		return browser.MustIncognito().MustPage()
	}

	log.Println("Starting scrape")
	stores := getStores(linkVar)

	getPhoneNumber := func(store Store) Store {
		page := pool.Get(create)
		defer pool.Put(page)

		page.Timeout(5 * time.Second).MustNavigate(store.StoreLocationLink).MustWaitOpen()

		storeInformation := page.MustElement("[aria-labelledby=expandedLocationCardLabel]")
		contactInformation := storeInformation.MustElement(".whiteSpace-noWrap")
		phoneNumber := contactInformation.MustElement("a").MustProperty("href").Str()
		phoneNumber = strings.TrimPrefix(phoneNumber, "tel:")
		store.PhoneNumber = phoneNumber
		return store
	}

	wg := sync.WaitGroup{}
	wg.Add(len(stores))
	for k, v := range stores {
		go func(v Store, k string) {
			defer wg.Done()
			v = getPhoneNumber(v)
			stores[k] = v
		}(v, k)
	}

	wg.Wait()
	pool.Cleanup(func(p *rod.Page) {
		p.MustClose()
	})

	var storesArray []Store

	for _, store := range stores {
		storesArray = append(storesArray, store)
	}

	fileName := fmt.Sprintf("stores/stores_%v.json", time.Now().Unix())
	file, _ := json.MarshalIndent(storesArray, "", " ")
	_ = os.WriteFile(fileName, file, 0644)
	log.Printf("Scrape finished wrote file to stores/%v", fileName)
}
