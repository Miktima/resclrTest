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
		URLPage   string
		CDN       string
		FlagImage int
	}
	var readFile string
	var ipar initPar
	fmt.Printf("Read from configuration file (conf.json)? (Y/n)")
	fmt.Scanln(&readFile)
	if readFile == "n" || readFile == "N" {
		fmt.Printf("Test image (1) of page (0)?")
		fmt.Scanln(&ipar.FlagImage)
		if ipar.FlagImage == 0 {
			fmt.Printf("Page URL to check jpeg images: ")
			fmt.Scanln(&ipar.URLPage)
			fmt.Printf("CDN on the page: ")
			fmt.Scanln(&ipar.CDN)
		} else if ipar.FlagImage == 1 {
			fmt.Printf("URL to image: ")
			fmt.Scanln(&ipar.URLPage)
		} else {
			fmt.Println("ERROR: Unknown option: ", ipar.FlagImage)
		}
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
		fmt.Println("FlagImage: ", ipar.FlagImage)
		if ipar.FlagImage == 0 {
			fmt.Println("Page URL to test jpeg images: ", ipar.URLPage)
			fmt.Println("CDN on the page: ", ipar.CDN)
		} else if ipar.FlagImage == 1 {
			fmt.Println("URL to image: ", ipar.URLPage)
		} else {
			fmt.Println("ERROR: Unknown option: ", ipar.FlagImage)
		}
	}
	if ipar.FlagImage == 0 {
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
	} else if ipar.FlagImage == 1 {
		for i := 0; i < 1500; i++ {
			go checkStatus(ipar.URLPage)
		}
	}
	var input string
	fmt.Scanln(&input)
}
