package main

import (
	"flag"
	"fmt"

	log "github.com/sirupsen/logrus"

	_ "github.com/mattn/go-sqlite3"
)

var (
	Tags     = make(map[string]int64)
	Words    = make(map[string]int64)
	dbPtr    *string
	debugPtr *bool
)

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
	exitCode := 0
	var err error

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
		infoSubcmd(args)
	case "import", "imp", "i":
		importSubcmd(args)
	case "search", "sea", "s":
		err = searchSubcmd(args)
	case "vacuum", "vac", "v":
		vacuumSubcmd(args)
	default:
		fmt.Fprintln(flag.CommandLine.Output(), "Unrecognized subcommand:", subcommand)
		flag.Usage()
		exitCode = 1
	}
	if err != nil {
		exitCode = 2
	}
	log.Exit(exitCode)
}
