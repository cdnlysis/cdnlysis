// +build !appengine

package main

import (
	"cdnlysis/db"
	"log"
	"sync"
)

var Settings config

func recursivelyWalk(marker *string) {
	var wg sync.WaitGroup

	it := NewIterator(Settings.S3.Prefix, *marker, &Settings)

	for !it.End() {
		file := it.Next()
		if db.HasVisited(file.Path) {
			log.Println("File", file.Path, "has been processed already")
			continue
		}

		*marker = file.Path

		wg.Add(1)

		go func(wg *sync.WaitGroup, file *LogFile) {
			defer wg.Done()
			log.Println("Processing file " + file.Path)
			ret := processFile(file)
			if ret == true {
				db.SetVisited(file.Path)
			}
		}(&wg, file)
	}

	wg.Wait()

	db.Update(Settings.S3.Prefix, *marker)

	if it.IsTruncated {
		log.Println("should fetch more")
		recursivelyWalk(marker)
	}

}

func main() {
	cliArgs(&Settings)

	db.InitDB(Settings.SyncProgress.Path)

	marker := db.LastMarker(Settings.S3.Prefix)

	if marker != "" {
		log.Println("Resuming state from:", marker)
	}

	recursivelyWalk(&marker)
}
