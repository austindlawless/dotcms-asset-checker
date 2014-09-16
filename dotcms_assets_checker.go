package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

var _ = log.Print

func main() {
	yamlPath := flag.String("config", "default.yaml", "path to yaml config")

	host := flag.String("host", "", "mysql host")
	db := flag.String("db", "", "mysql user")
	user := flag.String("user", "", "mysql user")
	pass := flag.String("pass", "", "mysql password")
	logPath := flag.String("log", "", "log path")
	assets := flag.String("assets", "", "dotcms assets path")
	backupStoragePath := flag.String("backupStoragePath", "", "path to place backup file")
	cmd := flag.String("cmd", "", "dotcms assets path")

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

	mysql := NewMySql(config.User, config.Pass, config.Host, config.Db)
	defer mysql.Close()

	switch *cmd {
	case "checkdatabase":
		// get files from db -> queue
		// check files from queue
		checkFiles(config, mysql, "db")
	case "checkextract":
		// read files from fs -> queue
		// check files from queue
		checkFiles(config, mysql, "fs")
	case "genextract":
		// get files from db -> queue
		// write queue to fs
		generateBackup(config, mysql)
	default:
		panic("-cmd is requred. options: checkdatabase, checkextract, genextract")
	}

}

func generateBackup(config Config, mysql *MySql) {
	fsQueue := make(chan string)
	doneSig := make(chan bool, 1)
	doneWorkSig := make(chan bool, 1)

	var channelWorker = &AssetChannelWorker{MySql: mysql, Config: config, FileChannel: fsQueue, DoneSignal: doneSig}
	var channelFsbackup = &AssetChannelFsbackup{Config: config, FileChannel: fsQueue, DoneSignal: doneWorkSig}

	channelFsbackup.Init()

	go channelWorker.ReadFromDatabase()
	go channelFsbackup.BackupFiles()

	<-doneWorkSig
	<-doneSig
}

func checkFiles(config Config, mysql *MySql, check string) {
	fsQueue := make(chan string)
	doneSig := make(chan bool, 1)
	doneWorkSig := make(chan bool, 1)

	var channelWorker = &AssetChannelWorker{MySql: mysql, Config: config, FileChannel: fsQueue, DoneSignal: doneSig}
	var channelChecker = &AssetChannelChecker{FileChannel: fsQueue, DoneSignal: doneWorkSig}

	if check == "fs" {
		go channelWorker.ReadFromFileSystem()
	} else {
		go channelWorker.ReadFromDatabase()
	}

	go channelChecker.CheckFiles()

	<-doneWorkSig
	<-doneSig
}

func check(isValid bool, err error) {
	checkError(err)

	if !isValid {
		os.Exit(1)
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
}
