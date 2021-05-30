package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"io/fs"
	"os"

	"path/filepath"
)

// store is used to map hash value to the matching full path name(s)
// duplicates are when the []string slice contains more than 1 element.
// The []string should bnever be empty nor nil.
var store map[[sha256.Size]byte]([]string)

func main() {
	store = make(map[[sha256.Size]byte]([]string), 16)

	fmt.Println("cli arguments :", verbose, filter, targets)
	walk()
	dump(false)
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
	if verbose {
		fmt.Printf("Checking : \t%s\n", path)
	}

	if d.IsDir() {
		return nil
	} else {
		h(path)
	}

	return nil
}

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
	// trick to copy slice to the array ...
	copy(arr[:], hh.Sum(nil))
	ps, ok := store[arr]
	ps = append(ps, p)
	store[arr] = ps
	if ok {
		fmt.Println("Duplicate exists !")
	}

	return nil
}

// dump the store database of hashes and duplicates
func dump(onlyDuplicates bool) {
	for k, vv := range store {
		if !onlyDuplicates || len(vv) > 1 {
			fmt.Printf("%x\n", k)
			for i, v := range vv {
				fmt.Printf("%d\t%s\n", i, v)
			}
		}
	}
}
