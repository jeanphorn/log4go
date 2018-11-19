# log4go

## Description

This repository is reconstructed from alecthomas's log4go, which is a logging package similar to log4j for the Go programming language.

Two new features are supported: one is Json config style, and the other is different output according to category.

## Features

-   **Log to console**
-   **Log to file, support rotate by size or time**
-   **log to network, support tcp and udp**
-   **support xml config**

---------------------------

-   **Support Json style configuration**
-   **Add Category for log**
    * Classify your logs for different output and different usage.
-   **Compatible with the old**
-   **Support json style config content beside filename**

## Usage

First, get the code from this repo. 

```go get github.com/jeanphorn/log4go```

Then import it to your project.

```import log "github.com/jeanphorn/log4go"```


## Examples

The config file is optional, if you don't set the config file, it would use the default console config.

Here is a Json config example:

```
{
    "console": {
        "enable": true,		// wether output the log
        "level": "FINE"		// log level: FINE, DEBUG, TRACE, INFO, WARNING,ERROR, CRITICAL
    },  
    "files": [{
        "enable": true,
        "level": "DEBUG",
        "filename":"./test.log",
        "category": "Test",			// different category log to different files
        "pattern": "[%D %T] [%C] [%L] (%S) %M"	// log output formmat
    },{ 
        "enable": false,
        "level": "DEBUG",
        "filename":"rotate_test.log",
        "category": "TestRotate",
        "pattern": "[%D %T] [%C] [%L] (%S) %M",
        "rotate": true,				// whether rotate the log
        "maxsize": "500M",
        "maxlines": "10K",
        "daily": true,
        "sanitize": true
    }], 
    "sockets": [{
        "enable": false,
        "level": "DEBUG",
        "category": "TestSocket",
        "pattern": "[%D %T] [%C] [%L] (%S) %M",
        "addr": "127.0.0.1:12124",
        "protocol":"udp"
    }]  
}
```

Code example:

```
package main

import (
	log "github.com/jeanphorn/log4go"
)

func main() {
	// load config file, it's optional
	// or log.LoadConfiguration("./example.json", "json")
	// config file could be json or xml
	log.LoadConfiguration("./example.json")

	log.LOGGER("Test").Info("category Test info test ...")
	log.LOGGER("Test").Info("category Test info test message: %s", "new test msg")
	log.LOGGER("Test").Debug("category Test debug test ...")

	// Other category not exist, test
	log.LOGGER("Other").Debug("category Other debug test ...")

	// socket log test
	log.LOGGER("TestSocket").Debug("category TestSocket debug test ...")

	// original log4go test
	log.Info("normal info test ...")
	log.Debug("normal debug test ...")

	log.Close()
}

```

The output like:

> [2017/11/15 14:35:11 CST] [Test] [INFO] (main.main:15) category Test info test ...     
> [2017/11/15 14:35:11 CST] [Test] [INFO] (main.main:16) category Test info test message: new test msg     
> [2017/11/15 14:35:11 CST] [Test] [DEBG] (main.main:17) category Test debug test ...     
> [2017/11/15 14:35:11 CST] [DEFAULT] [INFO] (main.main:26) normal info test ...     
> [2017/11/15 14:35:11 CST] [DEFAULT] [DEBG] (main.main:27) normal debug test ...    


## Thanks

Thanks alecthomas for providing the [original resource](https://github.com/alecthomas/log4go).
