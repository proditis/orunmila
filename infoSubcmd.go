package main

import (
	"database/sql"
	"flag"
	"fmt"

	log "github.com/sirupsen/logrus"
)

// parse args of the info subcommand and exec it
func infoSubcmd(args []string) {
	infoCmd := flag.NewFlagSet("info", flag.ExitOnError)

	infoCmd.Usage = func() {
		fmt.Fprint(infoCmd.Output(), "Display database system configuration information\n\n")
		fmt.Fprintf(infoCmd.Output(), "Usage of orunmila info:\n")
		infoCmd.PrintDefaults()
		fmt.Fprintln(infoCmd.Output(), "\texample: orunmila [-db <db_path>] info")
	}

	// nothing to parse, just there to trigger the usage menu
	infoCmd.Parse(args)

	dsn := fmt.Sprintf("file:%s?mode=rw", *dbPtr)
	db, err := sql.Open("sqlite3", dsn)
	check(err)
	defer db.Close()

	rows, err := db.Query("SELECT * FROM sysconfig")
	check(err)
	defer rows.Close()

	for rows.Next() {
		var name string
		var val string
		err = rows.Scan(&name, &val)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("[%s]: %s\n", name, val)
	}
	err = rows.Err()
	check(err)
}
