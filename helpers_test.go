package main

import (
	"os"
	"path/filepath"
	"testing"
)

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
	t.Fatal(`Error: not implemented`)
}

func TestGetTagId(t *testing.T) {
	t.Fatal(`Error: not implemented`)
}
func TestGetWordId(t *testing.T) {
	t.Fatal(`Error: not implemented`)
}

func TestCreateDB(t *testing.T) {
	t.Fatal(`Error: not implemented`)
}

func TestStringToArray(t *testing.T) {
	t.Fatal(`Error: not implemented`)
}

func TestPopulateTagIds(t *testing.T) {
	t.Fatal(`Error: not implemented`)
}
func Test(t *testing.T) {
	t.Fatal(`Error: not implemented`)
}
func TestRemoveEmptyTags(t *testing.T) {
	t.Fatal(`Error: not implemented`)
}
func TestImportTags(t *testing.T) {
	t.Fatal(`Error: not implemented`)
}
func TestImportFileWords(t *testing.T) {
	t.Fatal(`Error: not implemented`)
}
func TestSearchWordsByTagIds(t *testing.T) {
	t.Fatal(`Error: not implemented`)
}
func TestCreateDbFileifNotExists(t *testing.T) {
	t.Fatal(`Error: not implemented`)
}
