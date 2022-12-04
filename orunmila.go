package main

import (
	"bufio"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"

	_ "github.com/mattn/go-sqlite3"
)

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}
func createDB(dbname string) {
	db, err := sql.Open("sqlite3", dbname)
	check(err)
	defer db.Close()

	sqlStmt := `
	create table words (id integer not null primary key AUTOINCREMENT, name text NOT NULL UNIQUE);
	create table tags (id integer not null primary key AUTOINCREMENT, name text NOT NULL UNIQUE);
	create table wt (word_id integer not null , tag_id integer not null, FOREIGN KEY(word_id) REFERENCES words(id),FOREIGN KEY(tag_id) REFERENCES tags(id),PRIMARY KEY(word_id,tag_id));
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}
}

// Convert tags string into tags hash array
func tagsToArray(tags string) {
	// explode tags by comma
	// for each unique item push into a hash array the trimmed value ie tagsArray[tag]
}

// Import the tags into the database
func importTags(db sql.DB, tags string) {
	// ourArray=tagsToarray(tags)
	// perform insert of the tags if they dont exist and return the ID of that tag into the corresponding hash array
	// if the tag exists fetch its ID into the corresponding hash array
	// return the ID's
}

// Import the words from a given filename into the database
func importWords(db sql.DB, tags string, filename string) {

	file, err := os.Open(filename)
	if errors.Is(err, os.ErrNotExist) {
		log.Fatalln(err)
	}
	defer file.Close()

	tx, err := db.Begin()
	check(err)
	wordsStmt, err := tx.Prepare("insert or ignore into words(name) values(?)")
	check(err)
	defer wordsStmt.Close()

	tagsStmt, err := tx.Prepare("insert into tags(name) values(?)")
	check(err)
	defer tagsStmt.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		word := strings.TrimSpace(scanner.Text())
		if word != "" {
			fmt.Println("importing word:", word)
			_, err = wordsStmt.Exec(word)
			check(err)
			db.
		}
	}
	err = tx.Commit()
	check(err)
}

//
// Search for words matching tags
//
func searchWords(db sql.DB, tags string) {
	rows, err := db.Query("select name from words")
	check(err)
	defer rows.Close()
	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(name)
	}
	err = rows.Err()
	check(err)
}

func main() {
	dbPtr := flag.String("db", "orunmila.db", "the database filename (default: orunmila.db)")
	tagsPtr := flag.String("tags", "", "a comma separated list of the tags to use")
	wordsPtr := flag.String("words", "", "a comma separated list of words to look for")

	flag.Parse()

	fmt.Println("using db:", *dbPtr)
	fmt.Println("using tags:", *tagsPtr)
	fmt.Println("using tags:", *wordsPtr)

	// check if db file exists
	file, err := os.Open(*dbPtr)
	file.Close()
	if errors.Is(err, os.ErrNotExist) {
		log.Debugln("database does not exist, creating...")
		createDB(*dbPtr)
	}
	db, err := sql.Open("sqlite3", *dbPtr)
	check(err)
	defer db.Close()

	if flag.NArg() == 0 {
		fmt.Println("no filename given, performing a search")
		searchWords(*db, *tagsPtr)
	} else {
		fmt.Println("performing an import on the given files:", flag.Args())
		for i := 0; i < flag.NArg(); i++ {
			fmt.Println("importing file:", flag.Arg(i))
			importWords(*db, *tagsPtr, flag.Arg(i))
		}
	}

	os.Exit(0)
}
