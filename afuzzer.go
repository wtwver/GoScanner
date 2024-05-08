// TODO
// 1. Fix http not workoing
// 2. Add flag
// 3. Add database
// 4. mutiple wordlist

package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/queue"
)

func AppendIfMissing(slice []string, i string) []string {
	i = strings.TrimSpace(i)
	for _, ele := range slice {
		if len(i) == 1 || i == "" || ele == i {
			return slice
		}
	}
	return append(slice, i)
}

func learn(domain string) []int {
	c := colly.NewCollector(
	// colly.UserAgent(byte "Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:47.0) Gecko/20100101 Firefox/47.3")
	)
	// c.SetProxy("http://127.0.0.1:1000/")
	// extensions.RandomUserAgent(c)

	output := []int{}

	c.OnResponse(func(r *colly.Response) {
		// fmt.Println(string(r.Body))
		output = append(output, r.StatusCode)
		output = append(output, len(r.Body))

	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println(string(r.Body))
		output = append(output, r.StatusCode)
		output = append(output, len(r.Body))

	})

	c.Visit(fmt.Sprintf("http://%s/", domain))
	time.Sleep(5 * time.Second)
	c.Visit(fmt.Sprintf("http://%s/notexisted-asd15sad9561dsa.html", domain))
	c.Visit(fmt.Sprintf("http://%s/notedsa8re7w", domain))

	fmt.Println(output)

	return output
}

func resposneHandler(r *colly.Response, boring []int, found []string, q *queue.Queue, dictionary []string) []string {
	if !(r.StatusCode == boring[0] && len(r.Body) == boring[1] || r.StatusCode == boring[2] && len(r.Body) == boring[3]) {
		link := r.Request.URL.String()
		fmt.Print(time.Now().Format("01-02-2006 15:04:05.000000"), " ")
		fmt.Println(r.StatusCode, len(r.Body), r.Request.URL)
		if r.StatusCode == 403 || r.StatusCode == 200 {
			found = append(found, link)

			// if endpoint ends with words => recursion
			re, _ := regexp.Compile(`/[\w]*$`)
			if re.MatchString(link) {
				for _, j := range dictionary {
					url := fmt.Sprintf("%s/%s", link, j)
					q.AddURL(url)
				}
			}
		}
	}
	return found
}

func recusrion(urls []string) []string { // normalize dir
	// match file:  re, _ = regexp.Compile("\\.[a-z]{2,}$")
	// mat non extension endpoint:  re, _ = regexp.Compile("/.+$")
	re, _ := regexp.Compile(`/[\w]*$`) //match dir
	dir := []string{}

	for _, ele := range urls {
		if re.MatchString(ele) && strings.HasSuffix(ele, "/") {
			// fmt.Println("1", ele)/
			dir = append(dir, ele)
		} else if re.MatchString(ele) {
			// fmt.Println("2", ele+"/")
			dir = append(dir, ele+"/")
		}

	}
	return dir
}

func main() {

	dictionary := []string{}
	found := []string{}

	list := "./common.txt"
	domain := "testphp.vulnweb.com"
	ext := []string{"php", "js", "html"}
	boring := learn(domain)

	q, _ := queue.New(50, &queue.InMemoryQueueStorage{MaxSize: 10000000}) // thread ,  queue storage
	// fmt.Println(reflect.TypeOf(q))

	c := colly.NewCollector()
	// c.SetProxy("http://127.0.0.1:1000/")

	c.OnResponse(func(r *colly.Response) {
		found = resposneHandler(r, boring, found, q, dictionary)
	})

	c.OnError(func(r *colly.Response, _ error) {
		found = resposneHandler(r, boring, found, q, dictionary)
	})

	file, _ := os.Open(list)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		dictionary = AppendIfMissing(dictionary, scanner.Text())

		re, _ := regexp.Compile("\\.[a-z]{2,}$") // append extension if none
		if !re.MatchString(scanner.Text()) {
			for _, i := range ext {
				dictionary = AppendIfMissing(dictionary, scanner.Text()+"."+i)
			}
		}
	}

	for _, ele := range dictionary {
		url := fmt.Sprintf("http://%s/%s", domain, ele)
		q.AddURL(url)
	}
	q.Run(c)

	fmt.Println("found ===", found)
}
