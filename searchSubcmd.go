package main

import (
	"database/sql"
	"flag"
	"fmt"

	log "github.com/sirupsen/logrus"
)

// parse args of the search subcommand and exec it
func searchSubcmd(args []string) error {
	searchCmd := flag.NewFlagSet("search", flag.ContinueOnError)
	ShouldReturn := false
	searchCmd.Usage = func() {
		searchCmd.SetOutput(flag.CommandLine.Output())
		fmt.Fprint(searchCmd.Output(), "Display words matching an optional list of tags\n\n")
		fmt.Fprintf(searchCmd.Output(), "Usage of orunmila search:\n")
		fmt.Fprintf(searchCmd.Output(), "orunmila [-db <db_path>] [-debug] search [-tags OPTIONAL_TAGS]\n\n")
		searchCmd.PrintDefaults()
		ShouldReturn = true
	}

	var (
		tagsPtr = searchCmd.String("tags", "", "a comma separated list of the tags to use")
	)

	err := searchCmd.Parse(args)

	if ShouldReturn {
		log.Println(err)
		return err
	}
	log.Debugln("[searchSubcmd] using db:", *dbPtr)
	log.Debugln("[searchSubcmd] using tags:", *tagsPtr)
	log.Println("[searchSubcmd] performing a search")

	Tags = stringToArray(*tagsPtr)

	dsn := fmt.Sprintf("file:%s?mode=ro", *dbPtr)
	db, err := sql.Open("sqlite3", dsn)
	check(err)
	defer db.Close()

	searchWordsByTagIds(db, *tagsPtr)
	return nil
}
