package main

import (
	"io/ioutil"
	"os"
)

type GenerateBackup struct {
	MySql  *MySql
	Config Config
}

func (f *GenerateBackup) MakeFile() (bool, error) {
	var storageFile = f.Config.BackupStoragePath + "/asset_paths.txt"

	// Get files from Db
	file_data, err := f.getFiles()
	if err != nil {
		return false, err
	}

	// Attempt to remove if exists and create folder
	os.RemoveAll(f.Config.BackupStoragePath)
	err = os.MkdirAll(f.Config.BackupStoragePath, 0775)
	if err != nil {
		return false, err
	}

	// Create file via config
	_, err = os.Create(storageFile)
	if err != nil {
		return false, err
	}

	// Write files to storage file
	d1 := []byte(file_data)
	err = ioutil.WriteFile(storageFile, d1, 0644)
	if err != nil {
		return false, err
	}

	// Worked!
	return true, nil
}

func (f *GenerateBackup) getFiles() (string, error) {
	fields, err := f.MySql.db.Query("SELECT f.structure_inode, f.velocity_var_name FROM field f " +
		"JOIN structure s ON s.inode = f.structure_inode " +
		"WHERE f.field_type IN ('binary', 'image', 'file') AND s.structuretype=4 ORDER BY f.structure_inode;")

	defer fields.Close()

	if err != nil {
		return "", err
	}

	var file_data string
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
			return "", err
		}

		for contentlets.Next() {
			contentlets.Scan(&inode, &file_name)

			if file_name != "" {
				file_data += f.Config.Assets + "/" + inode[0:1] + "/" + inode[1:2] + "/" + inode + "/" + assets_folder + "/" + file_name + "\n"
			}
		}
	}

	return file_data, nil
}
