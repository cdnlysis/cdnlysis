package pipeline

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"io"
	"log"
	"regexp"
	"strings"

	"cdnlysis/backends"
	"cdnlysis/conf"
)

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
	msg = strings.ToLower(re_slug.ReplaceAllString(msg, backends.UNDERSCORE))
	return strings.Split(msg, backends.SPACE)
}

type TransformError struct {
	Path string
	Err  error
}

func Transform(
	file *LogFile,
	influxSink chan<- *backends.InfluxRecord,
	mongoSink chan<- *backends.MongoRecord,
	errc chan<- *TransformError,
) {

	log.Println(file.LogIdent(), "[Fetch]")

	buff, err := file.Get()
	if err != nil {
		log.Println(file.LogIdent(), "[Error] Cannot open File", err)
		errc <- &TransformError{file.Path, err}
		return
	}

	b := bytes.NewReader(buff)
	gzipReader, err2 := gzip.NewReader(b)
	if err2 != nil {
		log.Println(file.LogIdent(), "[Error] Cannot open GZIP", err2)
		errc <- &TransformError{file.Path, err2}
		return
	}

	defer gzipReader.Close()

	ix := 0
	bufReader := bufio.NewReader(gzipReader)

	var columns []string

	series := backends.InfluxRecord{
		conf.Settings.Logs.Prefix,
		columns,
		nil,
	}

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

		if ix == 2 {
			columns = findColumns(log_record)
		} else if ix > 2 {
			// Log Entries

			if conf.Settings.Backends.Influx {
				data := backends.MakeInfluxRecord(columns, log_record)
				series.Points = append(series.Points, data)
			}

			if conf.Settings.Backends.Mongo {
				record := backends.MakeMongoRecord(columns, log_record)
				mongoSink <- record
			}
		}
	}

	if conf.Settings.Backends.Influx {
		series.Columns = columns
		influxSink <- &series
	}
}
