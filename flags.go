package main

import (
	"flag"
	"fmt"
	"path/filepath"
)

var verbose bool
var filter string
var targets []string

func init() {
	flag.StringVar(&filter, "filter", "", "only process selected files or diresctories")
	flag.StringVar(&filter, "f", "", "")

	flag.BoolVar(&verbose, "verbose", false, "log verbose information")
	flag.BoolVar(&verbose, "v", false, "")

	flag.Parse()

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
