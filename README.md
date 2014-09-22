# dotCMS Assets Checker

[![Build Status](https://drone.io/github.com/BancVue/dotcms-asset-checker/status.png)](https://drone.io/github.com/BancVue/dotcms-asset-checker/latest)

Commands
--------

- checkdatabase: checks the assets against a running mysql instance
- checkextract: checks assets using a generated extract; requires the -backupstoragepath flag/yaml property
- genextract: generates an extract (newlined file paths) of assets using the database; requires the -backupstoragepath flag/yaml property; -backupstoragepath folder must already exist

Config Props
------------
Config props can either be added to a yaml and passed in via the -config /path/to/config.yaml flag or passed in via their associated CLI flag

	host: mysqlhost
	db: dotcms
	user: mysqluser
	pass: mysqlpass
	assets: /dotcms/assets/backup/123123125412
	backupstoragepath: /tmp/dotcms/asset_paths.txt

Examples
--------

	# Check assets via running MySql
	$ ./dotcms-assets-checker -config default.yaml -cmd checkdatabase

	# Generate Assets Extract
	$ mkdir -p /tmp/dotcms/backup
	$ ./dotcms-assets-checker -config default.yaml -cmd genextract -backupstoragepath /tmp/dotcms/backup/assets.txt

	# Check Existing Assets Extract
	$ ./dotcms-assets-checker -config default.yaml -cmd checkextract -backupstoragepath /tmp/dotcms/backup/assets.txt


Tests
-----

	$ make test
