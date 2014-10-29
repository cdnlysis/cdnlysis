package main

import (
	"bufio"
	"compress/gzip"
	"io"
	"log"

	influxdb "github.com/influxdb/influxdb/client"
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
			data := ParseRecord(log_record)
			series.Points = append(series.Points, data)
		}
	}

	addToInflux(&series)
}
