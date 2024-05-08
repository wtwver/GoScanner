package main

import (
	"fmt"
	"strconv"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/queue"
	"github.com/gocolly/redisstorage"
)

func main() {
	c := colly.NewCollector()

	storage := &redisstorage.Storage{
		Address: "206.189.88.45:6379", //change private ip address
		// Address:  "10.104.0.2:6379", //change private ip address
		Password: "",
		DB:       0,
		Prefix:   "colly_q1",
	}

	// err := storage.Clear()
	// fmt.Println(err)

	q, _ := queue.New(100, storage)

	c.OnResponse(func(r *colly.Response) {
		fmt.Println(strconv.Itoa(r.StatusCode), strconv.Itoa(len(r.Body)), r.Request.URL.String())
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println(strconv.Itoa(r.StatusCode), strconv.Itoa(len(r.Body)), r.Request.URL.String())
	})

	fmt.Println(q.Size())

	q.Run(c)
}
