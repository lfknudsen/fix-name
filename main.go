package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"
)

var rxDoublePlus = regexp.MustCompile(" {2}")
var rxPlus = regexp.MustCompile("[+]")
var rxMultiSpace = regexp.MustCompile("  +")

func handleDir(dir string, eradicateAllPluses bool, recursive bool, quiet bool, wg *sync.WaitGroup) (int64, int64) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		log.Panic(err)
	}

	count := new(atomic.Int64)
	errors := new(atomic.Int64)
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
					errors.Add(1)
					fmt.Println(err)
				} else {
					count.Add(1)
					if !quiet {
						fmt.Printf("%s  ->  %s\n", info.Name(), output)
					}
				}
			}
		} else if recursive && info.Mode().IsDir() {
			wg.Go(func() {
				subCount, subErrors := handleDir(dir+"/"+info.Name(), eradicateAllPluses, recursive, quiet, wg)
				count.Add(subCount)
				errors.Add(subErrors)
			})
		}
	}
	return count.Load(), errors.Load()
}

func main() {
	var eradicateAllPluses = false
	var recursive = false
	var quiet = false

	os.Args = os.Args[1:]
	for _, arg := range os.Args {
		if arg == "-n" {
			eradicateAllPluses = true
		} else if arg == "-r" {
			recursive = true
		} else if arg == "-q" {
			quiet = true
		} else {
			log.Fatalf("Unknown argument '%s'", arg)
		}
	}

	wg := new(sync.WaitGroup)
	var count int64 = 0
	var errors int64 = 0
	wg.Go(func() {
		subCount, subErrors := handleDir(".", eradicateAllPluses, recursive, quiet, wg)
		count += subCount
		errors += subErrors
	})
	wg.Wait()
	fmt.Printf("Renamed %d files. Encountered %d errors.\n", count, errors)
}
