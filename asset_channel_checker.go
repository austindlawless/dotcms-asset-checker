package main

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
)

// go routine
func CheckAssets(assets chan string, dSig chan error) {
	checker := &AssetChannelChecker{FileChannel: assets, DoneSignal: dSig}

	checker.CheckFiles()
}

type AssetChannelChecker struct {
	FileChannel chan string
	DoneSignal  chan error
}

func (f *AssetChannelChecker) CheckFiles() {
	var err error

	isValid := true

	for dir := range f.FileChannel {
		exists, err := f.exists(dir)

		if err != nil {
			err = err
		}

		if !exists {
			log.Println("MISSING: " + dir)
			isValid = false
		}

		files, _ := ioutil.ReadDir(dir)
		if len(files) <= 0 {
			log.Println("EMPTY DIR: " + dir)
			isValid = false
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
