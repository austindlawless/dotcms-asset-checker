# dotCMS Assets Checker

Commands
--------

- checkdatabase: checks the assets against a running mysql instance
- checkextract: checks assets using a generated extract; requires the -backupstoragepath flag/yaml property
- genextract: generates an extract (newlined file paths) of assets using the database; requires the -backupstoragepath flag/yaml property

Config Props
------------
Config props can either be added to a yaml and passed in via the -config /path/to/config.yaml flag or passed in via their associated CLI flag

```yaml
host: mysqlhost
db: dotcms
user: mysqluser
pass: mysqlpass
assets: /var/bv/apps/dotcms/assets/backup/123123125412
backupstoragepath: /tmp/dotcms/asset_paths.txt
```