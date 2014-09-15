package main

import (
	"bufio"
	"log"
	"os"
)

type AssetChecker struct {
	Config Config
}

func (f *AssetChecker) Check() (bool, error) {
	var isValid = true

	file, err := os.Open(f.Config.BackupStoragePath)
	defer file.Close()

	if err != nil {
		return false, err
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		exists, err := f.exists(scanner.Text())

		if err != nil {
			return false, err
		}

		if !exists {
			log.Println("Missing: " + scanner.Text())
			isValid = false
		}
	}

	return isValid, nil
}

func (f *AssetChecker) exists(path string) (bool, error) {
	_, err := os.Stat(path)

	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}
