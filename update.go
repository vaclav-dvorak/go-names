package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/PuerkitoBio/goquery"
	iconv "github.com/djimenez/iconv-go"
	"gopkg.in/yaml.v3"
)

func scrapeRodinaCz() []string {
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

	ret := make([]string, 0)
	doc.Find(".jmena_vse h2").Each(func(i int, s *goquery.Selection) {
		ret = append(ret, s.Find("a").Text())
	})
	return ret
}

func runUpdate() {
	var data dataset
	namesRodina := scrapeRodinaCz()
	namesNeco := []string{"Test"}
	data.Names = append(namesRodina, namesNeco...)
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
