package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
)

func checkStatus(url string) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Print("\n")
		fmt.Println("Link: ", url, " | Status:", resp.Status, "!!!!!!")
	} else {
		fmt.Print(".")
	}
}

func main() {
	// initial values
	type initPar struct {
		URLPage string
		CDN     string
	}
	var readFile string
	var ipar initPar
	fmt.Printf("Read from configuration file (conf.json)? (Y/n)")
	fmt.Scanln(&readFile)
	if readFile == "n" || readFile == "N" {
		fmt.Printf("Page URL to check jpeg images: ")
		fmt.Scanln(&ipar.URLPage)
		fmt.Printf("CDN on the page: ")
		fmt.Scanln(&ipar.CDN)
	} else {
		// Open our jsonFile
		byteValue, err := os.ReadFile("conf.json")
		// if we os.ReadFile returns an error then handle it
		if err != nil {
			fmt.Println(err)
		}
		// defer the closing of our jsonFile so that we can parse it later on
		err = json.Unmarshal(byteValue, &ipar)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("Page URL to test jpeg images: ", ipar.URLPage)
		fmt.Println("CDN on the page: ", ipar.CDN)
	}
	// Send an HTTP GET request to the urlPage web page
	resp, err := http.Get(ipar.URLPage)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()
	// find all matched links
	re := regexp.MustCompile(`(https://` + ipar.CDN + `[\w/.:]*(webp|jpg))`)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	content := string(body)
	links := re.FindAllString(content, -1)
	// loops through the links slice to check get status
	// increase number of requests by 10
	for i := 0; i < 10; i++ {
		for _, l := range links {
			go checkStatus(l)
		}
	}
	var input string
	fmt.Scanln(&input)
}
