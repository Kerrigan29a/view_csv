package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

var version string

func parseArgs(path *string, name *string, column *int, hasColumnNames *bool) {
	flag.Usage = func() {
		fmt.Printf("Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Println("")
		fmt.Printf("Version: %s\n", version)
	}
	flag.StringVar(path, "f", "", "File path")
	flag.StringVar(name, "n", "", "Column name. Implies -1 flag")
	flag.IntVar(column, "c", 0, "Column position. Negative numbers are allowed")
	flag.BoolVar(hasColumnNames, "1", false, "The first line has the names of the columns")
	flag.Parse()

	if *path == "" {
		fmt.Fprintf(os.Stderr, "Must supply a CSV file\n\n")
		flag.Usage()
		os.Exit(1)
	}
	if *name != "" {
		*hasColumnNames = true
	}
}

func main() {
	var path string
	var name string
	var column int
	var hasColumnNames bool
	parseArgs(&path, &name, &column, &hasColumnNames)
	reader := NewReader(path, hasColumnNames)
	if name != "" {
		check(reader.SelectNamedColumn(name))
	} else {
		reader.SelectColumn(column)
	}
	tui := NewTUI(reader)
	check(tui.Run())
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
