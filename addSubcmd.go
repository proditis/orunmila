package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

// parse args of the add subcommand and exec it
func addSubcmd(args []string) {
	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	addCmd.SetOutput(flag.CommandLine.Output())
	addCmd.Usage = func() {
		addCmd.SetOutput(flag.CommandLine.Output())
		fmt.Fprint(addCmd.Output(), "Add words to the database from the command line with optional tags\n\n")
		fmt.Fprintf(addCmd.Output(), "Usage of orunmila add:\n")
		fmt.Fprintf(addCmd.Output(), "orunmila [-db <db_path>] [-debug] add [-tags OPTIONAL_TAGS] words to add\n\n")
		fmt.Fprintln(addCmd.Output(), "  words strings\n\tspace separated words to add")
		addCmd.PrintDefaults()
	}

	var (
		tagsPtr = addCmd.String("tags", "", "a comma separated list of the tags to use")
	)

	addCmd.Parse(args)

	if len(args) == 0 {
		log.Error("[addSubcmd] you need to provide words to be added")
		addCmd.Usage()
		os.Exit(1)
	}

	log.Infoln("[addSubcmd] Adding the given words:", addCmd.Args())

	dsn := fmt.Sprintf("file:%s?mode=rw", *dbPtr)
	db, err := sql.Open("sqlite3", dsn)
	check(err)
	defer db.Close()

	Tags = stringToArray(*tagsPtr)
	importTags(db)

	tx, err := db.Begin()
	check(err)

	wordsStmt, err := tx.Prepare("insert or ignore into words(name) values(?)")
	check(err)
	defer wordsStmt.Close()

	wtStmt, err := tx.Prepare("insert or ignore into wt(word_id,tag_id) values(?,?)")
	check(err)
	defer wtStmt.Close()

	for i := 0; i < addCmd.NArg(); i++ {
		log.Println("[addSubcmd] adding word:", addCmd.Arg(i))
		word := strings.TrimSpace(addCmd.Arg(i))
		var word_id int64
		if word != "" {
			log.Debugln("[addSubcmd] importing word:", word)
			if word_id = getWordId(db, word); word_id <= 0 {
				result, err := wordsStmt.Exec(word)
				check(err)
				word_id, err = result.LastInsertId()
				check(err)
				if word_id <= 0 {
					log.Debugf("[addSubcmd] word %s already exists, fetching", word)
					word_id = getWordId(db, word)
					log.Println("Found word id:", word_id)
				}
			}

			log.Debugf("[addSubcmd] word: %s => id: %d\n", word, word_id)
			for tag, tag_id := range Tags {
				log.Debugf("[addSubcmd] adding wt(%d,%d) // %s %s", word_id, tag_id, word, tag)
				_, err = wtStmt.Exec(word_id, tag_id)
				if err != nil {
					log.Error(err)
				}
			}
			word_id = 0
		}
	}
	err = tx.Commit()
	check(err)
}
