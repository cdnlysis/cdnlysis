package main

import (
	"bufio"
	"compress/gzip"
	"io"
	"log"

	influxdb "github.com/influxdb/influxdb/client"
	"labix.org/v2/mgo"
)

func addToInflux(series *influxdb.Series) {
	conn, err := influxdb.New(&Settings.Influx)
	if err != nil {
		log.Println(err)
		return
	}

	if err := conn.WriteSeries([]*influxdb.Series{series}); err != nil {
		log.Println(err)
		return
	}
}

type mongoSeries []*LogEntry

func addToMongo(records *mongoSeries) {

	info := mgo.DialInfo{
		Addrs:    []string{Settings.Mongo.Host},
		Database: Settings.Mongo.Database,
		Direct:   true,
		Username: Settings.Mongo.Username,
		Password: Settings.Mongo.Password,
	}

	var err error

	session, err := mgo.DialWithInfo(&info)
	if err != nil {
		log.Println(err)
		return
	}

	conn := session.DB(Settings.Mongo.Database)
	coll := conn.C(Settings.Mongo.Collection)

	var interfaceSlice []interface{} = make([]interface{}, len(*records))
	for i, d := range *records {
		interfaceSlice[i] = d
	}

	err = coll.Insert(interfaceSlice...)
	if err != nil {
		log.Println(err)
	}
}

func processFile(file *LogFile) {
	reader, err := file.GetReader()
	if err != nil {
		log.Println(err)
		return
	}

	defer reader.Close()

	gzipReader, err2 := gzip.NewReader(reader)
	if err2 != nil {
		log.Println(err2)
		return
	}

	defer gzipReader.Close()

	ix := 0
	bufReader := bufio.NewReader(gzipReader)

	series := influxdb.Series{Settings.Logs.Prefix, COLUMNS, nil}
	mongo_records := mongoSeries{}

	for {
		ix++
		log_record, err := bufReader.ReadString('\n') //
		if err == io.EOF {
			//do something here
			break
		} else if err != nil {
			break
			// if you return error
		} else if ix > 2 {
			// Log Entries

			if Settings.Backends.Influx {
				data := InfluxRecord(log_record)
				series.Points = append(series.Points, data)
			}

			if Settings.Backends.Mongo {
				mongo_records = append(mongo_records, MongoRecord(log_record))
			}
		}
	}

	if Settings.Backends.Mongo {
		addToInflux(&series)
	}

	if Settings.Backends.Mongo {
		addToMongo(&mongo_records)
	}
}
