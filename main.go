package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"io/fs"
	"os"

	"path/filepath"
)

// version will be set at build time by the build.sh script provided
var version string

// store is used to map hash value to the matching full path name(s)
// duplicates are when the []string slice contains more than 1 element.
// The []string should never be empty nor nil.
var store map[[sha256.Size]byte]([]string)

func main() {

	fmt.Println("doubles - detecting files with identical content")
	fmt.Println("(c) Xavier Gandillot - 2020,2021")
	fmt.Println(version)

	store = make(map[[sha256.Size]byte]([]string), 16)

	if debug {
		fmt.Println("cli arguments --------------")
		fmt.Println("  Verbose    : ", verbose)
		fmt.Println("  Debug      : ", debug)
		fmt.Println("  Continuous : ", continuous)
		fmt.Println("  Filter     : ", filter)
		fmt.Println("  Targets    : ", targets)
		fmt.Println("----------------------------")
	}
	walk()
	if debug {
		dump(false)
	}
	if !continuous {
		summary()
	}
}

// Walk will walk the dirs and sub dirs selected.
// If none specified, will walk currnt dir.
func walk() {

	for _, d := range targets {

		e := filepath.WalkDir(d, wdf)
		if e != nil {
			fmt.Println(e)
		}

	}
}

func wdf(path string, d fs.DirEntry, err error) error {
	if err != nil {
		fmt.Println(err)
		return fs.SkipDir
	}
	if verbose || debug {
		fmt.Printf("Checking : \t%s\n", path)
	}

	if d.IsDir() {
		// its a dir
		if regex != nil && regex.Match([]byte(d.Name())) {
			if verbose || debug {
				fmt.Println("......... \tskiping !")
			}
			return fs.SkipDir
		}
		return nil
	} else {
		// its a file
		if regex != nil && regex.Match([]byte(d.Name())) {
			// match filter, skip !
			if verbose || debug {
				fmt.Println("......... \tskiping !")
			}
			return nil
		} else {
			// no filter match, handle !
			h(path)
		}
	}

	return nil
}

// h handles a valid file name
func h(p string) error {

	var arr [sha256.Size]byte

	f, err := os.Open(p)
	if err != nil {
		return err
	}
	defer f.Close()
	hh := sha256.New()
	if _, err := io.Copy(hh, f); err != nil {
		return err
	}
	// trick to copy slice to the array,
	// so it can be used as map key
	copy(arr[:], hh.Sum(nil))

	ps := store[arr]
	ps = append(ps, p)
	store[arr] = ps
	if len(ps) > 1 && continuous {
		if size {
			fmt.Printf("\nFile size : %d bytes\n", getSize(ps[0]))
		} else {
			fmt.Println()
		}
		for _, pp := range ps {
			fmt.Println(pp)
		}
	}

	return nil
}

// dump the store database of hashes and duplicates
func dump(onlyDuplicates bool) {
	fmt.Println("\n=========== store dump ===========")
	for k, vv := range store {
		if !onlyDuplicates || len(vv) > 1 {
			fmt.Printf("%x\n", k)
			for i, v := range vv {
				fmt.Printf("%d\t%s\n", i+1, v)
			}
		}
	}
}

// summary displays summary of duplicates
func summary() {
	fmt.Println("\nThe following groups of files have identical contents")
	for _, vv := range store {
		if len(vv) > 1 {
			if size {
				fmt.Printf("\nFile size : %d bytes\n", getSize(vv[0]))
			} else {
				fmt.Println()
			}
			for i, v := range vv {
				fmt.Printf("%d\t%s\n", i+1, v)
			}
		}
	}
}

// getSize for a named file.
func getSize(f string) int64 {
	fi, err := os.Stat(f)
	if err != nil {
		fmt.Println(err)
	}
	return fi.Size()
}
