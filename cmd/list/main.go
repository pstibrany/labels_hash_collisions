package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"labels_hash_collisions/pkg"
)

func main() {
	flag.Parse()

	for _, fn := range flag.Args() {
		printEntries(fn)
	}
}

func printEntries(filename string) {
	f, err := os.Open(filename)
	if err != nil {
		log.Println("failed to open", filename, "due to error:", err)
		return
	}

	defer f.Close()

	r := pkg.NewEntryReader(f)
	var ent pkg.Entry

	for ent, err = r.Next(); err == nil; ent, err = r.Next() {
		fmt.Println(ent.Hash, "\t", ent.Random)
	}

	if errors.Is(err, io.EOF) {
		return
	}
	log.Println("failed to read", filename, "due to error:", err)
}
