package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
)

func getDefaultDBPath() string {
	path, err := os.Getwd()
	check(err)

	return filepath.Join(path, "orunmila.db")
}

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

// Convert tags' ids from the db into a comma separated string
func TagsToIdsInString() string {
	var a []string
	for _, id := range Tags {
		log.Debugf("[TagsToIdsInString] id => %d\n", id)
		a = append(a, fmt.Sprint(id))
	}

	return strings.Join(a, ",")
}

// Gets the ID of a given tag
func getTagId(db *sql.DB, tag string) int64 {
	var id int64
	err := db.QueryRow("select id from tags where name = ?", tag).Scan(&id)
	if err != nil {
		id = -1
	}
	return id
}

// Gets the ID of a given word
func getWordId(db *sql.DB, word string) int64 {
	var id int64

	err := db.QueryRow("select id from words where name = ?", word).Scan(&id)
	if err != nil {
		id = -1
	}
	return id
}

// Creates the database schema
func createDB(dbname string) error {
	db, err := sql.Open("sqlite3", dbname)
	check(err)
	defer db.Close()

	sqlStmt := `
	create table IF NOT EXISTS words (id integer not null primary key AUTOINCREMENT, name text NOT NULL UNIQUE);
	create table IF NOT EXISTS tags (id integer not null primary key AUTOINCREMENT, name text NOT NULL UNIQUE);
	create table IF NOT EXISTS wt (word_id integer not null , tag_id integer not null, FOREIGN KEY(word_id) REFERENCES words(id),FOREIGN KEY(tag_id) REFERENCES tags(id),PRIMARY KEY(word_id,tag_id));
	create table IF NOT EXISTS sysconfig(name text not null primary key, val text);
	insert or ignore into sysconfig(name,val) values ("version","0.0.0"),("dbname","default");
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return err
	}
	return nil
}

// Split a string into a HASH map of the form Array[word]=-1
func stringToArray(inString string) map[string]int64 {
	// explode tags by comma
	var wordsArray = strings.Split(inString, ",")
	var _tMap = make(map[string]int64)

	// loop through unique items
	for _, s := range wordsArray {
		s = strings.TrimSpace(s)
		if s != "" && _tMap[s] != -1 {
			_tMap[s] = -1
		}
	}

	return _tMap
}

// Populate Tags map array with their corresponding id
// Tags[tag_name]=tag_id
func populateTagIds(db *sql.DB) {
	for tag, id := range Tags {
		if id <= 0 {
			// fetch tag id from database if exists and assign it
			Tags[tag] = getTagId(db, tag)
		}
	}

}

// Removes tags that have not been assigned an ID
// Tags[tag_name]=tag_id
func removeEmptyTags() {
	for tag, id := range Tags {
		if id <= 0 {
			log.Debugln("Removing not found tag:", tag)
			delete(Tags, tag)
		}
	}
}

// Import the tags into the database and populate Tags with tag_id
func importTags(db *sql.DB) {
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
			id, _ = result.LastInsertId()
		}
		Tags[tag] = id
		log.Println("Found tag id:", id)
	}
	err = tx.Commit()
	check(err)
}

// Import the words from a given filename into the database
func importFileWords(db *sql.DB, filename string) {

	importTags(db)
	log.Println(Tags)
	file, err := os.Open(filename)
	check(err)
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
	var lines = 0
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
				if word_id <= 0 {
					log.Debugf("word %s already exists, fetching", word)
					word_id = getWordId(db, word)
					log.Println("Found word id:", word_id)
				}
			}
			//
			log.Debugf("word: %s => id: %d\n", word, word_id)
			for tag, tag_id := range Tags {
				log.Debugf("adding wt(%d,%d) // %s %s", word_id, tag_id, word, tag)
				_, err = wtStmt.Exec(word_id, tag_id)
				if err != nil {
					log.Error(err)
				}
			}
			word_id = 0
		}
		if lines%4000 == 0 {
			log.Info("Lines:", lines)
			err = tx.Commit()
			check(err)
			tx, err = db.Begin()
			check(err)
			wordsStmt, err = tx.Prepare("insert or ignore into words(name) values(?)")
			check(err)
			defer wordsStmt.Close()

			wtStmt, err = tx.Prepare("insert or ignore into wt(word_id,tag_id) values(?,?)")
			check(err)
			defer wtStmt.Close()
		}
		lines++
	}
	err = tx.Commit()
	check(err)
}

// Search for the needle in the haystack
func containsValue(haystack map[string]int64, needle string) bool {
	for tagname := range haystack {
		if tagname == needle {
			return true
		}
	}
	return false
}

//
// Search for words matching tags
//
func searchWordsByTagIds(db *sql.DB, tags string, showTags bool) {
	queryStr := `select t1.name,(select group_concat(name,',') from tags where id in (select tag_id from wt where word_id=t1.id)) as tagged from words as t1`
	populateTagIds(db)
	removeEmptyTags()

	if len(tags) > 0 {
		userTags := strings.Split(tags, ",")
		log.Infoln("Using tags:", userTags)

		for _, CurrentUserTag := range userTags {
			if !containsValue(Tags, CurrentUserTag) {
				log.Warnf("tag %q not found in the db\n", CurrentUserTag)
			}
		}

		queryStr = fmt.Sprintf(queryStr+` left join wt as t2 on t2.word_id=t1.id WHERE t2.tag_id IN (%s) group by t1.id`, TagsToIdsInString())
	} else {
		log.Infoln("No tags were given")
	}
	rows, err := db.Query(queryStr)
	check(err)
	defer rows.Close()
	for rows.Next() {
		var name string
		var tags string
		err = rows.Scan(&name, &tags)
		if err != nil {
			log.Fatal(err)
		}
		if showTags {
			fmt.Println(name, tags)
		} else {
			fmt.Println(name)
		}

	}
	err = rows.Err()
	check(err)
}

// check if file exists
func isFileExists(filename string) bool {
	info, err := os.Stat(filename)
	if err == nil && !info.Mode().IsRegular() {
		log.Fatalf("%q is not a regular file", filename)
	} else if os.IsNotExist(err) {
		return false
	} else if err != nil {
		log.Fatal(err)
	}
	return true
}

// create the db file if it doesn't exists
func createDbFileifNotExists(dbPtr string) {
	if !isFileExists(dbPtr) {
		log.Debugln("database does not exist, creating...")
		createDB(dbPtr)
	}
}
