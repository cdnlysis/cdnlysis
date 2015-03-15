// +build !appengine

package cdnlysis

import (
	"errors"
	"log"

	"gopkg.in/cdnlysis/cdnlysis.v1/conf"
	"gopkg.in/cdnlysis/cdnlysis.v1/db"
)

func transform(
	files <-chan *LogFile,
	channel chan<- *LogRecord,
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

		Transform(file, channel, errc)
		db.SetVisited(file.Path)
	}
}

func keepGoing(marker *string, channel chan<- *LogRecord) {
	if *marker != "" {
		log.Println("Resuming state from:", *marker)
	}

	it := NewIter(*marker)

	incomingFiles := it.Produce(marker)

	expectedResponses := conf.Settings.Engine.Threads
	workerChan := make(chan error, expectedResponses)
	errc := make(chan *TransformError)

	go func() {
		log.Println("Waiting for Errors")
		for err := range errc {
			log.Println("Error:", err.Path, err.Err)
		}
	}()

	for p := 0; p < conf.Settings.Engine.Threads; p++ {
		go func() {
			transform(incomingFiles, channel, errc)
			workerChan <- nil
		}()
	}

	for {
		err := <-workerChan

		// Display the result.
		if err != nil {
			log.Println("Received error:", err)
		}

		expectedResponses--
		if expectedResponses == 0 {
			break
		}
	}

	db.Update(conf.Settings.S3.Prefix, *marker)

	if it.IsTruncated {
		Start(marker, channel)
	} else {
		if conf.Settings.Engine.Verbose {
			log.Println("Does not have more values")
		}
	}
}

func Start(prefix *string, channel chan<- *LogRecord) {
	var marker string

	if *prefix == "" {
		marker = db.LastMarker(conf.Settings.S3.Prefix)
	}

	keepGoing(&marker, channel)

	log.Println("Closing Channel")
	close(channel)
}

func Setup() {
	conf.GetConfig()

	log.Println(conf.Settings)

	db.InitDB(conf.Settings.SyncProgress.Path)
}
