package main

import (
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"

	utils "github.com/Simversity/gottp/utils"
)

const TAB = "\t"
const SPACE = " "
const UNDERSCORE = "_"

type LogEntry struct {
	Date           string  `bson:"date" json:"date"`
	Time           string  `bson:"time" json:"time"`
	EdgeLocation   string  `bson:"x_edge_location" json:"x_edge_location"`
	BytesSent      int64   `bson:"sc_bytes" json:"sc_bytes,string"`
	IP             string  `bson:"c_ip" json:"c_ip"`
	Method         string  `bson:"cs_method" json:"cs_method"`
	Host           string  `bson:"cs_host" json:"cs_host"`
	UriStem        string  `bson:"cs_uri_stem" json:"cs_uri_stem"`
	Status         int     `bson:"sc_status" json:"sc_status,string"`
	Referer        string  `bson:"cs_referer" json:"cs_referer"`
	UserAgent      string  `bson:"cs_user_agent" json:"cs_user_agent"`
	UriQuery       string  `bson:"cs_uri_query" json:"cs_uri_query"`
	Cookie         string  `bson:"cs_cookie" json:"cs_cookie"`
	EdgeResultType string  `bson:"x_edge_result_type" json:"x_edge_result_type"`
	EdgeRequestId  string  `bson:"x_edge_request_id" json:"x_edge_request_id"`
	HostHeader     string  `bson:"x_host_header" json:"x_host_header"`
	Protocol       string  `bson:"cs_protocol" json:"cs_protocol"`
	BytesReceived  int64   `bson:"cs_bytes" json:"cs_bytes,string"`
	TimeTaken      float64 `bson:"time_taken" json:"time_taken,string"`
}

func InfluxRecord(columns []string, log_record string) []interface{} {
	split := strings.Split(log_record, TAB)

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

func MongoRecord(columns []string, log_record string) *LogEntry {
	split := strings.Split(log_record, TAB)

	json_string := `{`

	var hasElem bool
	for ix, item := range columns {
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
