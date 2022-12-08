package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

// parse args of the describe subcommand and exec it
func describeSubcmd(args []string) {

	descibeCmd := flag.NewFlagSet("describe", flag.ExitOnError)

	descibeCmd.Usage = func() {
		descibeCmd.SetOutput(flag.CommandLine.Output())
		fmt.Fprint(descibeCmd.Output(), "Set the database description\n\n")
		fmt.Fprintf(descibeCmd.Output(), "Usage of orunmila descibe:\n")
		fmt.Fprintf(descibeCmd.Output(), "orunmila describe My Awesome Description\n\n")
		fmt.Fprintln(descibeCmd.Output(), "  words strings\n\tdescription")
		descibeCmd.PrintDefaults()
	}

	descibeCmd.Parse(args)

	if len(args) < 1 {
		log.Errorln("please provide a description")
		flag.Usage()
		os.Exit(1)
	}

	dsn := fmt.Sprintf("file:%s?mode=rw", *dbPtr)
	db, err := sql.Open("sqlite3", dsn)
	check(err)
	defer db.Close()

	tx, err := db.Begin()
	check(err)

	descStmt, err := tx.Prepare("INSERT OR REPLACE INTO sysconfig(name,val) values (?,?)")
	check(err)
	defer descStmt.Close()

	desc := strings.TrimSpace(strings.Join(args, " "))
	_, err = descStmt.Exec("description", desc)
	check(err)

	err = tx.Commit()
	check(err)
}
