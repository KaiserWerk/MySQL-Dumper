package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/JamesStewy/go-mysqldump"
	"log"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

/*
  - implement max. number of retained backups
  - give a running instance a unique ID, like 'md-hostname-node'
  - send backup to remote server
  - get restore commands to restore a local backup
*/

func main() {
	confFilePtr := flag.String("config", "", "The configuration file to read")
	flag.Parse()

	if *confFilePtr != "" {
		setConfigFile(*confFilePtr)
	}

	config, created, err := setupConfig()
	if err != nil {
		log.Println("could not set up configuration", err.Error())
		return
	}

	if created {
		log.Println("configuration file created; exiting")
		return
	}

	fn, err := createDump(config)
	if err != nil {
		log.Println("could not create first dump")
		return
	}
	fmt.Printf("created first backup with name '%s'\n", fn)

	t := time.NewTicker(time.Duration(config.BackupInterval) * time.Minute)
	for {
		select {
		case <-t.C:
			if fn, err = createDump(config); err != nil {
				log.Println("could not create dump:", err.Error())
			} else {
				fmt.Printf("Backup file '%s' created\n", fn)
			}
		}
	}
}

func createDump(conf *appConfig) (string, error) {
	parts := strings.Split(conf.DSN, "/")
	dumpFilenameFormat := fmt.Sprintf("%s_2006-01_02-15_04_05", parts[len(parts)-1])

	db, err := sql.Open("mysql", conf.DSN)
	if err != nil {
		return "", fmt.Errorf("error opening database: %s", err.Error())
	}

	dumper, err := mysqldump.Register(db, conf.BackupPath, dumpFilenameFormat)
	if err != nil {
		return "", fmt.Errorf("error registering database: %s", err)
	}

	resultFilename, err := dumper.Dump()
	if err != nil {
		return "", fmt.Errorf("error dumping: %s", err.Error())
	}

	return resultFilename, dumper.Close()
}
