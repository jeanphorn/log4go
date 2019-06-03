// Copyright (C) 2010, Kyle Lemons <kyle@kylelemons.net>.  All rights reserved.

package log4go

import (
	"fmt"
	"io"
	"os"
	"time"
)

// This is the standard writer that prints to standard output.
type ConsoleLogWriter struct {
	format string
	w      chan *LogRecord
}

// This creates a new ConsoleLogWriter
func NewConsoleLogWriter(stderr bool) *ConsoleLogWriter {
	consoleWriter := &ConsoleLogWriter{
		format: "[%T %D] [%C] [%L] (%S) %M",
		w:      make(chan *LogRecord, LogBufferLength),
	}
	out := os.Stdout
	if stderr {
		out = os.Stderr
	}
	go consoleWriter.run(out)
	return consoleWriter
}
func (c *ConsoleLogWriter) SetFormat(format string) {
	c.format = format
}
func (c *ConsoleLogWriter) run(out io.Writer) {
	for rec := range c.w {
		fmt.Fprint(out, FormatLogRecord(c.format, rec))
	}
}

// This is the ConsoleLogWriter's output method.  This will block if the output
// buffer is full.
func (c *ConsoleLogWriter) LogWrite(rec *LogRecord) {
	c.w <- rec
}

// Close stops the logger from sending messages to standard output.  Attempts to
// send log messages to this logger after a Close have undefined behavior.
func (c *ConsoleLogWriter) Close() {
	close(c.w)
	time.Sleep(50 * time.Millisecond) // Try to give console I/O time to complete
}
