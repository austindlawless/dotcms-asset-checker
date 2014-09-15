package main

import (
	"fmt"
	// "github.com/stretchr/testify/mock"
	"log"
	"os"
	"testing"
)

var _ = fmt.Print // For debugging; delete when done.
var _ = log.Print // For debugging; delete when done.

// wtf "cannot use NewMySql(config.User, config.Pass, config.Host, config.Db) (type *MySql) as type MySql in assignment"
// var mysql MySql
var config Config

func setup() {
	var err error
	config, err = getConfig(Config{}, "test.yaml")

	if err != nil {
		panic(err)
	}

	os.RemoveAll(config.Assets)

	// mysql = NewMySql(config.User, config.Pass, config.Host, config.Db)
	// defer mysql.Close()
}

func TestFileAssetsCheck(t *testing.T) {
	setup()

	mysql := NewMySql(config.User, config.Pass, config.Host, config.Db)
	defer mysql.Close()

	setupTestAssets(mysql)

	check := AssetsCheck{MySql: mysql, AssetsPath: config.Assets}

	valid, err := check.Check()

	if !valid {
		t.Error("Something was missing")
	}

	if err != nil {
		t.Error(err)
	}
}

func TestInvalidFileAssetsCheck(t *testing.T) {
	setup()

	mysql := NewMySql(config.User, config.Pass, config.Host, config.Db)
	defer mysql.Close()

	setupTestAssets(mysql)

	os.RemoveAll(config.Assets + "/0/1")

	check := AssetsCheck{MySql: mysql, AssetsPath: config.Assets}

	valid, _ := check.Check()

	if valid {
		t.Error("Check didn't find missing items")
	}
}

func setupTestAssets(mysql *MySql) {
	fields, _ := mysql.db.Query("SELECT f.structure_inode, f.velocity_var_name FROM field f " +
		"JOIN structure s ON s.inode = f.structure_inode " +
		"WHERE f.field_type IN ('binary', 'image', 'file') AND s.structuretype=4 ORDER BY f.structure_inode;")

	defer fields.Close()

	var structure_inode string
	var assets_folder string
	var inode string
	var file_name string

	for fields.Next() {
		fields.Scan(&structure_inode, &assets_folder)

		contentlets, _ := mysql.db.Query("SELECT cl.inode, cl.text3 AS assetToCheck "+
			"FROM contentlet cl "+
			"JOIN contentlet_version_info c ON c.identifier=cl.identifier AND c.working_inode=cl.inode "+
			"WHERE structure_inode=?", structure_inode)

		defer contentlets.Close()

		for contentlets.Next() {
			contentlets.Scan(&inode, &file_name)

			if file_name != "" {
				path := config.Assets + "/" + inode[0:1] + "/" + inode[1:2] + "/" + inode + "/" + assets_folder

				os.MkdirAll(path, 0775)
				os.Create(path + "/" + file_name)
			}
		}
	}
}
