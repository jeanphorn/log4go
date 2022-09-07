// Copyright (C) 2010, Kyle Lemons <kyle@kylelemons.net>.  All rights reserved.

package log4go

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strings"
)

const (
	FORMAT_DEFAULT = "[%Y-%m-%d %H:%M:%S.%c] [%L] (%A) %I"
	FORMAT_SHORT   = "[%m-%d %H:%M:%S] [%L] %I"
	FORMAT_ABBREV  = "[%L] %I"
)

type formatCacheType struct {
	LastUpdateSeconds      int64
	tMicroSecond           string
	tLongYear, tShortYear  string
	tMonth, tDay, tHour    string
	tMinute, tSecond       string
	tTimeZone              string
	tWeekday, tWeekdayName string
	tMonthName             string
	longDate               string
}

var formatCache = &formatCacheType{}

// FormatLogRecord ------------- Self Definitions ---------------
//	Commonly used format codes:
//
//	%Y  Year with century as a decimal number.
//	%m  Month as a decimal number [01,12].
//	%d  Day of the month as a decimal number [01,31].
//	%H  Hour (24-hour clock) as a decimal number [00,23].
//	%M  Minute as a decimal number [00,59].
//	%S  Second as a decimal number [00,61].
//	%o  Microseconds as a decimal number [000000,999999].
//	%z  Time zone offset from UTC.
//	%w  Locale's abbreviated weekday name.
//	%W  Locale's full weekday name.
//	%b  Locale's abbreviated month name.
//	%N  Locale's full month name.
//	%I  Received message.
//	%A  Source.
//	%C  field Category string of LogRecord type.
//	%L	Level (FNST, FINE, DEBG, TRAC, WARN, EROR, CRIT)
//	%D	Custom format: %D{2006-01-02T15:04:05}
//	Ignores unknown formats
//	Recommended: "[%m-%d %H:%M:%S] [%L] %I"
//	%p  Locale's equivalent of either AM or PM.
// ------------------------------------------------
func FormatLogRecord(format string, rec *LogRecord) string {
	if rec == nil {
		return "<nil>"
	}
	if len(format) == 0 {
		return ""
	}

	out := bytes.NewBuffer(make([]byte, 0, 64))
	secs := rec.Created.UnixNano() / 1e9

	cache := *formatCache
	if cache.LastUpdateSeconds != secs {
		month, day, year := rec.Created.Month(), rec.Created.Day(), rec.Created.Year()
		hour, minute, second := rec.Created.Hour(), rec.Created.Minute(), rec.Created.Second()
		monthName, weekday := rec.Created.Month().String(), rec.Created.Weekday()
		millisecond := rec.Created.UnixMicro() % 1e6
		zone, _ := rec.Created.Zone()
		updated := &formatCacheType{
			LastUpdateSeconds: secs,
			longDate:          fmt.Sprintf("%04d-%02d-%02d", year, month, day),
			tLongYear:         fmt.Sprintf("%04d", year),
			tShortYear:        fmt.Sprintf("%02d", year%100),
			tMonth:            fmt.Sprintf("%02d", month),
			tDay:              fmt.Sprintf("%02d", day),
			tHour:             fmt.Sprintf("%02d", hour),
			tMinute:           fmt.Sprintf("%02d", minute),
			tSecond:           fmt.Sprintf("%02d", second),
			tMicroSecond:      fmt.Sprintf("%06d", millisecond),
			tTimeZone:         fmt.Sprintf("%s", zone),
			tMonthName:        fmt.Sprintf("%s", monthName),
			tWeekday:          fmt.Sprintf("%s", weekday),
		}
		cache = *updated
		formatCache = updated
	}
	//custom format datetime pattern %D{2006-01-02T15:04:05}
	formatByte := changeDttmFormat(format, rec)
	// Split the string into pieces by % signs
	pieces := bytes.Split(formatByte, []byte{'%'})

	// Iterate over the pieces, replacing known formats
	for i, piece := range pieces {
		if i > 0 && len(piece) > 0 {
			switch piece[0] {
			case 'D':
				out.WriteString(cache.longDate)
			case 'Y':
				out.WriteString(cache.tLongYear)
			case 'y':
				out.WriteString(cache.tShortYear)
			case 'm':
				out.WriteString(cache.tMonth)
			case 'd':
				out.WriteString(cache.tDay)
			case 'H':
				out.WriteString(cache.tHour)
			case 'M':
				out.WriteString(cache.tMinute)
			case 'S':
				out.WriteString(cache.tSecond)
			case 'o':
				out.WriteString(cache.tMicroSecond)
			case 'z':
				out.WriteString(cache.tTimeZone)
			case 'N':
				out.WriteString(cache.tMonthName)
			case 'W':
				out.WriteString(cache.tWeekday)
			case 'A':
				out.WriteString(rec.Source)
			case 'a':
				slice := strings.Split(rec.Source, "/")
				out.WriteString(slice[len(slice)-1])
			case 'L':
				out.WriteString(levelStrings[rec.Level])
			case 'I':
				out.WriteString(rec.Message)
			case 'C':
				if len(rec.Category) == 0 {
					rec.Category = "DEFAULT"
				}
				out.WriteString(rec.Category)
			}
			if len(piece) > 1 {
				out.Write(piece[1:])
			}
		} else if len(piece) > 0 {
			out.Write(piece)
		}
	}
	out.WriteByte('\n')

	return out.String()
}

// This is the standard writer that prints to standard output.
type FormatLogWriter chan *LogRecord

// This creates a new FormatLogWriter
func NewFormatLogWriter(out io.Writer, format string) FormatLogWriter {
	records := make(FormatLogWriter, LogBufferLength)
	go records.run(out, format)
	return records
}

func (w FormatLogWriter) run(out io.Writer, format string) {
	defer recoverPanic()
	for rec := range w {
		fmt.Fprint(out, FormatLogRecord(format, rec))
	}
}

// This is the FormatLogWriter's output method.  This will block if the output
// buffer is full.
func (w FormatLogWriter) LogWrite(rec *LogRecord) {
	w <- rec
}

// Close stops the logger from sending messages to standard output.  Attempts to
// send log messages to this logger after a Close have undefined behavior.
func (w FormatLogWriter) Close() {
	close(w)
}

func changeDttmFormat(format string, rec *LogRecord) []byte {
	formatByte := []byte(format)
	r := regexp.MustCompile("\\%D\\{(.*?)\\}")
	i := 0
	formatByte = r.ReplaceAllFunc(formatByte, func(s []byte) []byte {
		if i < 2 {
			i++
			str := string(s)
			str = strings.Replace(str, "%D", "", -1)
			str = strings.Replace(str, "{", "", -1)
			str = strings.Replace(str, "}", "", -1)
			return []byte(rec.Created.Format(str))
		}
		return s
	})
	return formatByte
}
