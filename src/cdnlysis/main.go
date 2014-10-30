// +build !appengine

package main

import (
	"cdnlysis/db"
	"log"
	"sync"
)

var Settings config

func main() {
	cliArgs(&Settings)

	db.InitDB(Settings.SyncProgress.Path)

	var wg sync.WaitGroup

	marker := db.LastMarker(Settings.S3.Prefix)

	log.Println(marker)

	for it := NewIterator(Settings.S3.Prefix, marker, &Settings); !it.End(); {
		file := it.Next()

		db.Update(Settings.S3.Prefix, file.Path)

		wg.Add(1)

		go func(wg *sync.WaitGroup, file *LogFile) {
			defer wg.Done()
			log.Println("Processing file " + file.Path)
			processFile(file)
		}(&wg, file)
	}

	wg.Wait()
}
