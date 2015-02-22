package cdnlysis

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"io"
	"log"
	"regexp"
	"strings"
)

const TAB = "\t"
const SPACE = " "
const UNDERSCORE = "_"

func parseLogRecord(log_record string) *[]string {
	split := strings.Split(log_record, TAB)
	return &split
}

func cleanString(msg string) string {
	if strings.HasSuffix(msg, "\n") {
		return msg[:len(msg)-1]
	}

	return msg
}

func findColumns(msg string) *[]string {
	re_head := regexp.MustCompile("#Fields:\\s+")
	re_slug := regexp.MustCompile("[^\\w\\s]")

	msg = re_head.ReplaceAllString(msg, "")
	msg = strings.ToLower(re_slug.ReplaceAllString(msg, UNDERSCORE))

	split := strings.Split(msg, SPACE)
	return &split
}

type TransformError struct {
	Path string
	Err  error
}

func Transform(
	file *LogFile,
	channel chan<- *LogRecord,
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

	var columns *[]string

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

			record, err := NewRecord(
				columns,
				parseLogRecord(log_record),
			)

			if err != nil {
				errc <- &TransformError{file.Path, err}
				continue
			}

			channel <- record
		}
	}
}
