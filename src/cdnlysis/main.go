// +build !appengine

package main

import (
	"cdnlysis/backends"
	"cdnlysis/conf"
	"cdnlysis/db"
	"cdnlysis/pipeline"
	"errors"
	"log"
	"sync"
)

func transform(
	files <-chan *pipeline.LogFile,
	influxSink chan<- *backends.InfluxRecord,
	mongoSink chan<- *backends.MongoRecord,
	errc chan<- *pipeline.TransformError,
) {
	for file := range files {
		if db.HasVisited(file.Path) {
			errc <- &pipeline.TransformError{
				file.Path, errors.New("Already Processed"),
			}

			//log.Println(file.LogIdent(), "[Pipeline] Already processed")
			continue
		}

		log.Println(file.LogIdent(), "[Pipeline] Transforming")
		pipeline.Transform(file, influxSink, mongoSink, errc)
		db.SetVisited(file.Path)
	}
}

func influxAggregator(influxSink <-chan *backends.InfluxRecord) {
	var wg sync.WaitGroup

	for v := range influxSink {
		wg.Add(1)
		go func(rec *backends.InfluxRecord) {
			pipeline.AddToInflux(rec)
			wg.Done()
		}(v)
	}

	wg.Wait()
}

func mongoAggregator(mongoSink <-chan *backends.MongoRecord) {
	var wg sync.WaitGroup

	for v := range mongoSink {
		wg.Add(1)
		go func(rec *backends.MongoRecord) {
			mongoRecords := []interface{}{rec}
			pipeline.AddToMongo(mongoRecords)
			wg.Done()
		}(v)
	}

	wg.Wait()
}

func recursivelyWalk(marker *string) {
	it := pipeline.NewIter(*marker)

	incomingFiles := it.Produce(marker)

	var workerWaiter sync.WaitGroup

	errc := make(chan *pipeline.TransformError)

	//Channel to receive all values that need to be added to InfluxDB
	influxSink := make(chan *backends.InfluxRecord)

	//Channel to receive records to be added to MongoDB
	var err error
	err = pipeline.RefreshInfluxSession()
	if err != nil {
		log.Println(err)
		return
	}

	err = pipeline.RefreshMongoSession()
	if err != nil {
		log.Println(err)
		return
	}

	mongoSink := make(chan *backends.MongoRecord)

	workerWaiter.Add(conf.Settings.Engine.Threads)

	for p := 0; p < conf.Settings.Engine.Threads; p++ {
		go func() {
			transform(incomingFiles, influxSink, mongoSink, errc)
			workerWaiter.Done()
		}()
	}

	var resultGroup sync.WaitGroup
	resultGroup.Add(3)

	go func() {
		log.Println("Waiting for Errors")
		for err := range errc {
			log.Println(err.Path, err.Err)
		}
		resultGroup.Done()
	}()

	go func() {
		log.Println("Waiting for Influxers")
		influxAggregator(influxSink)
		resultGroup.Done()
	}()

	go func() {
		log.Println("Waiting for Mongo")
		mongoAggregator(mongoSink)
		resultGroup.Done()
	}()

	workerWaiter.Wait()

	close(errc)
	close(influxSink)
	close(mongoSink)

	resultGroup.Wait()

	db.Update(conf.Settings.S3.Prefix, *marker)

	if it.IsTruncated {
		if conf.Settings.Engine.Verbose {
			log.Println("Fetching next batch of values")
		}
		recursivelyWalk(marker)
	} else {
		if conf.Settings.Engine.Verbose {
			log.Println("Does not have more values")
		}
	}

}

func main() {
	conf.CliArgs()

	db.InitDB(conf.Settings.SyncProgress.Path)

	marker := db.LastMarker(conf.Settings.S3.Prefix)
	marker = ""

	if marker != "" && conf.Settings.Engine.Verbose {
		log.Println("Resuming state from:", marker)
	}

	recursivelyWalk(&marker)
}
