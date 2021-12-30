package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/JamesStewy/go-mysqldump"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var (
	dsn string = ""
)

func main() {
	dsnPtr := flag.String("dsn", "", "The DSN to use for the MySQL connection")
	flag.Parse()

	if d := os.Getenv("MYSQLDUMPER_DSN"); d != "" {
		dsn = d
	} else if *dsnPtr == "" {
		dsn = *dsnPtr
	}

	if dsn == "" {
		log.Println("DSN is not set!")
		return
	}

}

func createDump() error {
	dumpDir := "dumps"                                              // you should create this directory
	dumpFilenameFormat := fmt.Sprintf("%s-20060102T150405", dbname) // accepts time layout string and add .sql at the end of file

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", username, password, hostname, port, dbname))
	if err != nil {
		fmt.Println("Error opening database: ", err)
		return
	}

	// Register database with mysqldump
	dumper, err := mysqldump.Register(db, dumpDir, dumpFilenameFormat)
	if err != nil {
		fmt.Println("Error registering databse:", err)
		return
	}

	// Dump database to file
	resultFilename, err := dumper.Dump()
	if err != nil {
		fmt.Println("Error dumping:", err)
		return
	}
	fmt.Printf("File is saved to %s", resultFilename)

	// Close dumper and connected database
	dumper.Close()
}
