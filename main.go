package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gocolly/colly"
)

func ScrapeCountry(uri string) {

	defer func() {
		log.Println("Finished scraping " + uri)
	}()
	// Instantiate default collector
	c := colly.NewCollector()

	countryInformation := make(map[string]interface{})

	c.OnHTML("div.free-form-content__content", func(e *colly.HTMLElement) {
		sectionTitle := e.ChildText("h2")
		countryInformation[sectionTitle] = make(map[string]interface{})
		e.ForEach("div", func(_ int, element *colly.HTMLElement) {
			subsectionTitle := element.ChildText("h3")
			subsectionText := element.ChildText("p")

			if subsectionTitle != "" {
				countryInformation[sectionTitle].(map[string]interface{})[subsectionTitle] = subsectionText
			}
		})
	})

	c.Visit("https://www.cia.gov/the-world-factbook" + uri)

	content, err := json.Marshal(countryInformation)
	if err != nil {
		log.Println(err.Error())
	}

	countryName := uri[11 : len(uri)-1]

	os.WriteFile("results/"+countryName+".json", content, 0644)
}

func main() {
	res, err := http.Get("https://www.cia.gov/the-world-factbook/page-data/countries/page-data.json")
	if err != nil {
		log.Println(err.Error())
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(err.Error())
	}

	var jsonData interface{}
	json.Unmarshal(resBody, &jsonData)

	data := jsonData.(map[string]interface{})["result"].(map[string]interface{})["data"].(map[string]interface{})["countries"].(map[string]interface{})["edges"].([]interface{})

	for _, country := range data {
		countryUri := country.(map[string]interface{})["node"].(map[string]interface{})["uri"].(string)
		go ScrapeCountry(countryUri)
	}
	time.Sleep(10 * time.Second)
}
