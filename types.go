package cdnlysis

import "errors"

type LogRecord struct {
	Columns *[]string
	Values  *[]string
}

func NewRecord(columns *[]string, values *[]string) (*LogRecord, error) {
	if len(*columns) != len(*values) {
		return nil, errors.New("Columns and Values should be of the same length.")
	}

	return &LogRecord{columns, values}, nil
}
