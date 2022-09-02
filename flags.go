package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

var verbose bool
var filter string
var regex *regexp.Regexp // nil if no filter defined
var targets []string
var debug, size bool
var continuous bool

func init() {

	flag.StringVar(&filter, "filter", "", "regex pattern of basename of file or dir names to ignore. The pattern must match the full name, not just part of it. Matching dirs will be skiped entirely and recursively")
	flag.StringVar(&filter, "f", "", "")

	flag.BoolVar(&verbose, "verbose", false, "print verbose information")
	flag.BoolVar(&verbose, "v", false, "")

	flag.BoolVar(&continuous, "continuous", false, "print duplicates as they are found")
	flag.BoolVar(&continuous, "c", false, "")

	flag.BoolVar(&debug, "debug", false, "print debug information")
	flag.BoolVar(&debug, "d", false, "")

	flag.BoolVar(&size, "size", false, "print size information")
	flag.BoolVar(&size, "s", false, "")

	flag.Usage = func() {
		fmt.Println("Usage : doubles [OPTION FLAGS] PATH1 PATH2 PATH3 ... ")
		fmt.Println("By default, if PATH is not set, the working directory (.) is used")
		fmt.Println(version)
		flag.PrintDefaults()
		os.Exit(0)
	}

	flag.Parse()

	if filter != "" {
		regex = regexp.MustCompile(filter)
	}

	targets = flag.Args()
	// if none specified, use current dir
	if len(targets) == 0 {
		targets = append(targets, ".")
	}

	// now, we convert all targets to absolute paths
	for i, p := range targets {
		ap, err := filepath.Abs(p)
		if err != nil {
			fmt.Println(err)
		} else {
			targets[i] = ap
		}
	}
}
