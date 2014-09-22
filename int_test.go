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
	errors := make(chan error, 1)

	return config, mysql, fsQueue, errors
}

func TestFileValidAssetsCheck(t *testing.T) {
	config, _, fsQueue, errors := setup()

	go CheckAssets(config, fsQueue, errors)

	// setup dirs + files on fs
	os.MkdirAll(config.Assets+"/1/0", 0755)
	os.MkdirAll(config.Assets+"/2/4", 0755)
	os.Create(config.Assets + "/1/0/asdf.txt")
	os.Create(config.Assets + "/2/4/asdf.txt")

	// add dirs to queue
	fsQueue <- "1/0"
	fsQueue <- "2/4"

	close(fsQueue)

	err := <-errors

	if err != nil {
		t.Error("This should have passed")
	}
}

func TestEmptyDirAssetsCheck(t *testing.T) {
	config, _, fsQueue, errors := setup()

	go CheckAssets(config, fsQueue, errors)

	// setup dirs + files on fs
	os.MkdirAll(config.Assets+"/1/0", 0755)
	os.MkdirAll(config.Assets+"/2/4", 0755)
	os.Create(config.Assets + "/1/0/asdf.txt")

	// add dir + empty dir to queue
	fsQueue <- "1/0"
	fsQueue <- "2/4"

	close(fsQueue)

	err := <-errors

	if err == nil {
		t.Error("There should have been an error found")
	}
}

func TestExtractCreation(t *testing.T) {
	config, _, fsQueue, errors := setup()

	go CreateBackupExtract(config, fsQueue, errors)

	fsQueue <- "1/2/image"
	fsQueue <- "2/4/fileAsset"

	close(fsQueue)

	err := <-errors

	if err != nil {
		t.Error("This should have passed")
	}

	contents := getExtractContents(config)

	if contents != "1/2/image"+"2/4/fileAsset" {
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
