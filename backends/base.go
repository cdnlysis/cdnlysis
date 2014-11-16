package backends

import "strings"

const TAB = "\t"
const SPACE = " "
const UNDERSCORE = "_"

func parseLogRecord(log_record string) []string {
	return strings.Split(log_record, TAB)
}
