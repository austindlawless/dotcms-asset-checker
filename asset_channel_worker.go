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

	fields, err := f.MySql.db.Query("SELECT f.structure_inode, f.velocity_var_name FROM field f " +
		"JOIN structure s ON s.inode = f.structure_inode " +
		"WHERE f.field_type IN ('binary', 'image', 'file') AND s.structuretype=4 ORDER BY f.structure_inode;")

	defer fields.Close()

	if err != nil {
		log.Println(err)
		os.Exit(2)
	}

	var structure_inode string
	var assets_folder string
	var inode string
	var file_name string

	for fields.Next() {
		fields.Scan(&structure_inode, &assets_folder)

		contentlets, err := f.MySql.db.Query("SELECT cl.inode, cl.text3 AS assetToCheck "+
			"FROM contentlet cl "+
			"JOIN contentlet_version_info c ON c.identifier=cl.identifier AND c.working_inode=cl.inode "+
			"WHERE structure_inode=?", structure_inode)

		defer contentlets.Close()

		if err != nil {
			log.Println(err)
			os.Exit(2)
		}

		for contentlets.Next() {
			contentlets.Scan(&inode, &file_name)

			if file_name != "" {
				f.FileChannel <- f.Config.Assets + "/" + inode[0:1] + "/" + inode[1:2] + "/" + inode + "/" + assets_folder + "/" + file_name
			}
		}
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
