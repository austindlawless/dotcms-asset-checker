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
	fsQueue := make(chan string)
	doneSig := make(chan bool, 1)
	doneWorkSig := make(chan bool, 1)

	// Channel processors
	var channelWorker = &AssetChannelWorker{MySql: mysql, Config: config, FileChannel: fsQueue, DoneSignal: doneSig}
	var channelFsbackup = &AssetChannelFsbackup{Config: config, FileChannel: fsQueue, DoneSignal: doneWorkSig}
	var channelChecker = &AssetChannelChecker{FileChannel: fsQueue, DoneSignal: doneWorkSig}

	// Process command
	switch *cmd {
	case "checkdatabase":
		go channelWorker.ReadFromDatabase()
		go channelChecker.CheckFiles()

	case "checkextract":
		go channelWorker.ReadFromFileSystem()
		go channelChecker.CheckFiles()

	case "genextract":
		channelFsbackup.Init()

		go channelWorker.ReadFromDatabase()
		go channelFsbackup.BackupFiles()
	default:
		panic("-cmd is requred. options: checkdatabase, checkextract, genextract")
	}

	// Wait for done signals before exiting
	<-doneWorkSig
	<-doneSig
}

func processConfig() Config {
	yamlPath := flag.String("config", "default.yaml", "path to yaml config")

	host := flag.String("host", "", "mysql host")
	db := flag.String("db", "", "mysql user")
	user := flag.String("user", "", "mysql user")
	pass := flag.String("pass", "", "mysql password")
	logPath := flag.String("log", "", "log path")
	assets := flag.String("assets", "", "dotcms assets path")
	backupStoragePath := flag.String("backupStoragePath", "", "path to place backup file")

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

	// Config setup
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
