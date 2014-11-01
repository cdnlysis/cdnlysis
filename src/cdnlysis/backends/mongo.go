package backends

import utils "github.com/Simversity/gottp/utils"

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

func MongoRecord(columns []string, log_record string) *LogEntry {
	split := parseLogRecord(log_record)

	json_string := `{`

	var hasElem bool
	for ix, item := range columns {
		if hasElem {
			json_string += `, `
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
