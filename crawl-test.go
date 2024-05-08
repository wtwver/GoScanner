package main

import (
	"fmt"

	"github.com/gocolly/colly"
	"mvdan.cc/xurls/v2"
)

func main() {
	c := colly.NewCollector(
		colly.AllowedDomains("juice-shop-staging.herokuapp.com"),
	)

	c.OnHTML("script[src]", func(e *colly.HTMLElement) {
		link := e.Attr("src")
		fmt.Printf("Link found: %q -> %s\n", e.Text, link)
		c.Visit(e.Request.AbsoluteURL(link))
	})

	c.OnResponse(func(r *colly.Response) {
		// fmt.Println("Response", r.StatusCode)

		// rx1 := xurls.Strict() //extract strict urls
		// a := rx1.FindAllString(string(r.Body), -1)
		// fmt.Println(a)

		rx2, _ := xurls.StrictMatchingScheme("/") //extract relative urls
		b := rx2.FindAllString(string(r.Body), -1)

		for _, txt := range b {
			fmt.Println(txt)
			// add some magic function to clean the output
			// e.g. split on " ?
		}

	})

	// Start scraping on https://hackerspaces.org
	c.Visit("https://juice-shop-staging.herokuapp.com/")
}
