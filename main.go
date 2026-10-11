package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

func main() {
	entries, err := os.ReadDir(".")
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
		if !info.Mode().IsRegular() {
			continue
		}
		if strings.ContainsRune(info.Name(), '+') && !strings.ContainsRune(info.Name(), ' ') {
			output := info.Name()
			output = rxMultiSpace.ReplaceAllString(output, " ")
			output = rxPlus.ReplaceAllString(output, " ")
			output = rxDoublePlus.ReplaceAllString(output, " +")
			output = rxMultiSpace.ReplaceAllString(output, " ")
			output = strings.TrimSpace(output)
			err = os.Rename(entries[i].Name(), output)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Printf("%s  ->  %s\n", info.Name(), output)
			}
		}
	}
}
