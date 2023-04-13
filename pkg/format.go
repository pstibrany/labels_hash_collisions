package pkg

import (
	"bufio"
	"encoding/gob"
	"io"
	"os"

	"github.com/grafana/dskit/multierror"
)

type Entry struct {
	Hash   uint64
	Random string
}

type Writer struct {
	f   *os.File
	enc *gob.Encoder
	buf *bufio.Writer
}

func NewEntryWriter(filename string) (*Writer, error) {
	f, err := os.Create(filename)
	if err != nil {
		return nil, err
	}

	buf := bufio.NewWriterSize(f, 1024*1024)

	return &Writer{
		f:   f,
		buf: buf,
		enc: gob.NewEncoder(buf),
	}, nil
}

func (w *Writer) WriteEntry(ent Entry) error {
	return w.enc.Encode(ent)
}

func (w *Writer) Close() error {
	err1 := w.buf.Flush()
	err2 := w.f.Close()

	return multierror.New(err1, err2).Err()
}

func (w *Writer) Name() string {
	return w.f.Name()
}

type EntryReader struct {
	dec *gob.Decoder

	nextValid bool // if true, nextEntry and nextErr have the next entry
	nextEntry Entry
	nextErr   error
}

func NewEntryReader(r io.Reader) *EntryReader {
	buf := bufio.NewReaderSize(r, 100*1024)
	dec := gob.NewDecoder(buf)
	return &EntryReader{dec: dec}
}

// Peek returns next symbol or error, but also preserves them for subsequent Peek or Next calls.
func (sf *EntryReader) Peek() (Entry, error) {
	if sf.nextValid {
		return sf.nextEntry, sf.nextErr
	}

	sf.nextValid = true
	sf.nextEntry, sf.nextErr = sf.readNext()
	return sf.nextEntry, sf.nextErr
}

// Next advances the iterator and returns the next entry or error. io.EOF is returned at the end.
func (sf *EntryReader) Next() (Entry, error) {
	if sf.nextValid {
		defer func() {
			sf.nextValid = false
			sf.nextEntry = Entry{}
			sf.nextErr = nil
		}()
		return sf.nextEntry, sf.nextErr
	}

	return sf.readNext()
}

func (r *EntryReader) readNext() (Entry, error) {
	ent := Entry{}
	err := r.dec.Decode(&ent)
	return ent, err
}
