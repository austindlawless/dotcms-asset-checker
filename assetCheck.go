package main

import (
	"log"
	"os"
)

type AssetsCheck struct {
	MySql *MySql
}

func (f *AssetsCheck) Check() (error) {
	log.Println("Assets Checking")

	rows, err := f.MySql.db.Query("SELECT c.inode " +
									"FROM contentlet c" +
									"JOIN field f ON c.structure_inode=f.structure_inode" +
									"WHERE f.field_type IN ('binary', 'image', 'file');")

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	var inode string
	var path string

	for rows.Next() {
		rows.Scan(&inode)

		// @todo make this a config param
		path = "/var/bv/apps/dotcms/assets/" + inode[0:1] + "/" + inode[1:2] + "/" + inode

		// log.Println(inode + " -> " + path)

		exixsts, _ := exists(path)

		if err != nil {
			log.Println(err)
			os.Exit(1)
		}

		if exixsts == true {
			log.Println("Inode exists")
		} else {
			log.Println("NOT FOUND: " + inode)
		}
	}

	log.Println("Done")

	return nil
}

func exists(path string) (bool, error) {
    _, err := os.Stat(path)

    if err == nil {
    	return true, nil
    }

    if os.IsNotExist(err) { 
    	return false, nil 
    }

    return false, err
}