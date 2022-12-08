package main

import (
	"database/sql"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func init() {
}
func TestGetDefaultDBPath(t *testing.T) {
	path, _ := os.Getwd()
	msg := getDefaultDBPath()
	want := filepath.Join(path, "orunmila.db")
	if msg != want {
		t.Fatalf(`Error() = %q, want match for %#q, nil`, msg, want)
	}
}

func TestIsFileExists(t *testing.T) {
	if isFileExists("nonexistant.dbfile") {
		t.Fatalf(`Error: isFileExists("nonexistant.dbfile")`)
	}
	_, err := os.Create("empty.db")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove("empty.db")

	if !isFileExists("empty.db") {
		t.Fatal(`Error: isFileExists("empty.db")`)
	}
}

func TestTagsToString(t *testing.T) {
	Tags = make(map[string]int64)

	wants := ""
	if TagsToString() != wants {
		t.Fatal(`Error: TagsToString not empty`)
	}
	Tags["test"] = 1
	wants = "1"
	got := TagsToString()
	if got != wants {
		t.Fatalf(`Error: TagsToString got %s not %s`, got, wants)
	}
	Tags["test2"] = -1
	wants = "1,-1"
	got = TagsToString()
	// hash maps have random order
	if got != wants && got != "-1,1" {
		t.Fatalf(`Error: TagsToString not %v got %v`, wants, got)
	}

}

func TestGetTagId(t *testing.T) {
	// requires db access
	Tags = make(map[string]int64)
	t.Log(`SOFTFAIL: not implemented`)
}
func TestGetWordId(t *testing.T) {
	// requires db access
	t.Log(`SOFTFAIL: not implemented`)
}

func TestCreateDB(t *testing.T) {
	// cleanup any existing files
	os.Remove("random.db")

	// create database
	err := createDB("random.db")
	if err != nil {
		t.Fatalf(`Error: createDB failed %v`, err)
	}

	err = createDB("random.db")
	if err != nil {
		t.Fatalf(`Error: 2nd createDB failed, %v`, err)
	}
	defer os.Remove("random.db")
}

func TestStringToArray(t *testing.T) {
	var wants map[string]int64
	wants = map[string]int64{
		"a": -1,
		"b": -1,
		"c": -1,
	}
	got := stringToArray("a,b,c")
	if !reflect.DeepEqual(got, wants) {
		t.Fatalf(`Error: stringToArray wants %v got %v`, wants, got)
	}

	got = stringToArray("a,b")
	if reflect.DeepEqual(got, wants) {
		t.Fatalf(`Error: stringToArray wants %v got %v`, wants, got)
	}

	got = stringToArray("")
	wants = map[string]int64{}
	if !reflect.DeepEqual(got, wants) {
		t.Fatalf(`Error: stringToArray wants %v got %v`, wants, got)
	}
}

func TestPopulateTagIds(t *testing.T) {
	t.Log(`SOFTFAIL: not implemented`)
}
func TestRemoveEmptyTags(t *testing.T) {
	var wants map[string]int64
	wants = map[string]int64{}
	Tags = map[string]int64{
		"a": -1,
		"b": -1,
		"c": -1,
	}
	removeEmptyTags()
	if !reflect.DeepEqual(Tags, wants) {
		t.Fatalf(`Error: removeEmptyTags want %v got %v`, wants, Tags)
	}

	wants = map[string]int64{"b": 1}
	Tags = map[string]int64{
		"a": -1,
		"b": 1,
		"c": -1,
	}
	removeEmptyTags()
	if !reflect.DeepEqual(Tags, wants) {
		t.Fatalf(`Error: removeEmptyTags want %v got %v`, wants, Tags)
	}
}
func TestImportTags(t *testing.T) {
	var wants map[string]int64
	Tags = make(map[string]int64)
	Tags = map[string]int64{
		"a": -1,
		"b": -1,
		"c": -1,
	}
	wants = map[string]int64{
		"a": 1,
		"b": 2,
		"c": 3,
	}
	createDbFileifNotExists("random.db")
	defer os.Remove("random.db")

	db, err := sql.Open("sqlite3", "random.db")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	importTags(db)
	if !reflect.DeepEqual(Tags, wants) {
		t.Fatalf(`Error: importTags wants %v got %v`, wants, Tags)
	}

	Tags = map[string]int64{
		"a": -1,
		"b": -1,
		"c": -1,
	}
	importTags(db)
	if !reflect.DeepEqual(Tags, wants) {
		t.Fatalf(`Error: importTags didnt populate properly wants %v got %v`, wants, Tags)
	}

	Tags = map[string]int64{
		"a": -1,
		"b": -1,
		"c": -1,
		"d": -1,
	}
	wants["d"] = 4
	importTags(db)
	if !reflect.DeepEqual(Tags, wants) {
		t.Fatalf(`Error: importTags didnt populate properly with new entry wants %v got %v`, wants, Tags)
	}

}
func TestImportFileWords(t *testing.T) {
	t.Log(`SOFTFAIL: not implemented`)
}
func TestSearchWordsByTagIds(t *testing.T) {
	// populate words
	// populate tags
	// populate word_tags

	t.Log(`SOFTFAIL: not implemented`)
}

func TestCreateDbFileifNotExists(t *testing.T) {
	var dbname = "random.db"
	_, err := os.Stat(dbname)
	if err != nil && !os.IsNotExist(err) {
		t.Fatal(err)
	}
	createDbFileifNotExists(dbname)
	defer os.Remove(dbname)

	_, err = os.Stat(dbname)
	if os.IsNotExist(err) {
		t.Fatal(err)
	}
}
