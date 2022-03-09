package ttfm

// Adapted from https://github.com/antonfisher/nested-logrus-formatter

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	"github.com/sirupsen/logrus"
)

// LogFormatter - logrus formatter, implements logrus.Formatter
type LogFormatter struct {
	// NoColors - disable colors
	NoColors bool

	// NoFieldsColors - apply colors only to the level, default is level + fields
	NoFieldsColors bool

	// TrimMessages - trim whitespaces on messages
	TrimMessages bool
}

// Format an log entry
func (f *LogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	levelColor := getColorByLevel(entry.Level)

	// output buffer
	b := &bytes.Buffer{}

	// timestamp
	b.WriteString(entry.Time.Format("2006-01-02 15:04:05"))

	// message
	if f.TrimMessages {
		b.WriteString(" " + strings.TrimSpace(entry.Message))
	} else {
		b.WriteString(" " + entry.Message)
	}

	// level
	var level string
	level = strings.ToUpper(entry.Level.String())

	if !f.NoColors {
		fmt.Fprintf(b, "\x1b[%dm", levelColor)
	}

	b.WriteString(" [")
	b.WriteString(level[:4])
	b.WriteString("] ")

	if !f.NoColors && f.NoFieldsColors {
		b.WriteString("\x1b[0m")
	}

	// fields
	f.writeFields(b, entry)
	b.WriteString(" ")

	if !f.NoColors && !f.NoFieldsColors {
		b.WriteString("\x1b[0m")
	}

	f.writeCaller(b, entry)

	b.WriteByte('\n')

	return b.Bytes(), nil
}

func (f *LogFormatter) writeCaller(b *bytes.Buffer, entry *logrus.Entry) {
	if entry.HasCaller() {
		fmt.Fprintf(
			b,
			" (%s:%d %s)",
			entry.Caller.File,
			entry.Caller.Line,
			entry.Caller.Function,
		)
	}
}

func (f *LogFormatter) writeFields(b *bytes.Buffer, entry *logrus.Entry) {
	if len(entry.Data) != 0 {
		fields := make([]string, 0, len(entry.Data))
		for field := range entry.Data {
			fields = append(fields, field)
		}

		sort.Strings(fields)

		for _, field := range fields {
			f.writeField(b, entry, field)
		}
	}
}

func (f *LogFormatter) writeField(b *bytes.Buffer, entry *logrus.Entry, field string) {
	fmt.Fprintf(b, "%s=\"%v\"", field, entry.Data[field])
	b.WriteString(" ")
}

const (
	colorRed    = 31
	colorYellow = 33
	colorBlue   = 36
	colorGray   = 37
)

func getColorByLevel(level logrus.Level) int {
	switch level {
	case logrus.DebugLevel, logrus.TraceLevel:
		return colorGray
	case logrus.WarnLevel:
		return colorYellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		return colorRed
	default:
		return colorBlue
	}
}
