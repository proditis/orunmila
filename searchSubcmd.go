package main

import (
	"database/sql"
	"flag"
	"fmt"

	log "github.com/sirupsen/logrus"
)

// parse args of the search subcommand and exec it
func searchSubcmd(args []string) {
	searchCmd := flag.NewFlagSet("search", flag.ExitOnError)

	var (
		tagsPtr = searchCmd.String("tags", "", "a comma separated list of the tags to use")
	)

	searchCmd.Parse(args)

	log.Debugln("[searchSubcmd] using db:", *dbPtr)
	log.Debugln("[searchSubcmd] using tags:", *tagsPtr)
	log.Println("[searchSubcmd] performing a search")

	Tags = stringToArray(*tagsPtr)

	dsn := fmt.Sprintf("file:%s?mode=ro", *dbPtr)
	db, err := sql.Open("sqlite3", dsn)
	check(err)
	defer db.Close()

	searchWordsByTagIds(db, *tagsPtr)
}
