package main

import (
	"bytes"
	"flag"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSearchSubcmd(t *testing.T) {
	var buf bytes.Buffer
	flag.CommandLine.SetOutput(&buf)
	args := []string{"-defined"}

	err := searchSubcmd(args)
	assert.EqualError(t, err, `flag provided but not defined: -defined`)

	args = []string{"-help"}
	err = searchSubcmd(args)
	assert.EqualError(t, err, `flag: help requested`)
}
