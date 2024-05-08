package main

import (
	"bufio"
	"fmt"

	"os"
)

func main() {
	urlSeed := "demo.testfire.net"
	dictionary := []string{}
	file, _ := os.Open("common.txt")

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// dictionary = AppendIfMissing(dictionary, scanner.Text())
		dictionary = append(dictionary, scanner.Text())

		// re, _ := regexp.Compile("\\.[a-z]{2,}$") // append extension if none
		// if !re.MatchString(scanner.Text()) {
		// 	for _, i := range ext {
		// 		dictionary = AppendIfMissing(dictionary, scanner.Text()+"."+i)
		// 	}
		// }
	}
	for _, ele := range dictionary {
		url := fmt.Sprintf("http://%s/%s", urlSeed, ele)
		fmt.Println(url)
	}

	fmt.Println(len(dictionary))
}
