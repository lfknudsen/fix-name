package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"
)

// rxDoublePlus is the regexp which represents the substring to replace with
// doubleSpaceReplacement.
var rxDoublePlus = regexp.MustCompile(" {2}")

// rxChar is the regular expression version of charRune. It is used for replacement.
var rxChar = regexp.MustCompile("[+]")

// rxMultiSpace is just a regexp for more than one space in sequence. These will be trimmed down to
// a single space.
var rxMultiSpace = regexp.MustCompile("  +")

// doubleSpaceReplacement is the string to replace "  " with after all instances of charRune have been
// replaced with " "
var doubleSpaceReplacement = " +"

// charRune is the rune to search file-names for.
var charRune rune = '+'

func handleDir(dir string, wg *sync.WaitGroup) (int64, int64) {
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
			if strings.ContainsRune(info.Name(), charRune) && !strings.ContainsRune(info.Name(), ' ') {
				output := info.Name()
				output = rxMultiSpace.ReplaceAllString(output, " ")
				output = rxChar.ReplaceAllString(output, " ")
				if !*eradicateAll {
					output = rxDoublePlus.ReplaceAllString(output, doubleSpaceReplacement)
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
					if !*quiet {
						fmt.Printf("%s  ->  %s\n", info.Name(), output)
					}
				}
			}
		} else if *recursive && info.Mode().IsDir() {
			wg.Go(func() {
				subCount, subErrors := handleDir(dir+"/"+info.Name(), wg)
				count.Add(subCount)
				errors.Add(subErrors)
			})
		}
	}
	return count.Load(), errors.Load()
}

var help = flag.Bool("h", false, "Print this message and exit.")

var eradicateAll = flag.Bool("n", false, "If set to true, do not attempt to be smart about\n"+
	"intentional characters by converting \"++\" into \" +\".")

var recursive = flag.Bool("r", false, "If true, enter subdirectories recursively as well.")

var quiet = flag.Bool("q", false, "Refrain from printing the list of changed files.")

var char = flag.String("c", "+", "Character to replace with space.")

func printHelp() {
	fmt.Printf("Usage: %s [OPTIONS]\n\nOptions:\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {
	flag.Parse()
	if len(*char) > 1 {
		_, _ = fmt.Fprintf(os.Stderr, "Character to replace must be a single character.")
		printHelp()
		os.Exit(1)
	}
	rxChar = regexp.MustCompile(fmt.Sprintf("[%s]", *char))
	doubleSpaceReplacement = fmt.Sprintf(" %s", *char)
	charRune = []rune(*char)[0]

	wg := new(sync.WaitGroup)
	var count int64 = 0
	var errors int64 = 0
	wg.Go(func() {
		subCount, subErrors := handleDir(".", wg)
		count += subCount
		errors += subErrors
	})
	wg.Wait()
	fmt.Printf("Renamed %d files. Encountered %d errors.\n", count, errors)
}
