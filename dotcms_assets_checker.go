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

	flag.Parse()

	flagConfig := Config{
		Host:   *host,
		User:   *user,
		Db:     *db,
		Assets: *assets,
		Pass:   *pass,
		Log:    *logPath,
	}

	config, err := getConfig(flagConfig, *yamlPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	if config.Log != "" {
		f, err := os.OpenFile(config.Log, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			fmt.Println("error opening file: %v", err)
			os.Exit(2)
		}
		log.SetOutput(f)
	}

	mysql := NewMySql(config.User, config.Pass, config.Host, config.Db)
	defer mysql.Close()

	var checker = &AssetsCheck{MySql: mysql, AssetsPath: config.Assets}

	isValid, err := checker.Check()
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	if !isValid {
		os.Exit(1)
	}

}
