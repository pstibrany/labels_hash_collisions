package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"time"

	"labels_hash_collisions/pkg"
)

func main() {
	flag.Parse()

	filenames := flag.Args()

	fmt.Println("opening", len(filenames), "files")
	it, err := pkg.NewEntriesIterator(filenames)
	if err != nil {
		log.Fatalln("failed to open files", err)
	}
	defer func() { _ = it.Close() }()

	start := time.Now()
	var prevEntry pkg.Entry
	var cnt, duplicates, collisions int

	var ent pkg.Entry
	for ent, err = it.NextEntry(); err == nil; ent, err = it.NextEntry() {
		cnt++

		if cnt%1_000_000 == 0 {
			progress := ent.Hash / uint64(math.MaxUint64/100)
			fmt.Printf("running count: %d, duplicates: %d (%0.2g %%), unique: %d, collisions: %d, last entry hash: %016x, label: %s, progress: %2d%%, elapsed: %v\n",
				cnt, duplicates, 100*float64(duplicates)/float64(cnt), cnt-duplicates, collisions, ent.Hash, ent.Random, progress,
				time.Since(start).Truncate(time.Millisecond))
		}

		if ent.Hash < prevEntry.Hash {
			log.Fatalln("not sorted, prev hash:", prevEntry.Hash, "hash:", ent.Hash)
		}

		if prevEntry.Hash == ent.Hash {
			if prevEntry.Random == ent.Random {
				duplicates++
			} else {
				collisions++
				fmt.Println("collision", prevEntry.Hash, prevEntry.Random, ent.Random)
			}
		}

		prevEntry = ent
	}

	fmt.Printf("total count: %d, duplicates: %d (%0.2g %%), unique: %d, collisions: %d, elapsed: %v\n",
		cnt, duplicates, 100*float64(duplicates)/float64(cnt), cnt-duplicates, collisions,
		time.Since(start).Truncate(time.Millisecond))

	if errors.Is(err, io.EOF) {
		err = nil
	}

	if err != nil {
		log.Fatalln(err)
	}
}
