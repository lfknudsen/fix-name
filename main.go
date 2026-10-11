package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func handleDir(dir string, eradicateAllPluses bool, recursive bool) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}
	rxDoublePlus, _ := regexp.Compile(" {2}")
	rxPlus, _ := regexp.Compile("[+]")
	rxMultiSpace, _ := regexp.Compile("  +")

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
				err = os.Rename(entries[i].Name(), output)
				if err != nil {
					fmt.Println(err)
				} else {
					fmt.Printf("%s  ->  %s\n", info.Name(), output)
				}
			}
		} else if recursive && info.Mode().IsDir() {
			handleDir(filepath.Join(dir, info.Name()), eradicateAllPluses, recursive)
		}
	}
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

	handleDir(".", eradicateAllPluses, recursive)
}
