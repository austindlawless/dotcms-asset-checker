package main

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
)

// go routine
func CheckAssets(config Config, assets chan string, dSig chan error) {
	checker := &AssetChannelChecker{Config: config, FileChannel: assets, DoneSignal: dSig}

	checker.CheckFiles()
}

type AssetChannelChecker struct {
	FileChannel chan string
	DoneSignal  chan error
	Config      Config
}

func (f *AssetChannelChecker) CheckFiles() {
	var err error

	isValid := true

	for dir := range f.FileChannel {
		absDir := f.Config.Assets + "/" + dir

		exists, err := f.exists(absDir)

		if err != nil {
			err = err
		}

		if !exists {
			log.Println("MISSING: " + absDir)
			isValid = false
		} else {
			files, _ := ioutil.ReadDir(absDir)
			if len(files) <= 0 {
				log.Println("EMPTY DIR: " + absDir)
				isValid = false
			}
		}
	}

	if !isValid {
		err = errors.New("Missing files found")
	}

	f.DoneSignal <- err
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
