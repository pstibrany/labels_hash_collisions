package main

import (
	"context"
	"crypto/rand"
	"flag"
	"fmt"
	"log"
	"sort"

	"github.com/grafana/dskit/concurrency"
	"github.com/prometheus/prometheus/model/labels"

	"labels_hash_collisions/pkg"
)

// can be modified to use fewer characters from alphabet slice.
const bitsPerChar = 6
const alphabetLength = 1 << bitsPerChar
const mask = 0xff >> (8 - bitsPerChar)

var alphabet []byte

func init() {
	for b := 'a'; b <= 'z'; b++ {
		alphabet = append(alphabet, byte(b))
	}
	for b := 'A'; b <= 'Z'; b++ {
		alphabet = append(alphabet, byte(b))
	}
	// Put numbers last, then we can use bitsPerChar = 5 to only use letters (eg. when using random
	// strings for metric names, to avoid starting with a number).
	for b := '0'; b <= '9'; b++ {
		alphabet = append(alphabet, byte(b))
	}
	alphabet = append(alphabet, '_')
	alphabet = append(alphabet, '.')

	if len(alphabet) < alphabetLength {
		panic("invalid alphabet length")
	}
}

func main() {
	var (
		batchSize  int
		prefix     string
		filesCount int
		concurrent int
	)
	flag.IntVar(&batchSize, "s", 10_000_000, "Number of entries in file.")
	flag.StringVar(&prefix, "p", "hashes.", "Prefix for files")
	flag.IntVar(&filesCount, "b", 1000, "Number of files")
	flag.IntVar(&concurrent, "c", 4, "Number of concurrent file generations")
	flag.Parse()

	_ = concurrency.ForEachJob(context.Background(), filesCount, concurrent, func(ctx context.Context, idx int) error {
		fname := fmt.Sprintf("%s%05d", prefix, idx)

		log.Println("Generating entries for batch", idx, "file", fname)
		ents := generateRandomEntries(batchSize)

		w, err := pkg.NewEntryWriter(fname)
		if err != nil {
			log.Fatalln("failed to create file:", err)
		}

		for _, e := range ents {
			err = w.WriteEntry(e)
			if err != nil {
				log.Fatalln("failed to write entry:", err)
			}
		}

		err = w.Close()
		if err != nil {
			log.Fatalln("failed to close file:", err)
		}

		log.Println("Generated file", fname)
		return nil
	})
}

func generateRandomEntries(count int) []pkg.Entry {
	var result []pkg.Entry

	var buf, resultBuf []byte
	b := labels.NewBuilder(labels.Labels{})

	for i := 0; i < count; i++ {
		// We need enough bytes such that 256^N is higher than total number of generated entries.
		// With 5 bytes, we get 256^5 = 1 099 511 627 776 possible strings, which is enough to
		// generate our required number of unique entries to guarantee collisions with high probability
		// (this program defaults to 10B).
		var s string
		s, buf, resultBuf = generateRandomString(5, buf, resultBuf)

		// Here we can define how our series look like. We only store hash and random string in the files,
		// not full series (it would make files bigger and slow down finding duplicates)
		b.Reset(labels.Labels{})
		b.Set(labels.MetricName, "metric")
		b.Set("lbl", s)

		/*
			alternative:
				b.Set(labels.MetricName, "metric")
				b.Set("lbl1", s[:4])
				b.Set("lbl2", s[4:])

			alternative:
				b.Set(labels.MetricName, s[0:4]) // When generating metric name, use bitsPerChar=5, to avoid metric names with leading digits or dots.
				b.Set("lbl", s[4:])
		*/

		lbls := b.Labels()

		e := pkg.Entry{
			Hash:   lbls.Hash(),
			Random: s,
		}
		result = append(result, e)
	}

	// Sort by hashes, to make it easy to find collisions.
	sort.Slice(result, func(i, j int) bool {
		return result[i].Hash < result[j].Hash
	})

	return result
}

// Generate random string using specified number of bytes and generated alphabet.
// We only use bitsPerChar to find letter from the alphabet.
func generateRandomString(randBytes int, buf, resultBuf []byte) (string, []byte, []byte) {
	if cap(buf) < randBytes {
		buf = make([]byte, randBytes)
	}
	buf = buf[:randBytes]
	resultBuf = resultBuf[:0]

	_, err := rand.Read(buf)
	if err != nil {
		panic(err)
	}

	word := uint16(0) // higher byte is used as buffer for extra bits
	extraBits := 0    // number of valid bits stored in higher byte of word.

	for _, b := range buf {
		word = word << 8 // move the remaining bits to upper byte before mixing new byte
		word = word | uint16(b)
		extraBits += 8

		for extraBits >= bitsPerChar {
			ch := byte(word & mask)
			resultBuf = append(resultBuf, alphabet[ch])

			word = word >> bitsPerChar // this also moves buffered bits in upper byte
			extraBits -= bitsPerChar
		}
	}

	// Use remaining bits too, don't waste good random bits.
	for extraBits > 0 {
		ch := byte(word & mask)
		resultBuf = append(resultBuf, alphabet[ch])

		word = word >> bitsPerChar
		extraBits -= bitsPerChar
	}

	return string(resultBuf), buf, resultBuf
}
