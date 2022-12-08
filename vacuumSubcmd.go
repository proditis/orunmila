package main

import (
	"database/sql"
	"flag"
	"fmt"

	log "github.com/sirupsen/logrus"
)

// parse args of the vaccum subcommand and exec it
func vacuumSubcmd(args []string) {
	vacuumCmd := flag.NewFlagSet("vacuum", flag.ExitOnError)

	vacuumCmd.Usage = func() {
		fmt.Fprint(vacuumCmd.Output(), "Rebuild the database file, repacking it into a minimal amount of disk space\n\n")
		fmt.Fprintf(vacuumCmd.Output(), "Usage of orunmila vacuum:\n")
		vacuumCmd.PrintDefaults()
		fmt.Fprintln(vacuumCmd.Output(), "\texample: orunmila [-db <db_path>] vacuum")
	}

	// nothing to parse, just there to trigger the usage menu
	vacuumCmd.Parse(args)

	dsn := fmt.Sprintf("file:%s?mode=rw", *dbPtr)
	db, err := sql.Open("sqlite3", dsn)
	check(err)

	defer db.Close()
	_, err = db.Query("VACUUM")
	check(err)

	log.Println("database rebuilt successfully")
}
