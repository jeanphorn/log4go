package log4go

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/toolkits/file"
)

type ConsoleConfig struct {
	Enable  bool   `json:"enable"`
	Level   string `json:"level"`
	Pattern string `json:"pattern"`
}

type FileConfig struct {
	Enable   bool   `json:"enable"`
	Category string `json:"category"`
	Level    string `json:"level"`
	Filename string `json:"filename"`

	// %T - Time (15:04:05 MST)
	// %t - Time (15:04)
	// %D - Date (2006/01/02)
	// %d - Date (01/02/06)
	// %L - Level (FNST, FINE, DEBG, TRAC, WARN, EROR, CRIT)
	// %S - Source
	// %M - Message
	// %C - Category
	// It ignores unknown format strings (and removes them)
	// Recommended: "[%D %T] [%C] [%L] (%S) %M"//
	Pattern string `json:"pattern"`

	Rotate   bool   `json:"rotate"`
	Maxsize  string `json:"maxsize"`  // \d+[KMG]? Suffixes are in terms of 2**10
	Maxlines string `json:"maxlines"` //\d+[KMG]? Suffixes are in terms of thousands
	Daily    bool   `json:"daily"`    //Automatically rotates by day
	Sanitize bool	`json:"sanitize"` //Sanitize newlines to prevent log injection
}

type SocketConfig struct {
	Enable   bool   `json:"enable"`
	Category string `json:"category"`
	Level    string `json:"level"`
	Pattern  string `json:"pattern"`

	Addr     string `json:"addr"`
	Protocol string `json:"protocol"`
}

// LogConfig presents json log config struct
type LogConfig struct {
	Console *ConsoleConfig  `json:"console"`
	Files   []*FileConfig   `json:"files"`
	Sockets []*SocketConfig `json:"sockets"`
}

// LoadJsonConfiguration load log config from json file
// see examples/example.json for ducumentation
func (log Logger) LoadJsonConfiguration(filename string) {
	log.Close()
	dst := new(bytes.Buffer)
	var (
		lc      LogConfig
		content string
	)
	err := json.Compact(dst, []byte(filename))

	if err != nil {
		content, err = ReadFile(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "LoadJsonConfiguration: Error: Could not read %q: %s\n", filename, err)
			os.Exit(1)
		}
	} else {
		content = string(dst.Bytes())
	}

	err = json.Unmarshal([]byte(content), &lc)
	if err != nil {
		fmt.Fprintf(os.Stderr, "LoadJsonConfiguration: Error: Could not parse json configuration in %q: %s\n", filename, err)
		os.Exit(1)
	}

	if lc.Console.Enable {
		filt, _ := jsonToConsoleLogWriter(filename, lc.Console)
		log["stdout"] = &Filter{getLogLevel(lc.Console.Level), filt, "DEFAULT"}
	}

	for _, fc := range lc.Files {
		if !fc.Enable {
			continue
		}
		if len(fc.Category) == 0 {
			fmt.Fprintf(os.Stderr, "LoadJsonConfiguration: file category can not be empty in <%s>: ", filename)
			os.Exit(1)
		}

		filt, _ := jsonToFileLogWriter(filename, fc)
		log[fc.Category] = &Filter{getLogLevel(fc.Level), filt, fc.Category}
	}

	for _, sc := range lc.Sockets {
		if !sc.Enable {
			continue
		}
		if len(sc.Category) == 0 {
			fmt.Fprintf(os.Stderr, "LoadJsonConfiguration: file category can not be empty in <%s>: ", filename)
			os.Exit(1)
		}

		filt, _ := jsonToSocketLogWriter(filename, sc)
		log[sc.Category] = &Filter{getLogLevel(sc.Level), filt, sc.Category}
	}

}

func getLogLevel(l string) Level {
	var lvl Level
	switch l {
	case "FINEST":
		lvl = FINEST
	case "FINE":
		lvl = FINE
	case "DEBUG":
		lvl = DEBUG
	case "TRACE":
		lvl = TRACE
	case "INFO":
		lvl = INFO
	case "WARNING":
		lvl = WARNING
	case "ERROR":
		lvl = ERROR
	case "CRITICAL":
		lvl = CRITICAL
	default:
		fmt.Fprintf(os.Stderr, "LoadJsonConfiguration: Error: Required level <%s> for filter has unknown value: %s\n", "level", l)
		os.Exit(1)
	}
	return lvl
}

func jsonToConsoleLogWriter(filename string, cf *ConsoleConfig) (*ConsoleLogWriter, bool) {
	format := "[%D %T] [%C] [%L] (%S) %M"

	if len(cf.Pattern) > 0 {
		format = strings.Trim(cf.Pattern, " \r\n")
	}

	if !cf.Enable {
		return nil, true
	}

	clw := NewConsoleLogWriter()
	clw.SetFormat(format)

	return clw, true
}

func jsonToFileLogWriter(filename string, ff *FileConfig) (*FileLogWriter, bool) {
	file := "app.log"
	format := "[%D %T] [%C] [%L] (%S) %M"
	maxlines := 0
	maxsize := 0
	daily := false
	rotate := false
	sanitize := false

	if len(ff.Filename) > 0 {
		file = ff.Filename
	}
	if len(ff.Pattern) > 0 {
		format = strings.Trim(ff.Pattern, " \r\n")
	}
	if len(ff.Maxlines) > 0 {
		maxlines = strToNumSuffix(strings.Trim(ff.Maxlines, " \r\n"), 1000)
	}
	if len(ff.Maxsize) > 0 {
		maxsize = strToNumSuffix(strings.Trim(ff.Maxsize, " \r\n"), 1024)
	}
	daily = ff.Daily
	rotate = ff.Rotate
	sanitize = ff.Sanitize

	if !ff.Enable {
		return nil, true
	}

	flw := NewFileLogWriter(file, rotate, daily)
	flw.SetFormat(format)
	flw.SetRotateLines(maxlines)
	flw.SetRotateSize(maxsize)
	flw.SetSanitize(sanitize)
	return flw, true
}

func jsonToSocketLogWriter(filename string, sf *SocketConfig) (SocketLogWriter, bool) {
	endpoint := ""
	protocol := "tcp"

	if len(sf.Addr) == 0 {
		fmt.Fprintf(os.Stderr, "LoadConfiguration: Error: Required property \"%s\" for file filter missing in %s\n", "addr", filename)
		os.Exit(1)
	}
	endpoint = sf.Addr

	// set socket protocol
	if len(sf.Protocol) > 0 {
		if sf.Protocol != "tcp" && sf.Protocol != "udp" {
			fmt.Fprintf(os.Stderr, "LoadConfiguration: Error: Required property \"%s\" for file filter wrong type in %s, use default tcp instead.\n", "protocol", filename)
		} else {
			protocol = sf.Protocol
		}
	}

	if !sf.Enable {
		return nil, true
	}

	return NewSocketLogWriter(protocol, endpoint), true
}

func ReadFile(path string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("[%s] path empty", path)
	}

	if !file.IsExist(path) {
		return "", fmt.Errorf("config file %s is nonexistent", path)
	}

	configContent, err := file.ToTrimString(path)
	if err != nil {
		return "", fmt.Errorf("read file %s fail %s", path, err)
	}

	return configContent, nil
}
