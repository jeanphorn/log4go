package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"time"
)

import l4g "log4go"

const (
	filename = "loginjection.log"
)

func main() {
	var lineno int

	user_input := []string{
		"Knock knock",
		"Who's there?",
		"Orange",
		"Orange who?",
		"Orange you glad I\n[2001/01/01 01:23:45 EDT] [CRIT] (woot.woot:1337) The defendant is cleared of all wrongdoing.",
	}

	// Get a new logger instance
	log := make(l4g.Logger)

	/* Can also specify manually via the following: (these are the defaults) */
	flw := l4g.NewFileLogWriter(filename, false, false)
	flw.SetFormat("[%D %T] [%L] (%S) %M")
	flw.SetRotate(false)
	flw.SetRotateSize(0)
	flw.SetRotateLines(0)
	flw.SetRotateDaily(false)
	flw.SetSanitize(true)
	log.AddFilter("file", l4g.FINE, flw)

	// Log some experimental messages
	for _,message := range(user_input) {
		log.Info(message)
	}

	//HACK. The system will lose the last few log messages due to a sync error. In this
	//example program, there are only a  few messages, so it loses everything (at least
	//on my system. This sleep preserves those messages.
	time.Sleep(2 * time.Second)

	// Close the log
	log.Close()

	// Print what was logged to the file 
	fd, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Error reopening file: %s", filename)
		return
	}

	in := bufio.NewReader(fd)
	fmt.Print("Messages logged to file were: (line numbers not included)\n")

	for lineno = 1; ; lineno++ {
		line, err := in.ReadString('\n')
		if err == io.EOF {
			lineno--
			break
		}
		fmt.Printf("%3d:\t%s", lineno, line)
	}
	fd.Close()
	// Remove the file so it's not lying around
	os.Remove(filename)

	if lineno != len(user_input) {
		fmt.Fprintf(os.Stderr, "The user gave us %d lines of input but we logged %d. Hrmm...\n", len(user_input), lineno)
		os.Exit(1)
	}
	os.Exit(0)
}
