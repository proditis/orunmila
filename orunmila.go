package main

import (
	"bufio"
	"database/sql"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"

	_ "github.com/mattn/go-sqlite3"
)

var (
	Tags     = make(map[string]int64)
	Words    = make(map[string]int64)
	dbPtr    *string
	debugPtr *bool
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

// Convert tag ids into a comma separated string to be used with our query
func TagsToString() string {
	var a []string
	for _, id := range Tags {
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
func createDB(dbname string) {
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
		return
	}
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
	log.Printf("%s", TagsToString())
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

//
// Search for words matching tags
//
func searchWordsByTagIds(db *sql.DB, tags string) {
	queryStr := "select t1.name from words as t1"
	populateTagIds(db)
	removeEmptyTags()

	if len(Tags) > 0 {
		log.Infoln("Using tags:", Tags)
		queryStr = fmt.Sprintf(queryStr+" left join wt as t2 on t2.word_id=t1.id WHERE t2.tag_id IN (%s) group by t1.id", TagsToString())
	} else {
		log.Infoln("No tags were given or tags no found")
	}
	rows, err := db.Query(queryStr)
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

	db.Close()
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

// parse args of the vaccum subcommand and exec it
func vacuumSubcmd() {
	// parse the global args instead (debug && db path)

	dsn := fmt.Sprintf("file:%s?mode=rw", *dbPtr)
	db, err := sql.Open("sqlite3", dsn)
	check(err)

	defer db.Close()
	_, err = db.Query("VACUUM")
	check(err)

	log.Println("database rebuilt successfully")
	db.Close()
}

// parse args of the describe subcommand and exec it
func describeSubcmd(args []string) {

	flag := flag.NewFlagSet("describe", flag.ExitOnError)

	flag.Usage = func() {
		fmt.Fprint(flag.Output(), "Describe the database\n\n")
		fmt.Fprintf(flag.Output(), "Usage of orunmila descibe:\n")
		flag.PrintDefaults()
		fmt.Fprintln(flag.Output(), "  words strings\n\tspace separated list of words to add")
		fmt.Fprintln(flag.Output(), "\nexample: orunmila describe My Awesome Description")
	}

	flag.Parse(args)

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

	desc := strings.TrimSpace(strings.Join(flag.Args(), " "))
	_, err = descStmt.Exec("description", desc)
	check(err)

	err = tx.Commit()
	check(err)

	db.Close()
}

// parse args of the info subcommand and exec it
func infoSubcmd() {
	dsn := fmt.Sprintf("file:%s?mode=rw", *dbPtr)
	db, err := sql.Open("sqlite3", dsn)
	check(err)
	defer db.Close()

	rows, err := db.Query("SELECT * FROM sysconfig")
	check(err)
	defer rows.Close()

	for rows.Next() {
		var name string
		var val string
		err = rows.Scan(&name, &val)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("[%s]: %s\n", name, val)
	}
	err = rows.Err()
	check(err)

	db.Close()
}

// parse args of the add subcommand and exec it
func addSubcmd(args []string) {
	addCmd := flag.NewFlagSet("add", flag.ExitOnError)

	addCmd.Usage = func() {
		fmt.Fprint(addCmd.Output(), "Add words to the database from the command line with optional tags\n\n")
		fmt.Fprintf(addCmd.Output(), "Usage of orunmila add:\n")
		addCmd.PrintDefaults()
		fmt.Fprintln(addCmd.Output(), "  words strings\n\tspace separated list of words to add")
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

	db.Close()
}

// parse args of the import subcommand and exec it
func importSubcmd(args []string) {
	importCmd := flag.NewFlagSet("import", flag.ExitOnError)

	flag.Usage = func() {
		fmt.Fprint(importCmd.Output(), "Import a word file into the database with optional tags\n\n")
		fmt.Fprintf(importCmd.Output(), "Usage of orunmila import:\n")
		importCmd.PrintDefaults()
		fmt.Fprintln(importCmd.Output(), "  filename\n\tthe filename to read the words from")
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

	db.Close()
}

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

func main() {
	flag.Usage = func() {
		fmt.Fprint(flag.CommandLine.Output(), "~~~~~ orunmila word list manager ~~~~~\n\n")
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of orunmila:\n")
		flag.PrintDefaults()
		fmt.Fprintln(flag.CommandLine.Output(), "\nSubcommands")
		fmt.Fprintln(flag.CommandLine.Output(), "  add      Add words into the database with optional tags")
		fmt.Fprintln(flag.CommandLine.Output(), "  search   Searches the database for given words")
		fmt.Fprintln(flag.CommandLine.Output(), "  import   Import a wordlist file into the database")
		fmt.Fprintln(flag.CommandLine.Output(), "  info     Display database system configuration information")
		fmt.Fprintln(flag.CommandLine.Output(), "  describe Set the database description")
		fmt.Fprintln(flag.CommandLine.Output(), "  vacuum   Rebuild the database file, repacking it into a minimal amount of disk space")
	}

	dbPtr = flag.String("db", getDefaultDBPath(), "the database filename (default: orunmila.db")
	debugPtr = flag.Bool("debug", false, "enable debug")

	flag.Parse()

	if *debugPtr {
		log.SetLevel(log.DebugLevel)
	}

	args := flag.Args()
	var subcommand string
	if len(args) == 0 {
		log.Debugln("[main] No subcommand has been specified, defaulting to search")
		subcommand, args = "search", []string{}
	} else {
		subcommand, args = args[0], args[1:]
	}

	// poor guy, for now
	//log.Debugln("using words:", *wordsPtr)
	log.Debugln("[main] using db:", *dbPtr)
	log.Debugln("[main] debug mode:", *debugPtr)

	createDbFileifNotExists(*dbPtr)

	// poor guy, for now
	// Words = stringToArray(*wordsPtr)

	switch subcommand {
	case "add", "a":
		addSubcmd(args)
	case "describe", "des", "d":
		describeSubcmd(args)
	case "info":
		infoSubcmd()
	case "import", "imp", "i":
		importSubcmd(args)
	case "search", "sea", "s":
		searchSubcmd(args)
	case "vacuum", "vac", "v":
		vacuumSubcmd()
	default:
		log.Errorf("Unrecognized subcommand: %q", subcommand)
		flag.Usage()
		os.Exit(1)
		// TODO print help menu
	}

	os.Exit(0)
}
