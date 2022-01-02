package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	iconv "github.com/djimenez/iconv-go"
	"gopkg.in/yaml.v3"
)

func scrapeRodina(c chan string) {
	req, err := http.NewRequest(http.MethodGet, "https://www.rodina.cz/scripts/jmena/default.asp?muz=0", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("User-Agent", "ScraperBot - We read list od names")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if terr := res.Body.Close(); terr != nil {
			log.Fatal(terr)
		}
	}()

	utfBody, err := iconv.NewReader(res.Body, "windows-1250", "utf-8")
	if err != nil {
		log.Fatal(err)
	}

	doc, err := goquery.NewDocumentFromReader(utfBody)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find(".jmena_vse h2").Each(func(i int, s *goquery.Selection) {
		c <- s.Find("a").Text()
	})
}

func scrapeCentrum(c chan string) {
	req, err := http.NewRequest(http.MethodGet, "http://svatky.centrum.cz/jmenny-seznam/?gender=1", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("User-Agent", "ScraperBot - We read list od names")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if terr := res.Body.Close(); terr != nil {
			log.Fatal(terr)
		}
	}()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find("#list-names .name").Each(func(i int, s *goquery.Selection) {
		name := s.Find("a").Text()
		if !strings.Contains(name, " ") {
			c <- name
		}
	})
}

func scraper(web int) <-chan string {
	c := make(chan string) // ubuffered channel

	// Fire a goroutine to send values on the channel
	go func() {
		switch web {
		case 0:
			scrapeRodina(c)
		case 1:
			scrapeCentrum(c)
		}
		// close the channel when done; otherwise it leaks resources
		close(c)
	}()

	return c
}

/* The fan in pattern is an important pattern which combines
mulitple channels, returns a single channel from those channels
*/
func fanIn(chans ...<-chan string) chan string {
	var wg sync.WaitGroup

	c := make(chan string)

	// Closure to send values from a channel
	output := func(ch <-chan string) {
		for n := range ch {
			c <- n
		}
		wg.Done()
	}

	wg.Add(len(chans))

	// send values on c via differnt goroutines
	for _, ch := range chans {
		go output(ch)
	}

	// wait for all goroutines to finish before closing the channel
	go func() {
		wg.Wait()
		close(c)
	}()

	return c
}

func runUpdate() {
	c1 := scraper(0)
	c2 := scraper(1)

	merged := fanIn(c1, c2)
	names := make([]string, 0)
	for n := range merged { // range loop terminates once the chan is closed, otherwise it blocks if there is no value
		names = append(names, n)
	}

	var data dataset
	data.Names = unique(names)

	yamlData, err := yaml.Marshal(&data)
	if err != nil {
		log.Fatal(err)
	}

	if err := os.MkdirAll(filepath.Dir(file), 0755); err != nil {
		log.Fatal(err)
	}

	if err := ioutil.WriteFile(file, yamlData, 0644); err != nil {
		log.Fatal(err)
	}
}

func unique(slice []string) []string {
	// create a map with all the values as key
	uniqMap := make(map[string]struct{})
	for _, v := range slice {
		uniqMap[v] = struct{}{}
	}

	// turn the map keys into a slice
	uniqSlice := make([]string, 0, len(uniqMap))
	for v := range uniqMap {
		uniqSlice = append(uniqSlice, v)
	}
	return uniqSlice
}
