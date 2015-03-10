package cdnlysis

import (
	"errors"
	"reflect"
	"time"
)

type LogRecord struct {
	Columns *[]string
	Values  *[]string

	Time float64

	Date           string  `json:"date"`
	Timestamp      string  `json:"time"`
	EdgeLocation   string  `json:"x_edge_location"`
	BytesSent      int64   `json:"sc_bytes"`
	IP             string  `json:"c_ip"`
	Method         string  `json:"cs_method"`
	Host           string  `json:"cs_host_"`
	UriStem        string  `json:"cs_uri_stem"`
	Status         int     `json:"sc_status"`
	UriQuery       string  `json:"cs_uri_query"`
	EdgeResultType string  `json:"x_edge_result_type"`
	EdgeRequestId  string  `json:"x_edge_request_id"`
	HostHeader     string  `json:"x_host_header"`
	Protocol       string  `json:"cs_protocol"`
	BytesReceived  int64   `json:"cs_bytes"`
	TimeTaken      float64 `json:"time_taken"`
	Referer        string  `json:"cs_referer_"`
	UserAgent      string  `json:"cs_user_agent_"`
	Cookie         string  `json:"cs_cookie_"`
}

func (self *LogRecord) Convert() {
	typ := reflect.TypeOf(self)
	val := reflect.ValueOf(self)

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		val = val.Elem()
	}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)

		// Skip ignored and unexported fields in the struct
		if !val.Field(i).CanInterface() {
			continue
		}

		// Cannot Validate nested and embedded structs
		if field.Type.Kind() == reflect.Struct ||
			(field.Type.Kind() == reflect.Ptr &&
				field.Type.Elem().Kind() == reflect.Struct) {
			continue
		}

		name := field.Name
		if j := field.Tag.Get("json"); j != "" {
			name = j
		}

		index := -1
		for p, v := range *self.Columns {
			if v == name {
				index = p
				break
			}
		}

		if index == -1 {
			continue
		}

		v := val.Field(i)
		actualValue := (*self.Values)[index]
		convertValue(v, actualValue)
	}

	datetime := self.Date + "T" + self.Timestamp + "+00:00"
	t, _ := time.Parse(time.RFC3339, datetime)
	self.Time = float64(t.Unix() * 1000)
}

func NewRecord(columns *[]string, values *[]string) (*LogRecord, error) {
	if len(*columns) != len(*values) {
		return nil, errors.New("Columns and Values should be of the same length.")
	}

	record := &LogRecord{}
	record.Columns = columns
	record.Values = values

	return record, nil
}
