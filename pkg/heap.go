package pkg

import (
	"container/heap"
	"io"
	"os"

	"github.com/prometheus/prometheus/tsdb/errors"
)

// Implements heap.Interface
type readersHeap []*EntryReader

// Len implements sort.Interface.
func (s *readersHeap) Len() int {
	return len(*s)
}

// Less implements sort.Interface.
func (s *readersHeap) Less(i, j int) bool {
	iw, ierr := (*s)[i].Peek()
	if ierr != nil {
		// Entry with 0 hash will be first, so error will be returned quickly.
		iw = Entry{}
	}

	jw, jerr := (*s)[j].Peek()
	if jerr != nil {
		jw = Entry{}
	}

	return iw.Hash < jw.Hash
}

// Swap implements sort.Interface.
func (s *readersHeap) Swap(i, j int) {
	(*s)[i], (*s)[j] = (*s)[j], (*s)[i]
}

// Push implements heap.Interface. Push should add x as element Len().
func (s *readersHeap) Push(x interface{}) {
	*s = append(*s, x.(*EntryReader))
}

// Pop implements heap.Interface. Pop should remove and return element Len() - 1.
func (s *readersHeap) Pop() interface{} {
	l := len(*s)
	res := (*s)[l-1]
	*s = (*s)[:l-1]
	return res
}

type EntriesIterator struct {
	files []*os.File
	heap  readersHeap
}

func NewEntriesIterator(filenames []string) (*EntriesIterator, error) {
	files, err := openFiles(filenames)
	if err != nil {
		return nil, err
	}

	var symFiles []*EntryReader
	for _, f := range files {
		symFiles = append(symFiles, NewEntryReader(f))
	}

	h := &EntriesIterator{
		files: files,
		heap:  symFiles,
	}

	heap.Init(&h.heap)

	return h, nil
}

// NextEntry advances iterator forward, and returns next entry.
// If there is no next entry, returns err == io.EOF.
func (sit *EntriesIterator) NextEntry() (Entry, error) {
	for len(sit.heap) > 0 {
		result, err := sit.heap[0].Next()
		if err == io.EOF {
			// End of file, remove it from heap, and try next file.
			heap.Remove(&sit.heap, 0)
			continue
		}

		if err != nil {
			return Entry{}, err
		}

		heap.Fix(&sit.heap, 0)
		return result, nil
	}

	return Entry{}, io.EOF
}

// Close all files.
func (sit *EntriesIterator) Close() error {
	errs := errors.NewMulti()
	for _, f := range sit.files {
		errs.Add(f.Close())
	}
	return errs.Err()
}

func openFiles(filenames []string) ([]*os.File, error) {
	var result []*os.File

	for _, fn := range filenames {
		f, err := os.Open(fn)
		if err != nil {
			// Close files opened so far.
			for _, sf := range result {
				_ = sf.Close()
			}
			return nil, err
		}

		result = append(result, f)
	}
	return result, nil
}
