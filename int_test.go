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

func TestFileAssetsCheck(t *testing.T) {
	config, err := getConfig(Config{}, "test.yaml")

	if err != nil {
		t.Error(err)
	}

	mysql := NewMySql(config.User, config.Pass, config.Host, config.Db)
	defer mysql.Close()

	check := AssetsCheck{MySql: mysql, AssetsPath: config.Assets}

	_, err = check.Check()

	if err != nil {
		t.Error(err)
	}
}

func TestFileExistsOnDisk(t *testing.T) {
	config, err := getConfig(Config{}, "test.yaml")

	if err != nil {
		t.Error(err)
	}

	mysql := NewMySql(config.User, config.Pass, config.Host, config.Db)
	defer mysql.Close()

	check := AssetsCheck{MySql: mysql, AssetsPath: config.Assets}

	var file = "/tmp/testasset.txt"

	os.Create(file)

	exists, err := check.exists(file)

	if !exists {
		t.Error("File not found")
	}

	if err != nil {
		t.Error(err)
	}
}
