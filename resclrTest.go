package main

import (
	"encoding/json"
	"fmt"
	"hash/crc32"
	"io"
	"net/http"
	"os"
	"regexp"
)

func getHash(filename, url string) {
	var hFileSum32 chan uint32 = make(chan uint32)
	var hURLSum32 chan uint32 = make(chan uint32)
	go func() {
		bs, err := os.ReadFile(filename)
		if err != nil {
			return
		}
		hFile := crc32.NewIEEE()
		hFile.Write(bs)
		hFileSum32 <- hFile.Sum32()
	}()
	h1 := <-hFileSum32
	go func() {
		resp, err := http.Get(url)
		if err != nil {
			return
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return
		}
		hUrl := crc32.NewIEEE()
		hUrl.Write(body)
		hURLSum32 <- hUrl.Sum32()
	}()
	h2 := <-hURLSum32
	ver := h1 == h2
	switch ver {
	case true:
		fmt.Println(filename, h1, h2, ver)
	case false:
		fmt.Println(filename, h1, h2, ver, "!!!!!!!")
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
		fmt.Printf("Page URL to check css and js resources: ")
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
		fmt.Println("Page URL to check css and js resources: ", ipar.URLPage)
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
	// loops through the links slice and find corresponded files in git
	// Check if paths must be changed (for windows)
	reCdn := regexp.MustCompile(`//cdn[\w*].img`)
	for _, l := range links {
		lWoutCdn := reCdn.ReplaceAllString(l, "//img")
		fmt.Println(lWoutCdn)
		// reFile := regexp.MustCompile(`https://` + ipar.CDN + `(/[\w/.:]*(css|js))`)
		// paths := reFile.FindStringSubmatch(l)
		// file := paths[len(paths)-2]
		// go getHash(file, l)
	}
	var input string
	fmt.Scanln(&input)
}
