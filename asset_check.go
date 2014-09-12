package main

import (
	"log"
	"os"
)

type AssetsCheck struct {
	MySql      *MySql
	AssetsPath string
}

func (f *AssetsCheck) Check() error {
	fields, err := f.MySql.db.Query("SELECT f.structure_inode, f.field_type, f.velocity_var_name FROM field f " +
		"JOIN structure s ON s.inode = f.structure_inode " +
		"WHERE f.field_type IN ('binary', 'image', 'file') AND s.structuretype=4 ORDER BY f.structure_inode;")

	defer fields.Close()

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	var structure_inode string
	var field_type string
	var field_contentlet string

	for fields.Next() {
		fields.Scan(&structure_inode, &field_type, &field_contentlet)

		f.validateContentlets(structure_inode, field_contentlet)
	}

	return nil
}

func (f *AssetsCheck) validateContentlets(structure_inode string, asset_folder_name string) error {
	var inode string
	var assetToCheck string

	// Only select working nodes
	contentlets, err := f.MySql.db.Query("SELECT cl.inode, cl.text3 AS assetToCheck "+
		"FROM contentlet cl "+
		"JOIN contentlet_version_info c ON c.identifier=cl.identifier AND c.working_inode=cl.inode "+
		"WHERE structure_inode=?", structure_inode)

	defer contentlets.Close()

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	for contentlets.Next() {
		contentlets.Scan(&inode, &assetToCheck)

		if assetToCheck != "" {
			path := f.AssetsPath + "/" + inode[0:1] + "/" + inode[1:2] + "/" + inode + "/" + asset_folder_name + "/" + assetToCheck

			exixsts, _ := f.exists(path)

			if !exixsts {
				log.Println("NOT FOUND! Contentlet: " + inode + ", " + path)
			}
		}
	}

	return nil
}

func (f *AssetsCheck) exists(path string) (bool, error) {
	_, err := os.Stat(path)

	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}
