package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"sync"
)

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

func cmd(str string) (string, error) {
	cmd := exec.Command("bash", "-c", str)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func toolsHandler(urlSeed string, b interface{}) {
	fmt.Println(b)

	for _, tool := range interface2slice(b) {
		wg.Add(1)
		go func(tool string) {
			defer wg.Done()
			fmt.Println("running", tool)

			out, err := cmd("./sEngine.sh " + tool + " " + urlSeed)
			// tasks[taskID] += out
			// tasks[taskID] += "This is the last line" + fmt.Sprint(err)
			fmt.Println(string(out))
			fmt.Println(err)
		}(tool)
	}
}

// func scannerHandeler(){

// }

// var tools := make(map[string]interface{})
var tools = make(map[string][]string)
var wg sync.WaitGroup

func main() {
	// tools
	tools["pcrawler"] = []string{"self", "waybackurls", "gau"}
	tools["acrawler"] = []string{"self", "hakrawler"}
	tools["afuzzer"] = []string{"self", "ffuf"}
	tools["scanner"] = []string{"nuclei", "dalfox", "whatweb"}

	// var jsonData map[string][]string
	var jsonData map[string]interface{}
	data := []byte(`{"urlSeed":"demo.testfire.net","scanner":["nuclei","dalfox","whatweb"]}`)

	if e := json.Unmarshal(data, &jsonData); e != nil {
		fmt.Println(e.Error())
		return
	}

	// fmt.Println(jsonData)
	// fmt.Println(len(jsonData))

	var urlSeed string
	// var urlSeeds []string
	// check maclious input
	for key, slice := range jsonData {
		// fmt.Println(key, slice)
		// fmt.Println(e2[0])
		if key == "urlSeed" {
			urlSeed = slice.(string)
			// for _, ele := range slice {
			// 	urlSeeds = append(urlSeeds, slice.(string))
			// }
			continue
		} else if key == "projID" || key == "dependent" {
			continue
		} else if key == "pcrawler" || key == "acrawler" || key == "afuzzer" || key == "scanner" {
			for _, ele := range interface2slice(slice) {
				fmt.Println(ele)
				if !stringInSlice(ele, tools[key]) {
					// if !stringInSlice(ele, tools[key].([]string)) {
					// c.JSON(http.StatusBadRequest, gin.H{"msg": "The input tools are not legitimate"})
					fmt.Println("The tool", ele, "is not legit")
					return
				}
			}
		}
	}

	fmt.Println("===1223")
	fmt.Println(urlSeed)

	for key, slice := range jsonData {
		if key == "pcrawler" || key == "acrawler" || key == "afuzzer" || key == "scanner" {
			toolsHandler(urlSeed, slice)
		}
	}
	wg.Wait()

	// for len(jsonData) != len(taskIDs)+1 {
	// 	time.Sleep(1)
	// }
	// // fmt.Println("=============", taskIDs)
	// c.JSON(http.StatusOK, taskIDs)

}
