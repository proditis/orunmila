package main

import (
	"bytes"
	"database/sql"
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	_testDB = "TestSearchSubcmd.db"
	buf     bytes.Buffer
)

func init() {
	dbPtr = &_testDB
	flag.CommandLine.SetOutput(&buf)
}

// Test undefined argument
func TestSearchSubcmdUndefinedArgument(t *testing.T) {
	args := []string{"-undefined"}

	err := searchSubcmd(args)
	assert.EqualError(t, err, `flag provided but not defined: -undefined`)
}

// test help argument
func TestSearchSubcmdHelp(t *testing.T) {
	args := []string{"-help"}
	createDbFileifNotExists(*dbPtr)
	defer os.ReadDir(*dbPtr)
	err := searchSubcmd(args)
	assert.EqualError(t, err, `flag: help requested`)
}

// test tags
func TestSearchSubcmdTags(t *testing.T) {
	args := []string{"-tags", "a,b,c"}
	err := searchSubcmd(args)
	assert.NoError(t, err, `Failed to perform search with tags %v`, args)
}

func TestSearchSubcmdEmptyTags(t *testing.T) {
	args := []string{"-tags", ""}
	err := searchSubcmd(args)
	assert.NoError(t, err, `Failed to perform search with tags %v`, args)
}

func TestSearchSubcmdExistingTags(t *testing.T) {
	var Nrecords int64
	createDbFileifNotExists(*dbPtr)
	defer os.Remove(*dbPtr)

	args := []string{"-tags", "a,b,c", "word1", "word2"}
	addSubcmd(args)
	err := searchSubcmd(args[:2])
	assert.NoError(t, err, `Failed to perform search with tags %v`, args)

	dsn := fmt.Sprintf("file:%s?mode=ro", *dbPtr)
	db, err := sql.Open("sqlite3", dsn)
	assert.NoError(t, err, `Failed to open database %v`, dsn)
	defer db.Close()

	err = db.QueryRow("select count(*) from tags").Scan(&Nrecords)
	assert.NoError(t, err, `Failed to select count(*) from tags`)
	assert.Equal(t, int64(3), Nrecords, `Number of tag records returned did not match`)

	err = db.QueryRow("select count(*) from words").Scan(&Nrecords)
	assert.NoError(t, err, `Failed to select count(*) from words`)
	assert.Equal(t, int64(2), Nrecords, `Number of words records returned did not match`)

	err = db.QueryRow("select count(*) from wt").Scan(&Nrecords)
	assert.NoError(t, err, `Failed to select count(*) from wt`)
	assert.Equal(t, int64(2*3), Nrecords, `Number of words records returned did not match`)
}
