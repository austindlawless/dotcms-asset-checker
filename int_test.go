package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"testing"
)

var _ = fmt.Print // For debugging; delete when done.
var _ = log.Print // For debugging; delete when done.

func setup() (Config, *MySql, chan string, chan error) {
	config, err := getConfig(Config{}, "test.yaml")

	if err != nil {
		panic(err)
	}

	os.RemoveAll(config.Assets)
	os.MkdirAll(config.Assets, 0755)

	mysql := NewMySql(config.User, config.Pass, config.Host, config.Db)

	fsQueue := make(chan string)
	doneWorkSig := make(chan error, 1)

	return config, mysql, fsQueue, doneWorkSig
}

func TestFileValidAssetsCheck(t *testing.T) {
	config, _, fsQueue, doneWorkSig := setup()

	go CheckAssets(fsQueue, doneWorkSig)

	// setup dirs + files on fs
	os.Mkdir(config.Assets+"/1", 0755)
	os.Mkdir(config.Assets+"/2", 0755)
	os.Create(config.Assets + "/1/asdf.txt")
	os.Create(config.Assets + "/2/asdf.txt")

	// add dirs to queue
	fsQueue <- config.Assets + "/1"
	fsQueue <- config.Assets + "/2"

	close(fsQueue)

	err := <-doneWorkSig

	if err != nil {
		t.Error("This should have passed")
	}
}

func TestEmptyDirAssetsCheck(t *testing.T) {
	config, _, fsQueue, doneWorkSig := setup()

	go CheckAssets(fsQueue, doneWorkSig)

	// setup dirs + files on fs
	os.Mkdir(config.Assets+"/1", 0755)
	os.Mkdir(config.Assets+"/2", 0755)
	os.Create(config.Assets + "/1/asdf.txt")

	// add dir + empty dir to queue
	fsQueue <- config.Assets + "/1"
	fsQueue <- config.Assets + "/2"

	close(fsQueue)

	err := <-doneWorkSig

	if err == nil {
		t.Error("There should have been an error found")
	}
}

func TestExtractCreation(t *testing.T) {
	config, _, fsQueue, doneWorkSig := setup()

	go CreateBackupExtract(config, fsQueue, doneWorkSig)

	fsQueue <- config.Assets + "/1"
	fsQueue <- config.Assets + "/2"

	close(fsQueue)

	err := <-doneWorkSig

	if err != nil {
		t.Error("This should have passed")
	}

	contents := getExtractContents(config)

	if contents != config.Assets+"/1"+config.Assets+"/2" {
		t.Error("Weird file contents")
	}
}

func getExtractContents(config Config) string {
	var contents string

	files, _ := os.Open(config.BackupStoragePath)
	defer files.Close()

	scanner := bufio.NewScanner(files)
	for scanner.Scan() {
		contents += scanner.Text()
	}

	return contents
}
