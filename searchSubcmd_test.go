package main

import (
	"bytes"
	"flag"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test undefined argument
func TestSearchSubcmdUndefinedArgument(t *testing.T) {
	var buf bytes.Buffer
	flag.CommandLine.SetOutput(&buf)
	args := []string{"-undefined"}

	err := searchSubcmd(args)
	assert.EqualError(t, err, `flag provided but not defined: -undefined`)

}

// test help argument
func TestSearchSubcmdHelp(t *testing.T) {
	var buf bytes.Buffer
	flag.CommandLine.SetOutput(&buf)
	args := []string{"-help"}
	err := searchSubcmd(args)
	assert.EqualError(t, err, `flag: help requested`)
}

// test tags
func TestSearchSubcmdTags(t *testing.T) {
	var buf bytes.Buffer
	flag.CommandLine.SetOutput(&buf)
	args := []string{"-tags", "a,b,c"}
	err := searchSubcmd(args)
	t.Error(err)
}
