package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var rxDoublePlus = regexp.MustCompile(" {2}")
var rxPlus = regexp.MustCompile("[+]")
var rxMultiSpace = regexp.MustCompile("  +")

func handleDir(dir string, eradicateAllPluses bool, recursive bool) (count, errors int) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		log.Panic(err)
	}

	count = 0
	errors = 0
	for i := range entries {
		info, err := entries[i].Info()
		if err != nil {
			log.Print(err)
			continue
		}
		if info.Mode().IsRegular() {
			if strings.ContainsRune(info.Name(), '+') && !strings.ContainsRune(info.Name(), ' ') {
				output := info.Name()
				output = rxMultiSpace.ReplaceAllString(output, " ")
				output = rxPlus.ReplaceAllString(output, " ")
				if !eradicateAllPluses {
					output = rxDoublePlus.ReplaceAllString(output, " +")
				}
				output = rxMultiSpace.ReplaceAllString(output, " ")
				output = strings.TrimSpace(output)
				before := filepath.Join(dir, entries[i].Name())
				after := filepath.Join(dir, output)
				err = os.Rename(before, after)
				if err != nil {
					errors++
					fmt.Println(err)
				} else {
					count++
					fmt.Printf("%s  ->  %s\n", info.Name(), output)
				}
			}
		} else if recursive && info.Mode().IsDir() {
			subCount, subErrors := handleDir(dir+"/"+info.Name(), eradicateAllPluses, recursive)
			count += subCount
			errors += subErrors
		}
	}
	return count, errors
}

func main() {
	var eradicateAllPluses = false
	var recursive = false

	os.Args = os.Args[1:]
	for _, arg := range os.Args {
		if arg == "-n" {
			eradicateAllPluses = true
		} else if arg == "-r" {
			recursive = true
		} else {
			log.Fatalf("Unknown argument '%s'", arg)
		}
	}

	count, errors := handleDir(".", eradicateAllPluses, recursive)
	fmt.Printf("Renamed %d files. Encountered %d errors.\n", count, errors)
}
