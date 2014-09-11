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

	// rows, err := f.MySql.db.Query("SELECT c.inode " +
	// 								"FROM contentlet c " +
	// 								"JOIN field f ON c.structure_inode=f.structure_inode " +
	// 								"WHERE f.field_type IN ('binary', 'image', 'file');")

	fields, err := f.MySql.db.Query("SELECT f.structure_inode, f.field_type, f.velocity_var_name FROM field f " +
									 "JOIN structure s ON s.inode = f.structure_inode " +
									 "WHERE f.field_type IN ('binary', 'image', 'file') AND s.structuretype=4 ORDER BY f.structure_inode;")

	defer fields.Close()

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	var inode string
	var assetToCheck string
	var structure_inode string
	var field_type string
	var field_contentlet string

	for fields.Next() {
		fields.Scan(&structure_inode, &field_type, &field_contentlet)

		// if field_type == "binary" {
		// 	continue
		// }

		// log.Println(structure_inode + ", " + field_type + ", " + field_contentlet)		
		// log.Println("Checking " + structure_inode + "....")
		// log.Println("SELECT inode, "+field_contentlet+" AS assetToCheck FROM contentlet WHERE structure_inode='"+structure_inode+"'")

		// contentlets, err := f.MySql.db.Query("SELECT inode, "+field_contentlet+" AS assetToCheck FROM contentlet WHERE structure_inode=?", structure_inode)
		contentlets, err := f.MySql.db.Query("SELECT inode, text3 AS assetToCheck FROM contentlet WHERE structure_inode=?", structure_inode)

		defer contentlets.Close()

		if err != nil {
			log.Println(err)
			os.Exit(1)
		}

		for contentlets.Next() {
			contentlets.Scan(&inode, &assetToCheck)

			if assetToCheck != "" {
				path := "/var/bv/apps/dotcms/assets/" + inode[0:1] + "/" + inode[1:2] + "/" + inode + "/" + field_contentlet + "/" + assetToCheck

				exixsts, _ := exists(path)

				if exixsts == true {
					log.Println("Inode exists")
				} else {
					log.Println("NOT FOUND! Contentlet: " + inode)
				}

				// log.Println("	" + path)
			}
		}

		// for contentlets
		// @todo make this a config param
		// path = "/var/bv/apps/dotcms/assets/" + inode[0:1] + "/" + inode[1:2] + "/" + inode

		// // log.Println(inode + " -> " + path)

		// exixsts, _ := exists(path)

		// if err != nil {
		// 	log.Println(err)
		// 	os.Exit(1)
		// }

		// if exixsts == true {
		// 	log.Println("Inode exists")
		// } else {
		// 	log.Println("NOT FOUND: " + inode)
		// }
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