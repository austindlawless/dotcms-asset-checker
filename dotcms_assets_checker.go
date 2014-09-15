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
		generateBackup(config, mysql)
		checkFiles(config)
	case "checkextract":
		// read files from fs -> queue
		// check files from queue
		checkFiles(config)
	case "genextract":
		// get files from db -> queue
		// write queue to fs
		generateBackup(config, mysql)
	default:
		panic("-cmd is requred. options: dbcheck, backupcheck, genbackup")
	}

}

func generateBackup(config Config, mysql *MySql) {
	var generator = &GenerateBackup{MySql: mysql, Config: config}

	isValid, err := generator.MakeFile()

	check(isValid, err)
}

func checkFiles(config Config) {
	var checker = &AssetChecker{Config: config}

	isValid, err := checker.Check()

	check(isValid, err)
}

func doDbCheck(config Config, mysql *MySql) {
	var checker = &AssetsCheck{MySql: mysql, AssetsPath: config.Assets}

	isValid, err := checker.Check()

	check(isValid, err)
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
