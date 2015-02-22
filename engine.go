// +build !appengine

package cdnlysis

import (
	"errors"
	"log"

	"github.com/meson10/cdnlysis/conf"
	"github.com/meson10/cdnlysis/db"
)

func transform(
	files <-chan *LogFile,
	channels *[]LogRecordChannel,
	errc chan<- *TransformError,
) {
	for file := range files {
		if db.HasVisited(file.Path) {
			errc <- &TransformError{
				file.Path,
				errors.New("Already Processed"),
			}
			continue
		}

		log.Println(file.LogIdent(), "[Pipeline] Transforming")

		Transform(file, channels, errc)
		db.SetVisited(file.Path)
	}
}

func Start(marker *string, channels *[]LogRecordChannel) {
	if *marker != "" {
		log.Println("Resuming state from:", *marker)
	}

	it := NewIter(*marker)

	incomingFiles := it.Produce(marker)

	expectedResponses := conf.Settings.Engine.Threads
	workerChan := make(chan error, expectedResponses)
	errc := make(chan *TransformError)
	allDone := make(chan bool, 1)

	go func() {
		log.Println("Waiting for Errors")
		for err := range errc {
			log.Println(err.Path, err.Err)
		}
		allDone <- true
	}()

	for p := 0; p < conf.Settings.Engine.Threads; p++ {
		go func() {
			transform(incomingFiles, channels, errc)
			workerChan <- nil
		}()
	}

	for {
		err := <-workerChan

		// Display the result.
		if err != nil {
			log.Println("Received error:", err)
		} else {
			log.Println("Received nil error")
		}

		expectedResponses--
		if expectedResponses == 0 {
			break
		}
	}

	<-allDone

	db.Update(conf.Settings.S3.Prefix, *marker)

	if it.IsTruncated {
		Start(marker, channels)
	} else {
		if conf.Settings.Engine.Verbose {
			log.Println("Does not have more values")
		}
	}

}

func Setup() {
	conf.GetConfig()

	log.Println(conf.Settings)

	db.InitDB(conf.Settings.SyncProgress.Path)
}
