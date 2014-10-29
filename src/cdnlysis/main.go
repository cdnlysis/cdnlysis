// +build !appengine

package main

import (
	"log"
	"sync"
)

func findMatching(cfg *config) {
	for it := Iterator("plain", cfg); !it.End(); {
		x := it.Next()
		log.Println(x)
	}
}

func processFile(wg *sync.WaitGroup, cfg *config, path string) {
	defer wg.Done()

	log.Println("Processing File: ", path)
}

func main() {
	var Settings config
	cliArgs(&Settings)

	findMatching(&Settings)
	/*
		if len(files) == 0 {
			return
		}

		var wg sync.WaitGroup

		for _, file := range files {
			wg.Add(1)
			go processFile(&wg, &Settings, file)
		}

		wg.Wait()
	*/
}
