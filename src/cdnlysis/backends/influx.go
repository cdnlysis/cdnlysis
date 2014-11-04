package backends

import (
	"log"
	"net/url"
	"strconv"
	"time"

	influxdb "github.com/influxdb/influxdb/client"
)

type InfluxRecord influxdb.Series

func MakeInfluxRecord(columns []string, log_record string) []interface{} {
	split := parseLogRecord(log_record)
	record := []interface{}{}

	datetime := ""

	for ix, item := range columns {
		val := split[ix]

		if item == "date" {
			datetime += val
		}

		if item == "cs_referer" ||
			item == "cs_user_agent" ||
			item == "cs_uri_query" {
			unescaped, err := url.QueryUnescape(val)
			if err != nil {
				log.Println("Cannot unescape value", err)
			}

			record = append(record, unescaped)
			continue
		}

		if item == "sc_bytes" ||
			item == "sc_status" ||
			item == "cs_bytes" {
			conv, _ := strconv.ParseInt(val, 10, 64)
			record = append(record, conv)

		} else if item == "time_taken" {
			conv, _ := strconv.ParseFloat(val, 64)
			record = append(record, conv)

		} else if item == "time" {
			datetime += "T" + val + "+00:00"
			t, _ := time.Parse(time.RFC3339, datetime)
			record = append(record, t.Unix()*1000)

		} else {
			record = append(record, val)
		}
	}

	return record
}
