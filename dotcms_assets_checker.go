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
	workerErrors := make(chan error, 1)
	errors := make(chan error, 1)

	// Process command
	switch *cmd {
	case "checkdatabase":
		go AssetsFromDatabase(mysql, config, filesQueue, workerErrors)
		go CheckAssets(config, filesQueue, errors)

	case "checkextract":
		go AssetsFromExtract(config, filesQueue, workerErrors)
		go CheckAssets(config, filesQueue, errors)

	case "genextract":
		go AssetsFromDatabase(mysql, config, filesQueue, workerErrors)
		go CreateBackupExtract(config, filesQueue, errors)

	default:
		panic("-cmd is requred. options: checkdatabase, checkextract, genextract")
	}

	// Wait for done signals before exiting
	err := <-errors
	wrkErr := <-workerErrors

	// critical errors
	if wrkErr != nil {
		log.Println(err)
		os.Exit(2)
	}

	// analysis errors
	if err != nil {
		log.Println(err)
		os.Exit(1)
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
