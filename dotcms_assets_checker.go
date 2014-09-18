package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

var _ = log.Print

func main() {
	cmd := flag.String("cmd", "", "dotcms assets path")

	// Get the config
	config := processConfig()

	// Setup mysql
	mysql := NewMySql(config.User, config.Pass, config.Host, config.Db)
	defer mysql.Close()

	// Setup channels
	filesQueue := make(chan string)
	doneSig := make(chan bool, 1)
	doneWorkSig := make(chan error, 1)

	// Process command
	switch *cmd {
	case "checkdatabase":
		go AssetsFromDatabase(mysql, config, filesQueue, doneSig)
		go CheckAssets(filesQueue, doneWorkSig)

	case "checkextract":
		go AssetsFromExtract(config, filesQueue, doneSig)
		go CheckAssets(filesQueue, doneWorkSig)

	case "genextract":
		go AssetsFromDatabase(mysql, config, filesQueue, doneSig)
		go CreateBackupExtract(config, filesQueue, doneWorkSig)

	default:
		panic("-cmd is requred. options: checkdatabase, checkextract, genextract")
	}

	// Wait for done signals before exiting
	err := <-doneWorkSig
	<-doneSig

	// asserts error types ect...
	if err != nil {
		log.Println(err)
		os.Exit(2)
	}
}

func processConfig() Config {
	yamlPath := flag.String("config", "default.yaml", "path to yaml config")

	host := flag.String("host", "", "mysql host")
	db := flag.String("db", "", "mysql user")
	user := flag.String("user", "", "mysql user")
	pass := flag.String("pass", "", "mysql password")
	logPath := flag.String("log", "", "log path")
	assets := flag.String("assets", "", "dotcms assets path")
	backupStoragePath := flag.String("backupstoragepath", "", "path to place backup file")

	flag.Parse()

	flagConfig := Config{
		Host:              *host,
		User:              *user,
		Db:                *db,
		Assets:            *assets,
		Pass:              *pass,
		Log:               *logPath,
		BackupStoragePath: *backupStoragePath,
	}

	config, err := getConfig(flagConfig, *yamlPath)
	checkError(err)

	if config.Log != "" {
		f, err := os.OpenFile(config.Log, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

		checkError(err)

		log.SetOutput(f)
	}

	return config
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
}
