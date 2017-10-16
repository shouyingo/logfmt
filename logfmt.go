package logfmt

import (
	"bytes"
	"fmt"
	"sort"

	"github.com/sirupsen/logrus"
)

var (
	SortKeys = false
)

var (
	tags = [...]string{
		logrus.PanicLevel: "[PNC] ",
		logrus.FatalLevel: "[FAT] ",
		logrus.ErrorLevel: "[ERR] ",
		logrus.WarnLevel:  "[WRN] ",
		logrus.InfoLevel:  "[INF] ",
		logrus.DebugLevel: "[DBG] ",
	}
)

var DefaultFormatter *PlainFormatter

type PlainFormatter struct{}

func (*PlainFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	buf := entry.Buffer
	if buf == nil {
		buf = bytes.NewBuffer(nil)
	}

	// time
	var t [32]byte
	buf.Write(entry.Time.AppendFormat(t[:0], "2006-01-02/15:04:05.000 "))

	// level
	buf.WriteString(tags[entry.Level])

	if len(entry.Data) > 0 {
		// fields
		buf.WriteByte('{')
		if !SortKeys {
			more := false
			for k, v := range entry.Data {
				if more {
					buf.WriteString(", ")
				} else {
					more = true
				}
				buf.WriteString(k)
				buf.WriteString(": ")
				fmt.Fprint(buf, v)
			}
		} else {
			keys := make([]string, len(entry.Data))
			p := 0
			for k := range entry.Data {
				keys[p] = k
				p++
			}
			keys = keys[:p]
			sort.Strings(keys)
			more := false
			for _, k := range keys {
				if more {
					buf.WriteString(", ")
				} else {
					more = true
				}
				buf.WriteString(k)
				buf.WriteString(": ")
				fmt.Fprint(buf, entry.Data[k])
			}
		}
		buf.WriteString("} ")
	}

	// message
	buf.WriteString(entry.Message)

	buf.WriteByte('\n')
	return buf.Bytes(), nil
}
