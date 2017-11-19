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
	log.Info("nomal info test ...")
	log.Debug("nomal debug test ...")

	log.Close()
}
