package main

import (
	"fmt"
	// "github.com/stretchr/testify/mock"
	"log"
	"os"
	"testing"
)

var _ = fmt.Print // For debugging; delete when done.
var _ = log.Print // For debugging; delete when done.

func setup() (Config, *MySql) {
	config, err := getConfig(Config{}, "test.yaml")

	if err != nil {
		panic(err)
	}

	os.RemoveAll(config.Assets)
	os.MkdirAll(config.Assets, 0755)

	mysql := NewMySql(config.User, config.Pass, config.Host, config.Db)

	return config, mysql
}

func TestFileValidAssetsCheck(t *testing.T) {
	config, _ := setup()

	fsQueue := make(chan string)
	doneWorkSig := make(chan error, 1)

	var channelChecker = &AssetChannelChecker{FileChannel: fsQueue, DoneSignal: doneWorkSig}
	go channelChecker.CheckFiles()

	queueFile(config.Assets+"/somefile1.txt", fsQueue, true)
	queueFile(config.Assets+"/somefile2.txt", fsQueue, true)
	queueFile(config.Assets+"/somefile3.txt", fsQueue, true)

	close(fsQueue)

	<-doneWorkSig
}

func TestFileInvalidAssetsCheck(t *testing.T) {
	config, _ := setup()

	fsQueue := make(chan string)
	doneWorkSig := make(chan error, 1)

	var channelChecker = &AssetChannelChecker{FileChannel: fsQueue, DoneSignal: doneWorkSig}
	go channelChecker.CheckFiles()

	queueFile(config.Assets+"/somefile1.txt", fsQueue, true)
	queueFile(config.Assets+"/somefile2.txt", fsQueue, false)
	queueFile(config.Assets+"/somefile3.txt", fsQueue, true)

	close(fsQueue)

	err := <-doneWorkSig

	log.Println(err)

	if err == nil {
		t.Error("There should have been an error found")
	}
}

func TestExtractCreation(t *testing.T) {
}

func queueFile(file string, queue chan string, create bool) {
	if create {
		os.Create(file)
	}

	queue <- file
}
