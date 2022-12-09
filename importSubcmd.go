package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

// parse args of the import subcommand and exec it
func importSubcmd(args []string) {
	importCmd := flag.NewFlagSet("import", flag.ExitOnError)
	importCmd.SetOutput(flag.CommandLine.Output())

	importCmd.Usage = func() {
		fmt.Fprint(importCmd.Output(), "Import a word file into the database with optional tags\n\n")
		fmt.Fprintln(importCmd.Output(), "Usage of orunmila import:")
		fmt.Fprintf(importCmd.Output(), "orunmila [-db <db_path>] [-debug] import -tags [tag]... [filesname]...\n\n")
		fmt.Fprintln(importCmd.Output(), "  filename\n\tthe filename(s) to read the words from")
		importCmd.PrintDefaults()
	}

	var (
		tagsPtr = importCmd.String("tags", "", "a comma separated list of the tags to use")
	)

	importCmd.Parse(args)

	if len(args) == 0 {
		log.Error("[importSubcmd] you need to provide at least a filename")
		importCmd.Usage()
		os.Exit(1)
	}

	log.Println("performing an import on the given files:", importCmd.Args())

	dsn := fmt.Sprintf("file:%s?mode=rw", *dbPtr)
	db, err := sql.Open("sqlite3", dsn)
	check(err)
	defer db.Close()

	Tags = stringToArray(*tagsPtr)

	for i := 0; i < importCmd.NArg(); i++ {
		log.Println("[importSubcmd] importing file:", importCmd.Arg(i))
		if isFileExists(importCmd.Arg(i)) {
			importFileWords(db, importCmd.Arg(i))
		} else {
			log.Warnf("[importSubcmd] %q does not exists.", importCmd.Arg(i))
		}
	}
}
