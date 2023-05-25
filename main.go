package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/gocolly/colly"
)

func main() {
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

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL.String())
	})

	c.Visit("https://www.cia.gov/the-world-factbook/countries/argentina")

	content, err := json.Marshal(countryInformation)
	if err != nil {
		log.Println(err.Error())
	}

	os.WriteFile("country.json", content, 0644)
}
