package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func createDB(dbname string) {
	db, err := sql.Open("sqlite3", dbname)
	if err != nil {
		log.Fatal(err)
	}
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
func importTags(db sql.DB, tags string) {
	// explode tags by comma
	// perform insert of the tags i they dont exist
}
func importWords(db sql.DB, tags string, filename string) {

	// TODO: Open file
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare("insert into words(name) values(?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	for i := 0; i < 100; i++ {
		_, err = stmt.Exec(fmt.Sprintf("%03d", i))
		if err != nil {
			log.Fatal(err)
		}
	}
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
}

func searchWords(db sql.DB, tags string) {
	rows, err := db.Query("select id, name from words")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var name string
		err = rows.Scan(&id, &name)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(id, name)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	dbPtr := flag.String("db", "foo.db", "the database filename")
	tagsPtr := flag.String("tags", "", "a comma separated list of the tags to use")

	flag.Parse()

	fmt.Println("using db:", *dbPtr)
	fmt.Println("using tags:", *tagsPtr)

	// FIXME: check fo db existence first
	file, err := os.Open(*dbPtr)
	file.Close()
	if errors.Is(err, os.ErrNotExist) {
		fmt.Println("database does not exist, creating...")
		createDB(*dbPtr)
	}
	db, err := sql.Open("sqlite3", *dbPtr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if flag.NArg() == 0 {
		fmt.Println("no filename given, performing a search")
		searchWords(*db, *tagsPtr)
	} else {
		fmt.Println("performing an import on the given files:", flag.Args())
		for i := 0; i < flag.NArg(); i++ {
			fmt.Println("importing:", flag.Arg(i))
			importWords(*db, *tagsPtr, flag.Arg(i))
		}
	}

	os.Exit(0)
}
