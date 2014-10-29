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
	"x_edge_location",
	"sc_bytes",
	"c_ip",
	"cs_method",
	"cs_host",
	"cs_uri_stem",
	"sc_status",
	"cs_referer",
	"cs_user_agent",
	"cs_uri_query",
	"cs_cookie",
	"x_edge_result_type",
	"x_edge_request_id",
	"x_host_header",
	"cs_protocol",
	"cs_bytes",
	"time_taken",
}

type LogEntry struct {
	Date           string  `json:"date"`
	Time           string  `json:"time"`
	EdgeLocation   string  `json:"x_edge_location"`
	BytesSent      int64   `json:"sc_bytes,string"`
	IP             string  `json:"c_ip"`
	Method         string  `json:"cs_method"`
	Host           string  `json:"cs_host"`
	UriStem        string  `json:"cs_uri_stem"`
	Status         int64   `json:"sc_status,string"`
	Referer        string  `json:"cs_referer"`
	UserAgent      string  `json:"cs_user_agent"`
	UriQuery       string  `json:"cs_uri_query"`
	Cookie         string  `json:"cs_cookie"`
	EdgeResultType string  `json:"x_edge_result_type"`
	EdgeRequestId  string  `json:"x_edge_request_id"`
	HostHeader     string  `json:"x_host_header"`
	Protocol       string  `json:"cs_protocol"`
	BytesReceived  int64   `json:"cs_bytes,string"`
	TimeTaken      float64 `json:"time_taken,string"`
}

func ParseRecord(log_record string) []interface{} {
	if strings.HasSuffix(log_record, "\n") {
		log_record = log_record[:len(log_record)-1]
	}

	split := strings.Split(log_record, TAB)

	record := []interface{}{}

	datetime := ""

	for ix, item := range COLUMNS {
		val := split[ix]

		if item == "date" {
			datetime += val
		}

		if item == "cs_referer" ||
			item == "cs_user_agent" ||
			item == "cs_uri_query" {
			unescaped, err := url.QueryUnescape(val)
			if err != nil {
				log.Println(err)
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
