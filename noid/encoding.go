package noid

// This file handles converting a minter to or from a serializable value

import (
	"encoding/json"
	"io"
	"os"
)

// Instead of making a minter expose everything and serializing tons of
// internals, we really only care about the sequence and template data
type SerializeableMinter struct {
	Template string
	Sequence uint64
}

func (m *Minter) WriteJSON(w io.Writer) error {
	sm := SerializeableMinter{Template: m.Template(), Sequence: m.Sequence()}
	enc := json.NewEncoder(w)
	return enc.Encode(sm)
}

func NewMinterFromJSON(r io.Reader) (*Minter, error) {
	sm := SerializeableMinter{}
	dec := json.NewDecoder(r)
	if err := dec.Decode(&sm); err != nil {
		return nil, err
	}

	return NewSequencedMinter(sm.Template, sm.Sequence)
}

// Reads the given file and converts its JSON data into a minter
func NewMinterFromJSONFile(filename string) (*Minter, error) {
	// First make sure we can open the file
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return NewMinterFromJSON(f)
}
