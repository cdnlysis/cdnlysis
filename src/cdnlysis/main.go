// +build !appengine

package main

import (
	"cdnlysis/backends"
	"cdnlysis/conf"
	"cdnlysis/db"
	"cdnlysis/pipeline"
	"log"
	"sync"
)

func transform(
	files <-chan *pipeline.LogFile,
	influxSink chan<- *backends.InfluxRecord,
	mongoSink chan<- *backends.MongoRecord,
) {
	for file := range files {
		if db.HasVisited(file.Path) {
			log.Println("[Already Processed]", file.Path)
			continue
		}

		log.Println("[Dispatch]", file.Path)
		pipeline.Transform(file, influxSink, mongoSink)
		db.SetVisited(file.Path)
	}

}

func influxAggregator(
	influxSink <-chan *backends.InfluxRecord,
	rg *sync.WaitGroup,
) {
	for v := range influxSink {
		rg.Add(1)
		go func(rec *backends.InfluxRecord) {
			pipeline.AddToInflux(rec)
			rg.Done()
		}(v)
	}
}

func mongoAggregator(
	mongoSink <-chan *backends.MongoRecord,
	rg *sync.WaitGroup,
) {

	mongoRecords := []interface{}{}
	rg.Add(1)

	for v := range mongoSink {
		mongoRecords = append(mongoRecords, v)
	}

	if len(mongoRecords) > 0 {
		pipeline.AddToMongo(mongoRecords)
	}

	rg.Done()
}

func recursivelyWalk(marker *string) {
	it := pipeline.NewIter(*marker)

	incomingFiles := it.Produce(marker)

	var workerWaiter sync.WaitGroup

	//Channel to receive all values that need to be added to InfluxDB
	influxSink := make(chan *backends.InfluxRecord)

	//Channel to receive records to be added to MongoDB
	mongoSink := make(chan *backends.MongoRecord)

	workerWaiter.Add(conf.Settings.Engine.Threads)

	for p := 0; p < conf.Settings.Engine.Threads; p++ {
		go func() {
			transform(incomingFiles, influxSink, mongoSink)
			workerWaiter.Done()
		}()
	}

	var resultGroup sync.WaitGroup

	go influxAggregator(influxSink, &resultGroup)
	go mongoAggregator(mongoSink, &resultGroup)

	workerWaiter.Wait()
	close(influxSink)
	close(mongoSink)
	resultGroup.Wait()

	db.Update(conf.Settings.S3.Prefix, *marker)

	if it.IsTruncated {
		if conf.Settings.Engine.Verbose {
			log.Println("should fetch more")
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

	if marker != "" && conf.Settings.Engine.Verbose {
		log.Println("Resuming state from:", marker)
	}

	recursivelyWalk(&marker)
}
