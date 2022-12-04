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

var Tags = make(map[string]int64)

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func TagsToString() string {
	var a []string
	for _, id := range Tags {
		a = append(a, fmt.Sprint(id))
	}
	return strings.Join(a, ",")
}

// Gets the ID of a given tag
func getTagId(db sql.DB, tag string) int64 {
	var id int64
	stmt, err := db.Prepare("select id from tags where name = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	stmt.QueryRow(tag).Scan(&id)
	return id
}

// Gets the ID of a given word
func getWordId(db sql.DB, word string) int64 {
	var id int64
	stmt, err := db.Prepare("select id from words where name = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	stmt.QueryRow(word).Scan(&id)
	return id
}

// Creates the database schema
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

// Converts tags string into tags hash array
func tagsToArray(tags string) map[string]int64 {
	// explode tags by comma
	var tagsArray = strings.Split(tags, ",")

	// loop through unique items
	for _, s := range tagsArray {
		if Tags[s] != -1 {
			Tags[s] = -1
		}
	}
	return Tags
}

// Import the tags into the database
func importTags(db sql.DB) {
	tx, err := db.Begin()
	check(err)
	tagsStmt, err := tx.Prepare("insert or ignore into tags(name) values(?)")
	check(err)
	defer tagsStmt.Close()

	for tag, id := range Tags {
		if id <= 0 {
			id = getTagId(db, tag)
		}
		if id <= 0 {
			result, err := tagsStmt.Exec(tag)
			check(err)
			id, err = result.LastInsertId()
		}
		Tags[tag] = id
		log.Println("Found tag id:", id)
	}
	err = tx.Commit()
	check(err)
	// perform insert of the tags if they dont exist and return the ID of that tag into the corresponding hash array
	// if the tag exists fetch its ID into the corresponding hash array
	// return the ID's
}

// Import the words from a given filename into the database
func importWords(db sql.DB, tags string, filename string) {

	tagsToArray(tags)
	importTags(db)
	log.Println(Tags)
	log.Printf("%s", TagsToString())
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

	wtStmt, err := tx.Prepare("insert or ignore into wt(word_id,tag_id) values(?,?)")
	check(err)
	defer wtStmt.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		word := strings.TrimSpace(scanner.Text())
		var word_id int64
		if word != "" {
			log.Debugln("importing word:", word)
			if word_id = getWordId(db, word); word_id <= 0 {
				result, err := wordsStmt.Exec(word)
				check(err)
				word_id, err = result.LastInsertId()
				check(err)
				if word_id == 0 {
					log.Debugf("word %s already exists, fetching", word)
					word_id = getWordId(db, word)
					log.Println("Found word id:", word_id)
				}
			}
			//
			log.Printf("word: %s => id: %d\n", word, word_id)
			for tag, tag_id := range Tags {
				log.Printf("adding wt(%d,%d) // %s %s", word_id, tag_id, word, tag)
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

//
// Search for words matching tags
//
func searchWords(db sql.DB, tags string) {
	tagsToArray(tags)
	importTags(db)
	log.Println(Tags)
	rows, err := db.Query(fmt.Sprintf("select t1.name from words as t1 left join wt as t2 on t2.word_id=t1.id WHERE t2.tag_id IN (%s) group by t1.id", TagsToString()))
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
	debugPtr := flag.Bool("debug", false, "enable debug")

	flag.Parse()
	if *debugPtr {
		log.SetLevel(log.DebugLevel)
	}

	log.Debugln("using db:", *dbPtr)
	log.Debugln("using tags:", *tagsPtr)
	log.Debugln("using words:", *wordsPtr)
	log.Debugln("debug:", *debugPtr)
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
		log.Println("no filename given, performing a search")
		searchWords(*db, *tagsPtr)
	} else {
		log.Println("performing an import on the given files:", flag.Args())
		for i := 0; i < flag.NArg(); i++ {
			log.Println("importing file:", flag.Arg(i))
			importWords(*db, *tagsPtr, flag.Arg(i))
		}
	}

	os.Exit(0)
}
