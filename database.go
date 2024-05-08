package main

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type Record struct {
	time     string `json:"time"`
	respCode string `json:"respCode"`
	respLen  string `json:"respLen"`
	url      string `json:"url"`
}

func getTaskid() int {
	taskID := rand.Intn(10000)
	for tasks[taskID] != "" {
		fmt.Println("regenerating taskid")
		taskID = rand.Intn(10000)
	}
	// tasks[taskID] = "running"
	return taskID
}

func cmd(str string) (string, error) {
	cmd := exec.Command("cmd", "/C", str)
	out, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Println("linux detected")
		cmd = exec.Command("bash", "-c", str)
		out, err = cmd.CombinedOutput()
	}

	return string(out), err
}
func cmd_debug(str string) (string, error) {
	cmd := exec.Command("cmd", "/C", str)
	out, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Println("linux detected")
		cmd = exec.Command("bash", "-c", str)
		out, err = cmd.CombinedOutput()
	}

	return string(out), err
}

//real time pipe output
func cmd2(id int, str string) {
	cmd := exec.Command("cmd", "/C", str)
	_, err := cmd.CombinedOutput()

	var unix bool
	if err != nil {
		unix = true
	}

	cmd = exec.Command("cmd", "/C", str)
	stdout, err := cmd.StdoutPipe()
	if unix {
		cmd = exec.Command("bash", "-c", str)
		stdout, err = cmd.StdoutPipe()
	}

	cmd.Start()

	scanner := bufio.NewScanner(stdout)
	// scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		m := scanner.Text()
		tasks[id] += "\n" + m
	}
	cmd.Wait()
	tasks[id] += "\nThis is the last line"
	_, err = cmd_debug(str)
	if err != nil {
		tasks[id] += fmt.Sprint(err)
	}

}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func interface2slice(a interface{}) []string {
	var slice []string
	for _, ele := range a.([]interface{}) {
		slice = append(slice, ele.(string))
	}
	return slice
}

func toolsHandler(urlSeed string, b interface{}) []string {
	fmt.Println(b)
	var tasksIDs []string

	for _, tool := range interface2slice(b) {
		wg.Add(1)
		go func(tool string) {
			defer wg.Done()
			fmt.Println("running", tool)

			taskID := getTaskid()
			fmt.Println(taskID)
			tasksIDs = append(tasksIDs, strconv.Itoa(taskID))

			// out, err := cmd("./sEngine.sh " + tool + " " + urlSeed)
			// tasks[taskID] = out + "\nThis is the last line" + fmt.Sprint(err)
			// fmt.Println(string(out))
			// fmt.Println(err)

			cmd2(taskID, "make "+tool+" target="+urlSeed)
		}(tool)
	}
	for len(tasksIDs) != len(b.([]interface{})) {
		time.Sleep(1)
	}
	return tasksIDs
	// return []int{1111, 2222}
}

var wg sync.WaitGroup

var tasks map[int]string
var tools = make(map[string]interface{})
var tools2 = make(map[string][]string)

func main() {
	tools["pcrawler"] = []string{"self", "waybackurls", "gau"}
	tools["acrawler"] = []string{"self", "hakrawler"}
	tools["afuzzer"] = []string{"self", "ffuf"}
	tools["scanner"] = []string{"nuclei", "dalfox", "whatweb"}

	tools2["pcrawler"] = []string{"self_pcrawler", "waybackurls", "gau"}
	tools2["acrawler"] = []string{"self_acrawler", "hakrawler", "gospider"}
	tools2["afuzzer"] = []string{"self_afuzzer", "ffuf"}
	tools2["scanner"] = []string{"nuclei", "dalfox", "whatweb", "zap", "nikto", "inql", "iis_shortname_scanner", "wpscan", "bfac", "headi", "arjun", "xsstrike", "smuggler", "intrigue-ident", "wafw00f", "favfreak"}

	tasks = make(map[int]string)
	// fmt.Println("%T\n", tasks)

	// db, err := sql.Open("mysql", "root:password@123@tcp(database:3306)/test")
	db, err := sql.Open("mysql", "root:password@123@tcp(206.189.88.45:1000)/test")
	err = db.Ping()
	if err != nil {
		db, err = sql.Open("mysql", "root:password@123@tcp(127.0.0.1:1000)/test")
	}
	dbcon := false
	if err != nil {
		dbcon = true
		statement, _ := db.Prepare("CREATE TABLE IF NOT EXISTS log (time TEXT , respCode TEXT, respLen TEXT, url TEXT)") //Fix me PRIMARY KEY
		statement.Exec()
	}
	// statement, _ = db.Prepare("DELETE FROM log")
	// statement.Exec()

	// statement, err = db.Prepare("INSERT INTO log (time, respCode, respLen, url) VALUES (?, ?, ?, ?)")
	// statement.Exec("02-05-2021 04:11:31.709543", "200", "1502", "http://www.google.com")

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
	},
	))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})
	r.Static("/asset", "./frontend/asset")
	r.Use(static.Serve("/", static.LocalFile("./frontend/", false)))

	r.GET("/api/gitpull", func(c *gin.Context) {
		cmd := exec.Command("git", "pull")
		out, _ := cmd.CombinedOutput()
		fmt.Println(string(out))

		c.JSON(200, gin.H{
			"message": string(out),
		})
	})
	r.POST("/api/fypHandler", func(c *gin.Context) {
		var jsonData map[string]string
		data, _ := ioutil.ReadAll(c.Request.Body)
		if e := json.Unmarshal(data, &jsonData); e != nil {
			c.JSON(http.StatusBadRequest, gin.H{"err": e.Error(), "request": data})
			return
		}

		scanId := getTaskid()
		urlSeed := jsonData["urlSeed"]

		// for key, slice := range jsonData {
		// 	if key == "pcrawler" || key == "acrawler" || key == "afuzzer" || key == "scanner" {
		// 		taskIDs = append(taskIDs, toolsHandler(urlSeed, slice)...)
		// 	}
		// }
		// j, _ := json.Marshal(taskIDs)
		// tasks[scanId] = string(j)
		// fmt.Println(tasks[scanId])
		// // tasks[scanId] = fmt.Sprint(j)

		// c.JSON(200, gin.H{"scanId": scanId, "tasks[scanId]": tasks[scanId]})
		var taskIDs []string

		go func() {
			pcrawler := getTaskid()
			afuzzer := getTaskid()
			scanner := getTaskid()
			parafuzzer := getTaskid()
			xsstrike := getTaskid()
			dalfox := getTaskid()
			taskIDs = append(taskIDs, strconv.Itoa(pcrawler))
			taskIDs = append(taskIDs, strconv.Itoa(afuzzer))
			taskIDs = append(taskIDs, strconv.Itoa(scanner))
			taskIDs = append(taskIDs, strconv.Itoa(parafuzzer))
			taskIDs = append(taskIDs, strconv.Itoa(xsstrike))
			taskIDs = append(taskIDs, strconv.Itoa(dalfox))
			j, _ := json.Marshal(taskIDs)
			tasks[scanId] = string(j)

			// acrawler := getTaskid()
			// cmd2(acrawler, "make "+"self_acrawler"+" target="+urlSeed)

			cmd2(pcrawler, "make fyp_pcrawler target="+urlSeed)
			fmt.Println(tasks[pcrawler])
			cmd2(afuzzer, "make fyp_afuzzer target="+urlSeed)
			fmt.Println(tasks[afuzzer])
			// out, _ := cmd("make fyp_parafuzzer")
			// fmt.Println(out)

			go func() {
				cmd2(xsstrike, "make fyp_xsstrike")
			}()
			cmd2(parafuzzer, "make fyp_parafuzzer")
			go func() {
				cmd2(dalfox, "make fyp_dalfox")
			}()
			go func() {
				cmd2(scanner, "make fyp_self_scanner")
			}()

		}()
		c.JSON(200, gin.H{"message": tasks[scanId], "scanId": scanId})
		// c.JSON(200, gin.H{"message": pcrawler + acrawler + afuzzer, "scanID": "123"})
	})
	r.POST("/api/frameworkHandler", func(c *gin.Context) {
		var jsonData map[string]interface{}
		data, _ := ioutil.ReadAll(c.Request.Body)
		if e := json.Unmarshal(data, &jsonData); e != nil {
			c.JSON(http.StatusBadRequest, gin.H{"err": e.Error(), "request": data})
			return
		}

		scanId := getTaskid()
		// c.JSON(200, gin.H{"scanId": scanId})

		var taskIDs []string
		var urlSeed string
		for key, slice := range jsonData {
			if key == "urlSeed" {
				urlSeed = slice.(string)
			} else if key == "pcrawler" || key == "acrawler" || key == "afuzzer" || key == "scanner" {
				for _, ele := range interface2slice(slice) {
					if !stringInSlice(ele, tools2[key]) {
						c.JSON(http.StatusBadRequest, gin.H{"msg": "The input tool" + ele + "is not legitimate"})
						return
					}
				}
			}
		}

		for key, slice := range jsonData {
			if key == "pcrawler" || key == "acrawler" || key == "afuzzer" || key == "scanner" {
				taskIDs = append(taskIDs, toolsHandler(urlSeed, slice)...)
			}
		}
		j, _ := json.Marshal(taskIDs)
		tasks[scanId] = string(j)
		fmt.Println(tasks[scanId])
		// tasks[scanId] = fmt.Sprint(j)

		c.JSON(200, gin.H{"scanId": scanId, "tasks[scanId]": tasks[scanId]})
		// wg.Wait()
	})

	r.POST("/api/commandHandler", func(c *gin.Context) {
		taskIDs := make(map[string]int)
		var jsonData map[string]interface{}
		data, _ := ioutil.ReadAll(c.Request.Body)
		if e := json.Unmarshal(data, &jsonData); e != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": e.Error()})
			return
		}
		fmt.Println(jsonData)
		fmt.Println(len(jsonData))

		// checkValidTools()
		for e1, e2 := range jsonData {
			if e1 == "urlSeed" || e1 == "projID" || e1 == "dependent" {
				continue
			}
			if !stringInSlice(e2.(string), tools[e1].([]string)) {
				c.JSON(http.StatusBadRequest, gin.H{"msg": "The input tools are not legitimate"})
				return
			}
		}

		dependent := false
		if !dependent {
			// for _, ele := range interface2slice(slice) {
			// 	toolsHandler()
			// }
			go func() {
				// fmt.Println(jsonData["asdsadasd"])
				if jsonData["acrawler"] != nil {
					taskID := getTaskid()
					taskIDs["acrawler"] = taskID
					if jsonData["acrawler"] == "self" {
						cmd := exec.Command("go", "run", "aCrawler.go", "-p", "-u", fmt.Sprint(jsonData["urlSeed"]))
						out, _ := cmd.CombinedOutput()

						tasks[taskID] += string(out)
						tasks[taskID] += "This is the last line"
						fmt.Println(string(out))
						// } else if jsonData["acrawler"] == "haralwer" {
					} else if jsonData["acrawler"] != "" {
						// cmd := exec.Command(sEngine, fmt.Sprint(jsonData["acrawler"]), fmt.Sprint(jsonData["urlSeed"]))
						// out, _ := cmd.CombinedOutput()
						out, err := cmd("./sEngine.sh " + fmt.Sprint(jsonData["acrawler"]) + " " + fmt.Sprint(jsonData["urlSeed"]))

						tasks[taskID] += out
						tasks[taskID] += "This is the last line" + fmt.Sprint(err)
						fmt.Println(string(out))
					}
				}
			}()

			go func() {
				if jsonData["afuzzer"] != nil {
					taskID := getTaskid()
					taskIDs["afuzzer"] = taskID
					if jsonData["afuzzer"] == "self" {
						cmd := exec.Command("go", "run", "aCrawler.go", "-w", "test.txt", "-p", "-u", fmt.Sprint(jsonData["urlSeed"]))
						out, _ := cmd.CombinedOutput()

						tasks[taskID] += string(out)
						tasks[taskID] += "This is the last line"
						fmt.Println(string(out))
						// } else if jsonData["acrawler"] == "haralwer" {
					} else {
						// cmd := exec.Command(sEngine, fmt.Sprint(jsonData["afuzzer"]), fmt.Sprint(jsonData["urlSeed"]))
						// out, _ := cmd.CombinedOutput()
						out, err := cmd("./sEngine.sh " + fmt.Sprint(jsonData["afuzzer"]) + " " + fmt.Sprint(jsonData["urlSeed"]))

						tasks[taskID] += out
						tasks[taskID] += "This is the last line" + fmt.Sprint(err)
						fmt.Println(string(out))
					}
				}
			}()

			go func() {
				if jsonData["pcrawler"] != nil {
					taskID := getTaskid()
					taskIDs["pcrawler"] = taskID

					if jsonData["pcrawler"] == "self" {
						cmd := exec.Command("curl", "pcrawl.hb1.workers.dev/?url="+fmt.Sprint(jsonData["urlSeed"]), "-s")
						out, _ := cmd.CombinedOutput()

						tasks[taskID] += string(out)
						tasks[taskID] += "This is the last line" + fmt.Sprint(err)
						fmt.Println(string(out))
					} else {
						// cmd := exec.Command(sEngine, fmt.Sprint(jsonData["pcrawler"]), fmt.Sprint(jsonData["urlSeed"]))
						// out, _ := cmd.CombinedOutput()
						out, err := cmd("./sEngine.sh " + fmt.Sprint(jsonData["pcrawler"]) + " " + fmt.Sprint(jsonData["urlSeed"]))

						tasks[taskID] += out
						tasks[taskID] += "This is the last line" + fmt.Sprint(err)
						fmt.Println(string(out))
					}
					// } else if jsonData["pcrawler"] == "gau" {
					// 	cmd := exec.Command("./sEngine.sh", "gau", fmt.Sprint(jsonData["urlSeed"]))
					// 	out, _ := cmd.CombinedOutput()

					// 	tasks[taskID] += string(out)
					// 	fmt.Println(string(out))
					// }
				}
			}()

			go func() {
				if jsonData["scanner"] != nil {
					taskID := getTaskid()
					taskIDs["scanner"] = taskID

					if jsonData["scanner"] == "nuclei" {
						// cmd := exec.Command("./sEngine.sh", "nuclei", "demo.testfire.net")
						out, err := cmd("./sEngine.sh nuclei " + fmt.Sprint(jsonData["urlSeed"]))

						tasks[taskID] += out
						tasks[taskID] += "This is the last line" + fmt.Sprint(err)
						fmt.Println(string(out))
					} else if jsonData["scanner"] == "dalfox" {
						// cmd := exec.Command(sEngine, "dalfox", fmt.Sprint(jsonData["urlSeed"]))
						// out, _ := cmd.CombinedOutput()
						out, err := cmd("./sEngine.sh dalfox " + fmt.Sprint(jsonData["urlSeed"]))

						tasks[taskID] += out
						tasks[taskID] += "This is the last line" + fmt.Sprint(err)
						fmt.Println(string(out))
					} else {
						out, err := cmd("./sEngine.sh " + fmt.Sprint(jsonData["scanner"]) + " " + fmt.Sprint(jsonData["urlSeed"]))

						tasks[taskID] += out
						tasks[taskID] += "This is the last line" + fmt.Sprint(err)
						fmt.Println(string(out))
					}
				}
			}()
		}
		for len(jsonData) != len(taskIDs)+1 {
			time.Sleep(1)
		}
		// fmt.Println("=============", taskIDs)
		c.JSON(http.StatusOK, taskIDs)

	})
	r.POST("/api/createProj", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong1",
		})
	})
	r.GET("/api/getTasks/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			fmt.Println(err)
		}
		c.JSON(200, gin.H{"message": fmt.Sprintf(tasks[id])})
	})
	r.GET("/api/getTasksId/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			fmt.Println(err)
		}
		var jsonData interface{}
		data := []byte(tasks[id])

		if e := json.Unmarshal(data, &jsonData); e != nil {
			fmt.Println(e.Error())
			return
		}

		c.JSON(200, gin.H{"message": jsonData})
	})

	r.GET("/api/getTools", func(c *gin.Context) {
		c.JSON(200, tools2)
	})
	r.GET("/api/getResult", func(c *gin.Context) {

		if dbcon {

			results, err := db.Query("SELECT time, respCode,respLen,url FROM log")
			if err != nil {
				panic(err.Error())
			}
			var records []string
			for results.Next() {
				var record Record
				err = results.Scan(&record.time, &record.respCode, &record.respLen, &record.url)
				if err != nil {
					panic(err.Error())
				}

				haha := fmt.Sprint(record)
				records = append(records, haha)
			}
			// fmt.Println(haha)
			c.JSON(200, gin.H{
				"message": records,
			})
		} else {
			c.JSON(200, gin.H{
				"message": "db connectio err",
			})
		}
	})

	out, _ := cmd("curl ifconfig.me")
	fmt.Println(out)
	if strings.Contains(out, "206.189.90.86") {
		r.RunTLS(":443", "/etc/letsencrypt/live/hb1.tech/fullchain.pem", "/etc/letsencrypt/live/hb1.tech/privkey.pem")
	} else {
		r.Run(":80")
	}

}
