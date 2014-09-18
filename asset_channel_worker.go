package main

import (
	"bufio"
	"log"
	"os"
)

func AssetsFromDatabase(mysql *MySql, config Config, assets chan string, dSig chan bool) {
	var channelWorker = &AssetChannelWorker{MySql: mysql, Config: config, FileChannel: assets, DoneSignal: dSig}

	channelWorker.ReadFromDatabase()
}

func AssetsFromExtract(config Config, assets chan string, dSig chan bool) {
	var channelWorker = &AssetChannelWorker{MySql: nil, Config: config, FileChannel: assets, DoneSignal: dSig}

	channelWorker.ReadFromFileSystem()
}

type AssetChannelWorker struct {
	MySql       *MySql
	Config      Config
	FileChannel chan string
	DoneSignal  chan bool
}

func (f *AssetChannelWorker) ReadFromDatabase() {
	var err error

	fields, err := f.MySql.db.Query("SELECT CONCAT('" + f.Config.Assets + "', '/', SUBSTRING(c.inode, 1, 1), '/', SUBSTRING(c.inode, 2, 1), '/', c.inode, '/', f.velocity_var_name) AS assetPath " +
		"FROM contentlet c " +
		"INNER JOIN contentlet_version_info ci ON ci.identifier=c.identifier AND ci.working_inode=c.inode " +
		"INNER JOIN field f ON c.structure_inode=f.structure_inode " +
		"WHERE f.field_type IN ('binary', 'file') AND f.`required`=1 " +
		"ORDER BY c.inode;")

	defer fields.Close()

	if err != nil {
		log.Println(err)
		os.Exit(2)
	}

	var assets_folder string

	for fields.Next() {
		fields.Scan(&assets_folder)

		f.FileChannel <- assets_folder
	}

	close(f.FileChannel)

	f.DoneSignal <- true
}

func (f *AssetChannelWorker) ReadFromFileSystem() {
	file, err := os.Open(f.Config.BackupStoragePath)
	defer file.Close()

	if err != nil {
		log.Println(err)
		os.Exit(2)
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		f.FileChannel <- scanner.Text()
	}

	close(f.FileChannel)

	f.DoneSignal <- true
}
