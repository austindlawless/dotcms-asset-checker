package main

import (
	"log"
	"os"
)

type AssetChannelChecker struct {
	FileChannel chan string
	DoneSignal  chan bool
}

func (f *AssetChannelChecker) CheckFiles() {
	var isValid = true

	for file := range f.FileChannel {
		exists, err := f.exists(file)

		if err != nil {
			log.Println(err)
			os.Exit(2)
		}

		if !exists {
			log.Println("MISSING: " + file)
			isValid = false
		}
	}

	f.DoneSignal <- true

	if !isValid {
		os.Exit(1)
	}
}

func (f *AssetChannelChecker) exists(path string) (bool, error) {
	_, err := os.Stat(path)

	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}
