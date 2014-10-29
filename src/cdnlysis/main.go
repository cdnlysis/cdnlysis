// +build !appengine

package main

import "sync"

var Settings config

func main() {
	cliArgs(&Settings)

	var wg sync.WaitGroup

	for it := NewIterator(Settings.S3.Prefix, &Settings); !it.End(); {
		file := it.Next()

		wg.Add(1)

		go func(wg *sync.WaitGroup, file *LogFile) {
			defer wg.Done()
			processFile(file)
		}(&wg, file)
	}

	wg.Wait()
}
