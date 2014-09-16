package main

import (
	"io/ioutil"
	"log"
	"os"
)

type AssetChannelFsbackup struct {
	Config      Config
	FileChannel chan string
	DoneSignal  chan bool
}

func (f *AssetChannelFsbackup) Init() {
	// Attempt to remove if exists and create folder
	os.RemoveAll(f.Config.BackupStoragePath)

	// Create file via config
	_, err := os.Create(f.Config.BackupStoragePath)
	if err != nil {
		log.Println(err)
		os.Exit(2)
	}
}

func (f *AssetChannelFsbackup) BackupFiles() error {
	var file_contents string

	// Read queue & store files in memory
	for file := range f.FileChannel {
		file_contents += file
	}

	// Write files to storage file
	d1 := []byte(file_contents)
	err := ioutil.WriteFile(f.Config.BackupStoragePath, d1, 0644)
	if err != nil {
		log.Println(err)
		os.Exit(2)
	}

	f.DoneSignal <- true

	return nil
}
