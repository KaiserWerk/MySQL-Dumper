# MySQL-Dumper
A tiny (sidecar) app that can automatically backup your app's MySQL database and optionally store it remotely.
It is supposed to be run against a central MySQL Dump instance but it can also be
run standalone if you do not need the management/restore features. 

* Autmatically creates a backup file (*.sql) from a database using the supplied DSN
  at a configurable interval (min. 1 minute)
* 