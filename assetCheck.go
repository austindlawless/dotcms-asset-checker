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

	rows, err := f.MySql.db.Query("SELECT firstname, lastname FROM user_;")

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	var firstname string
	var lastname string

	for rows.Next() {
		rows.Scan(&firstname, &lastname)

		log.Println(firstname + " " + lastname)
	}

	return nil
}