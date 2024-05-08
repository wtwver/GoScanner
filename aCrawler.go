// TODO

// handle onError: fuzz on 403 500
// handle http/https

// 2. Add recursive crawl
// Cannot handle subdomains
// Cannot handle /index //index unique
// Check repeated endpoints

// handle no response request

package main

import (
	"bufio"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
	// "github.com/gocolly/colly/proxy"
	// "github.com/gocolly/redisstorage"
	"github.com/gocolly/redisstorage"

	"mvdan.cc/xurls/v2"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gocolly/colly/extensions"
	"github.com/gocolly/colly/queue"
)

func getUniqueWords(input string) []string {
	re, _ := regexp.Compile("[a-zA-Z0-9_\\-\\.]+")
	uni_words := []string{}

	for _, word := range re.FindAllString(input, -1) {
		if strings.Index(word, ".") > -1 {
			for _, long_word := range strings.Split(word, ".") {
				// fmt.Println(long_word)
				uni_words = AppendIfMissing(uni_words, long_word)
			}
		} else {
			uni_words = AppendIfMissing(uni_words, word)
		}
	}

	return uni_words
}

func getBlacklist(url string) []string {
	c := colly.NewCollector()
	words := []string{}

	c.OnResponse(func(r *colly.Response) {
		words = getUniqueWords(string(r.Body))
	})

	c.Visit(url)

	return words
}

func AppendIfMissing(slice []string, i string) []string {
	i = strings.TrimSpace(i)
	for _, ele := range slice {
		if len(i) == 1 || i == "" || ele == i {
			return slice
		}
	}
	return append(slice, i)
}

func in(str string, slice []string) bool {
	for _, ele := range slice {
		if str == ele {
			return true
		}
	}
	return false
}

func collectURL(resp string, absolute []string, relative []string) ([]string, []string) {
	// absolute := []string{}
	// relative := []string{}

	rx1 := xurls.Strict() //extract strict urls

	a := rx1.FindAllString(resp, -1)
	for _, txt := range a {
		txt := strings.Split(txt, "\"")[0]
		absolute = AppendIfMissing(absolute, txt)
	}

	rx2, _ := xurls.StrictMatchingScheme("/") //extract relative urls

	b := rx2.FindAllString(resp, -1)
	for _, txt := range b {
		txt := strings.Split(txt, "\"")[0]
		relative = AppendIfMissing(relative, txt)
	}

	// fmt.Println(absolute)
	return absolute, relative

}

func collectUniqueWords(resp string, blacklist_on bool, blacklist []string, wordlist []string) []string {
	if !blacklist_on {
		for _, i := range getUniqueWords(resp) {
			wordlist = append(wordlist, i)
		}
	} else {
		for _, i := range getUniqueWords(resp) {
			if !in(i, blacklist) {
				wordlist = append(wordlist, i) //getUniqueWords ensured Unique
			}
		}
	}
	return wordlist

}

func insert(db *sql.DB, time string, respCode string, respLen string, url string) {
	statement, _ := db.Prepare("INSERT INTO log (time, respCode, respLen, url) VALUES (?, ?, ?, ?)")
	statement.Exec(time, respCode, respLen, url)
}

func learn(domain string) []int {
	c := colly.NewCollector()
	c.WithTransport(&http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	})
	// c.SetProxy("http://127.0.0.1:1000/")
	extensions.RandomUserAgent(c)

	if len(*proxy) > 0 {
		c.SetProxy("http://" + *proxy)
	}
	output := []int{}

	c.OnResponse(func(r *colly.Response) {
		output = append(output, r.StatusCode)
		output = append(output, len(r.Body))
	})

	c.OnError(func(r *colly.Response, err error) {
		output = append(output, r.StatusCode)
		output = append(output, len(r.Body))
	})

	c.Visit(fmt.Sprintf("http://%s/", domain))
	c.Visit(fmt.Sprintf("http://%s/notexisted-asd15sad9561dsa.html", domain))
	time.Sleep(3 * time.Second)
	c.Visit(fmt.Sprintf("http://%s/notedsa8re7w", domain))
	// c.Visit(fmt.Sprintf("https://%s/", domain))
	// c.Visit(fmt.Sprintf("https://%s/notexisted-asd15sad9561dsa.html", domain))
	// time.Sleep(3 * time.Second)
	// c.Visit(fmt.Sprintf("https://%s/notedsa8re7w", domain))

	return output
}

func bored(r *colly.Response, boring []int) bool {
	if !(r.StatusCode == boring[0] && len(r.Body) == boring[1] || r.StatusCode == boring[2] && len(r.Body) == boring[3]) {
		return false
	}
	return true
}

//print for debug message
func print(str string) {
	if !*plain {
		fmt.Println(str)
	}
}

func printsl(str string) {
	if !*plain {
		fmt.Print(str)
	}
}

var plain *bool
var randomagent *bool
var proxy *string

func main() {
	var urlSeed string
	var httpurl string
	flag.StringVar(&urlSeed, "u", "demo.testfire.net", "Enter a url seed")
	flag.StringVar(&httpurl, "url", "http://demo.testfire.net", "Enter a full url seed with protocol")
	scope := flag.Bool("s", true, "In scope domain only")
	plain = flag.Bool("p", false, "plain output")
	randomagent = flag.Bool("ra", false, "enable random agent ")
	blacklist_on := flag.Bool("b", true, "Collect and filter blacklist")
	output := flag.String("o", "txt", "Print output")
	proxy = flag.String("pr", "", "Proxy and port in format 127.0.0.1:8080")
	dict := flag.String("w", "", " Supply worlist to enable fuzzing mode")
	pcrawler := flag.String("pc", "", " Supply pcrawler result to check active")
	redis := flag.String("r", "", " Redis queue server e.g. 206.189.88.45:6379")
	role := flag.String("role", "master", "default master, enter worker to enable worker mode")
	depth := flag.Int("d", 1, "crawling depth")
	thread := flag.Int("t", 30, "number of thread")
	flag.Parse()

	if !*plain {
		print("[*] Running on " + *role + " mode")
		print("[*] Running aCrawler on " + urlSeed)
		print("[*] Scope only: " + fmt.Sprint(*scope))
		print("[*] Crawling Depth: " + fmt.Sprint(*depth))
		print("[*] Thread: " + fmt.Sprint(*thread))
		print("[*] Proxy: " + fmt.Sprint(*proxy))

		// fmt.Println("[*] Running on " + *role + " mode")
		// fmt.Println("[*] Running aCrawler on " + urlSeed)
		// fmt.Println("[*] Scope only: " + fmt.Sprint(*scope))
		// fmt.Println("[*] Crawling Depth: " + fmt.Sprint(*depth))
		// fmt.Println("[*] Thread: " + fmt.Sprint(*thread))
	}

	// statement, _ := db.Prepare("CREATE TABLE IF NOT EXISTS log (time TEXT , respCode TEXT, respLen TEXT, url TEXT)") //Fix me PRIMARY KEY
	// statement.Exec()

	storage := &redisstorage.Storage{
		Address: *redis,
		// Address:  "206.189.88.45:6379",
		Password: "",
		DB:       0,
		Prefix:   "colly_q1",
	}

	q, err := queue.New(*thread, &queue.InMemoryQueueStorage{MaxSize: 10000000})
	if len(*redis) > 0 {
		q, err = queue.New(*thread, storage)
		if q == nil {
			log.Fatal(err)
		}
	}

	print("queue size:" + fmt.Sprint(q.Size()))

	wordlist := []string{}
	absolute := []string{}
	relative := []string{}
	active := []string{}
	dictionary := []string{}
	boring := learn(urlSeed)
	for boring[0] == 0 {
		boring = learn(urlSeed)
		print(fmt.Sprint(boring))
	}
	blacklist := getBlacklist("https://tools.ietf.org/html/rfc1866")
	ext := []string{"php", "js", "html"}

	if len(*dict) > 0 {
		file, err := os.Open(*dict)
		if err != nil {
			print(fmt.Sprint("[*] Your worlist ", *dict, "is not valid. Using default wordlist instead."))
			file, _ = os.Open("common.txt")
		}
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

		// fmt.Println("finnished intialize dictionary")
		for _, ele := range dictionary {
			go func(ele string) {
				url := fmt.Sprintf("http://%s/%s", urlSeed, ele)
				q.AddURL(url)
			}(ele)
		}
		// fmt.Println("finnished queue")

	}

	if len(*pcrawler) > 0 {
		file, err := os.Open(*pcrawler)
		if err != nil {
			print(fmt.Sprint("[*] Your pcrawl result ", *pcrawler, "is not valid. Ignored"))
		} else {
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				// fmt.Println("read", scanner.Text())
				q.AddURL(scanner.Text())
			}
		}
	}

	c := colly.NewCollector()
	c.WithTransport(&http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	})

	if *scope {
		c.AllowedDomains = []string{urlSeed}
	}
	c.MaxDepth = *depth
	if *randomagent {
		extensions.RandomUserAgent(c)
	}

	if len(*proxy) > 0 {
		c.SetProxy("http://" + *proxy)
	}
	// rp, err := proxy.RoundRobinProxySwitcher("http://206.189.88.45:8899", "http://143.198.196.17:8899", "http://143.198.200.57:8899")
	// if err != nil {
	// 	fmt.Println("========Proxy err. Check your proxy")
	// 	log.Fatal(err)
	// }
	// c.SetProxyFunc(rp)

	c.OnResponse(func(r *colly.Response) {
		// if !(r.StatusCode == boring[0] && len(r.Body) == boring[1] || r.StatusCode == boring[2] && len(r.Body) == boring[3]) {
		if !bored(r, boring) {
			link := r.Request.URL.String()

			wordlist = collectUniqueWords(string(r.Body), *blacklist_on, blacklist, wordlist)

			absolute, relative = collectURL(string(r.Body), absolute, relative)
			for _, ele := range relative {
				q.AddURL("http://" + urlSeed + ele)
			}
			for _, ele := range absolute {
				q.AddURL(ele)
			}

			printsl(fmt.Sprint(r.StatusCode, len(r.Body), " "))
			fmt.Println(r.Request.URL.String())
			active = AppendIfMissing(active, r.Request.URL.String())
			// insert(db, time.Now().Format("01-02-2006 15:04:05.000000"), strconv.Itoa(r.StatusCode), strconv.Itoa(len(r.Body)), r.Request.URL.String())

			if len(*dict) > 0 { //Enter fuzz mode IF wordlist provide
				re, _ := regexp.Compile(`/[\w]*$`)
				if re.MatchString(link) {
					for _, j := range dictionary {
						url := fmt.Sprintf("%s/%s", link, j)
						q.AddURL(url)
					}
				}
			}
		}
		// fmt.Println("queue size:" + fmt.Sprint(q.Size()))

	})

	c.OnHTML("script[src]", func(e *colly.HTMLElement) {
		link := e.Attr("src")
		c.Visit(e.Request.AbsoluteURL(link))
	})

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		c.Visit(e.Request.AbsoluteURL(link))
	})

	c.OnError(func(r *colly.Response, err error) {
		print(fmt.Sprint(strconv.Itoa(r.StatusCode), " ", strconv.Itoa(len(r.Body)), " ", r.Request.URL.String()))
	})

	// Add urlseed to queue only when worker mode is off
	if *role == "master" {
		// q.AddURL("https://" + urlSeed)
		// q.AddURL("https://" + urlSeed + "/robots.txt")
		q.AddURL("http://" + urlSeed)
		q.AddURL("http://" + urlSeed + "/robots.txt")
		q.AddURL("http://" + urlSeed + "/wavsep/active/index-xss.jsp")
	}

	// ['sitemap.xml', 'robots.txt', 'crossdomain.xml', 'clientaccesspolicy.xml', 'phpmyadmin', 'pma', 'myadmin',
	//  '.svn', '.ssh', '.git', 'CVS', 'info.php', 'phpinfo.php', 'test.php', 'php.php', 'Thumbs.db', 'CHANGELOG',
	//  '.DS_Store', 'composer.lock', 'composer.json', '.hg', '.hgignore', '.gitignore', 'access.log', '.bash_history',
	//  '.bash_profile', '.htaccess', '.htpasswd', '.mysql_history', '.passwd', '.htconfig', '.htusers'])

	q.Run(c)

	if *output == "txt" {
		// fmt.Println(relative)
		// for _, ele := range absolute {
		// 	fmt.Println(ele)
		// }
		// // fmt.Println(wordlist)

		print(fmt.Sprint("[*] relative path", len(relative)))
		print(fmt.Sprint("[*] absolute path", len(absolute)))
		print(fmt.Sprint("[*] unqiue wordlist", len(wordlist)))
		print(fmt.Sprint("[*] blacklist", len(blacklist)))
	} else if *output == "json" {
		j1, _ := json.Marshal(relative)
		print(string(j1))
	}

}
