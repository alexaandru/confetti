package confetti

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
)

// jsonLoader loads config from a JSON file at the given path.
type jsonLoader string

var ErrUnknownFields = errors.New("unknown fields in config")

func (j jsonLoader) Load(config any, ownConfig *confetti) (err error) {
	f, err := os.Open(string(j))
	if err != nil {
		return
	}
	defer f.Close() //nolint:errcheck // ok

	var errOnUnknown bool

	if ownConfig != nil {
		errOnUnknown = ownConfig.errOnUnknown
	}

	return loadJSON(f, config, errOnUnknown)
}

func loadJSON(r io.ReadSeeker, config any, errOnUnknown bool) (err error) {
	// First pass: decode and populate all known fields, ignore unknowns.
	dec := json.NewDecoder(r)
	if err = dec.Decode(config); err != nil {
		return
	}

	if !errOnUnknown {
		return
	}

	// Second pass: rewind and check for unknown fields.
	if _, err = r.Seek(0, io.SeekStart); err != nil {
		return
	}

	dec = json.NewDecoder(r)

	dec.DisallowUnknownFields()

	if err = dec.Decode(config); err != nil {
		return fmt.Errorf("%w: %w", ErrUnknownFields, err)
	}

	return
}
