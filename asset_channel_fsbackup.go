package main

import (
	"io/ioutil"
	"log"
	"os"
)

func CreateBackupExtract(config Config, assets chan string, dSig chan error) {
	backup := &AssetChannelFsbackup{Config: config, FileChannel: assets, DoneSignal: dSig}

	backup.Init()

	backup.BackupFiles()
}

type AssetChannelFsbackup struct {
	Config      Config
	FileChannel chan string
	DoneSignal  chan error
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

func (f *AssetChannelFsbackup) BackupFiles() {
	var file_contents string

	// Read queue & store files in memory
	for file := range f.FileChannel {
		file_contents += file + "\n"
	}

	// Write files to storage file
	d1 := []byte(file_contents)
	err := ioutil.WriteFile(f.Config.BackupStoragePath, d1, 0644)

	f.DoneSignal <- err
}
