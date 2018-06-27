package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

// Cache for thread history
var Cache []*Thread

func isInCache(t *Thread) bool {
	for _, v := range Cache {
		return t.Name == v.Name && t.Href == v.Href
	}
	return false
}

func addToCache(t *Thread) {
	Cache = append(Cache, t)
}

func loadCache() {
	file, err := os.OpenFile("cache.json", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)

	if err != nil {
		log.Fatal("Unable to load cache")
	}

	byteValues, _ := ioutil.ReadAll(file)

	if len(byteValues) > 0 {
		if err = json.Unmarshal(byteValues, &Cache); err != nil {
			file.Close()
			Progress.Prefix(fmt.Sprintf("Unable to parse cache: %s", err))
		}
	}

	defer file.Close()
}

func saveCache() {
	byteValues, err := json.Marshal(&Cache)

	if err != nil {
		log.Fatal("Unable to encode cache")
	}

	if err = ioutil.WriteFile("cache.json", byteValues, 0777); err != nil {
		log.Fatal("Unable to write cache: ", err)
	}
}
