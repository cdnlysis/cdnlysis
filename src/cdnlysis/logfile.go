package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"io"
	"log"
	"regexp"
	"strings"

	influxdb "github.com/influxdb/influxdb/client"
	"labix.org/v2/mgo"
)

func addToInflux(series *influxdb.Series) {
	conn, err := influxdb.New(&Settings.Influx)
	if err != nil {
		log.Println("Cannot connect to Influx", err)
		return
	}

	if err := conn.WriteSeries([]*influxdb.Series{series}); err != nil {
		log.Println("Cannot add to Influx", err)
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
		log.Println("Cannot connect to Mongo:", err)
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
		log.Println("Cannot Insert document in collection", err)
	}
}

func cleanString(msg string) string {
	if strings.HasSuffix(msg, "\n") {
		return msg[:len(msg)-1]
	}

	return msg
}

func findColumns(msg string) []string {
	re_head := regexp.MustCompile("#Fields:\\s+")
	re_slug := regexp.MustCompile("[^\\w\\s]")

	msg = re_head.ReplaceAllString(msg, "")
	msg = strings.ToLower(re_slug.ReplaceAllString(msg, UNDERSCORE))
	return strings.Split(msg, SPACE)
}

func processFile(file *LogFile) bool {
	buff, err := file.Get()
	if err != nil {
		log.Println("Cannot GetFIle", err)
		return false
	}

	b := bytes.NewReader(buff)
	gzipReader, err2 := gzip.NewReader(b)
	if err2 != nil {
		log.Println("Cannot make GZIP Reader", err2)
		return false
	}

	defer gzipReader.Close()

	ix := 0
	bufReader := bufio.NewReader(gzipReader)

	var columns []string

	series := influxdb.Series{Settings.Logs.Prefix, columns, nil}
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
		}

		log_record = cleanString(log_record)

		if ix == 1 {
			columns = findColumns(log_record)
		} else if ix > 2 {
			// Log Entries

			if Settings.Backends.Influx {
				data := InfluxRecord(columns, log_record)
				series.Points = append(series.Points, data)
			}

			if Settings.Backends.Mongo {
				data := MongoRecord(columns, log_record)
				mongo_records = append(mongo_records, data)
			}
		}
	}

	if Settings.Backends.Mongo {
		series.Columns = columns
		addToInflux(&series)
	}

	if Settings.Backends.Mongo {
		addToMongo(&mongo_records)
	}

	return true
}
