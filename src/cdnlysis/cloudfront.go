package main

import (
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const TAB = "\t"

var COLUMNS = []string{
	"date",
	"time",
	"x-edge-location",
	"sc-bytes",
	"c-ip",
	"cs-method",
	"cs(Host)",
	"cs-uri-stem",
	"sc-status",
	"cs(Referer)",
	"cs(User-Agent)",
	"cs-uri-query",
	"cs(Cookie)",
	"x-edge-result-type",
	"x-edge-request-id",
	"x-host-header",
	"cs-protocol",
	"cs-bytes",
	"time-taken",
}

type LogEntry struct {
	Date           string  `json:"date"`
	Time           string  `json:"time"`
	EdgeLocation   string  `json:"x-edge-location"`
	BytesSent      int64   `json:"sc-bytes,string"`
	IP             string  `json:"c-ip"`
	Method         string  `json:"cs-method"`
	Host           string  `json:"cs(Host)"`
	UriStem        string  `json:"cs-uri-stem"`
	Status         int64   `json:"sc-status,string"`
	Referer        string  `json:"cs(Referer)"`
	UserAgent      string  `json:"cs(User-Agent)"`
	UriQuery       string  `json:"cs-uri-query"`
	Cookie         string  `json:"cs(Cookie)"`
	EdgeResultType string  `json:"x-edge-result-type"`
	EdgeRequestId  string  `json:"x-edge-request-id"`
	HostHeader     string  `json:"x-host-header"`
	Protocol       string  `json:"cs-protocol"`
	BytesReceived  int64   `json:"cs-bytes,string"`
	TimeTaken      float64 `json:"time-taken,string"`
}

func ParseRecord(log_record string) []interface{} {
	if strings.HasSuffix(log_record, "\n") {
		log_record = log_record[:len(log_record)-1]
	}

	split := strings.Split(log_record, TAB)

	record := []interface{}{}

	datetime := ""

	for ix, item := range COLUMNS {
		log.Println(item)

		val := split[ix]

		if item == "date" {
			datetime += val
		}

		if item == "cs(Referer)" ||
			item == "cs(User-Agent)" ||
			item == "cs-uri-query" {
			unescaped, err := url.QueryUnescape(val)
			if err != nil {
				log.Println(err)
			}

			record = append(record, unescaped)
			continue
		}

		if item == "sc-bytes" ||
			item == "sc-status" ||
			item == "cs-bytes" {
			conv, _ := strconv.ParseInt(val, 10, 64)
			record = append(record, conv)

		} else if item == "time-taken" {
			conv, _ := strconv.ParseFloat(val, 64)
			record = append(record, conv)

		} else if item == "time" {
			datetime += "T" + val + "+00:00"
			t, _ := time.Parse(time.RFC3339, datetime)
			record = append(record, t.Unix())

		} else {
			record = append(record, val)
		}
	}

	return record
}

/*
func ParseRecord(log_record string) *LogEntry {
	if strings.HasSuffix(log_record, "\n") {
		log_record = log_record[:len(log_record)-1]
	}

	split := strings.Split(log_record, TAB)

	json_string := `{`

	var hasElem bool
	for ix, item := range COLUMNS {
		if hasElem {
			json_string += ", "
		} else {
			hasElem = true
		}
		json_string += (`"` + item + `"`)
		json_string += ":"
		json_string += (`"` + split[ix] + `"`)
	}

	json_string += `}`

	var entry LogEntry
	utils.Decoder([]byte(json_string), &entry)
	return &entry
}
*/
