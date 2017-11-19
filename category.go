package log4go

import (
	"fmt"
	"runtime"
	"strings"
	"time"
)

// LOGGER get the log Filter by category
func LOGGER(category string) *Filter {
	f, ok := Global[category]
	if !ok {
		f = Global["stdout"]
		f.Category = "DEFAULT"
	} else {
		f.Category = category
	}
	return f
}

// Send a formatted log message internally
func (f *Filter) intLogf(lvl Level, format string, args ...interface{}) {
	skip := true

	// Determine if any logging will be done
	if lvl >= f.Level {
		skip = false
	}
	if skip {
		return
	}

	// Determine caller func
	pc, _, lineno, ok := runtime.Caller(2)
	src := ""
	if ok {
		src = fmt.Sprintf("%s:%d", runtime.FuncForPC(pc).Name(), lineno)
	}

	msg := format
	if len(args) > 0 {
		msg = fmt.Sprintf(format, args...)
	}

	// Make the log record
	rec := &LogRecord{
		Level:    lvl,
		Created:  time.Now(),
		Source:   src,
		Message:  msg,
		Category: f.Category,
	}

	// Dispatch the logs
	/*for _, filt := range log {
		if lvl < filt.Level {
			continue
		}
		filt.LogWrite(rec)
	}
	*/
	default_filter := Global["stdout"]

	if lvl > default_filter.Level {
		default_filter.LogWrite(rec)
	}

	if f.Category != "DEFAULT" && f.Category != "stdout" {
		f.LogWrite(rec)
	}

}

// Send a closure log message internally
func (f *Filter) intLogc(lvl Level, closure func() string) {
	skip := true

	// Determine if any logging will be done
	if lvl >= f.Level {
		skip = false
	}
	if skip {
		return
	}

	// Determine caller func
	pc, _, lineno, ok := runtime.Caller(2)
	src := ""
	if ok {
		src = fmt.Sprintf("%s:%d", runtime.FuncForPC(pc).Name(), lineno)
	}

	// Make the log record
	rec := &LogRecord{
		Level:    lvl,
		Created:  time.Now(),
		Source:   src,
		Message:  closure(),
		Category: f.Category,
	}

	default_filter := Global["stdout"]

	if lvl > default_filter.Level {
		default_filter.LogWrite(rec)
	}

	if f.Category != "DEFAULT" && f.Category != "stdout" {
		f.LogWrite(rec)
	}
}

// Send a log message with manual level, source, and message.
func (f *Filter) Log(lvl Level, source, message string) {
	skip := true

	// Determine if any logging will be done
	if lvl >= f.Level {
		skip = false
	}
	if skip {
		return
	}

	// Make the log record
	rec := &LogRecord{
		Level:    lvl,
		Created:  time.Now(),
		Source:   source,
		Message:  message,
		Category: f.Category,
	}

	default_filter := Global["stdout"]

	if lvl > default_filter.Level {
		default_filter.LogWrite(rec)
	}

	if f.Category != "DEFAULT" && f.Category != "stdout" {
		f.LogWrite(rec)
	}
}

// Logf logs a formatted log message at the given log level, using the caller as
// its source.
func (f *Filter) Logf(lvl Level, format string, args ...interface{}) {
	f.intLogf(lvl, format, args...)
}

// Logc logs a string returned by the closure at the given log level, using the caller as
// its source.  If no log message would be written, the closure is never called.
func (f *Filter) Logc(lvl Level, closure func() string) {
	f.intLogc(lvl, closure)
}

// Finest logs a message at the finest log level.
// See Debug for an explanation of the arguments.
func (f *Filter) Finest(arg0 interface{}, args ...interface{}) {
	const (
		lvl = FINEST
	)
	switch first := arg0.(type) {
	case string:
		// Use the string as a format string
		f.intLogf(lvl, first, args...)
	case func() string:
		// Log the closure (no other arguments used)
		f.intLogc(lvl, first)
	default:
		// Build a format string so that it will be similar to Sprint
		f.intLogf(lvl, fmt.Sprint(arg0)+strings.Repeat(" %v", len(args)), args...)
	}
}

// Fine logs a message at the fine log level.
// See Debug for an explanation of the arguments.
func (f *Filter) Fine(arg0 interface{}, args ...interface{}) {
	const (
		lvl = FINE
	)
	switch first := arg0.(type) {
	case string:
		// Use the string as a format string
		f.intLogf(lvl, first, args...)
	case func() string:
		// f the closure (no other arguments used)
		f.intLogc(lvl, first)
	default:
		// Build a format string so that it will be similar to Sprint
		f.intLogf(lvl, fmt.Sprint(arg0)+strings.Repeat(" %v", len(args)), args...)
	}
}

// Debug is a utility method for debug f messages.
// The behavior of Debug depends on the first argument:
// - arg0 is a string
//   When given a string as the first argument, this behaves like ff but with
//   the DEBUG f level: the first argument is interpreted as a format for the
//   latter arguments.
// - arg0 is a func()string
//   When given a closure of type func()string, this fs the string returned by
//   the closure iff it will be fged.  The closure runs at most one time.
// - arg0 is interface{}
//   When given anything else, the f message will be each of the arguments
//   formatted with %v and separated by spaces (ala Sprint).
func (f *Filter) Debug(arg0 interface{}, args ...interface{}) {
	const (
		lvl = DEBUG
	)
	switch first := arg0.(type) {
	case string:
		// Use the string as a format string
		f.intLogf(lvl, first, args...)
	case func() string:
		// f the closure (no other arguments used)
		f.intLogc(lvl, first)
	default:
		// Build a format string so that it will be similar to Sprint
		f.intLogf(lvl, fmt.Sprint(arg0)+strings.Repeat(" %v", len(args)), args...)
	}
}

// Trace fs a message at the trace f level.
// See Debug for an explanation of the arguments.
func (f *Filter) Trace(arg0 interface{}, args ...interface{}) {
	const (
		lvl = TRACE
	)
	switch first := arg0.(type) {
	case string:
		// Use the string as a format string
		f.intLogf(lvl, first, args...)
	case func() string:
		// f the closure (no other arguments used)
		f.intLogc(lvl, first)
	default:
		// Build a format string so that it will be similar to Sprint
		f.intLogf(lvl, fmt.Sprint(arg0)+strings.Repeat(" %v", len(args)), args...)
	}
}

// Info fs a message at the info f level.
// See Debug for an explanation of the arguments.
func (f *Filter) Info(arg0 interface{}, args ...interface{}) {
	const (
		lvl = INFO
	)
	switch first := arg0.(type) {
	case string:
		// Use the string as a format string
		f.intLogf(lvl, first, args...)
	case func() string:
		// f the closure (no other arguments used)
		f.intLogc(lvl, first)
	default:
		// Build a format string so that it will be similar to Sprint
		f.intLogf(lvl, fmt.Sprint(arg0)+strings.Repeat(" %v", len(args)), args...)
	}
}

// Warn fs a message at the warning f level and returns the formatted error.
// At the warning level and higher, there is no performance benefit if the
// message is not actually fged, because all formats are processed and all
// closures are executed to format the error message.
// See Debug for further explanation of the arguments.
func (f *Filter) Warn(arg0 interface{}, args ...interface{}) {
	const (
		lvl = WARNING
	)
	var msg string
	switch first := arg0.(type) {
	case string:
		// Use the string as a format string
		msg = fmt.Sprintf(first, args...)
	case func() string:
		// f the closure (no other arguments used)
		msg = first()
	default:
		// Build a format string so that it will be similar to Sprint
		msg = fmt.Sprintf(fmt.Sprint(first)+strings.Repeat(" %v", len(args)), args...)
	}
	f.intLogf(lvl, msg)
}

// Error fs a message at the error f level and returns the formatted error,
// See Warn for an explanation of the performance and Debug for an explanation
// of the parameters.
func (f *Filter) Error(arg0 interface{}, args ...interface{}) {
	const (
		lvl = ERROR
	)
	var msg string
	switch first := arg0.(type) {
	case string:
		// Use the string as a format string
		msg = fmt.Sprintf(first, args...)
	case func() string:
		// f the closure (no other arguments used)
		msg = first()
	default:
		// Build a format string so that it will be similar to Sprint
		msg = fmt.Sprintf(fmt.Sprint(first)+strings.Repeat(" %v", len(args)), args...)
	}
	f.intLogf(lvl, msg)
}

// Critical fs a message at the critical f level and returns the formatted error,
// See Warn for an explanation of the performance and Debug for an explanation
// of the parameters.
func (f *Filter) Critical(arg0 interface{}, args ...interface{}) {
	const (
		lvl = CRITICAL
	)
	var msg string
	switch first := arg0.(type) {
	case string:
		// Use the string as a format string
		msg = fmt.Sprintf(first, args...)
	case func() string:
		// f the closure (no other arguments used)
		msg = first()
	default:
		// Build a format string so that it will be similar to Sprint
		msg = fmt.Sprintf(fmt.Sprint(first)+strings.Repeat(" %v", len(args)), args...)
	}
	f.intLogf(lvl, msg)
}
