package confetti

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

// jsonLoader loads config from a JSON file, []byte, or io.Reader.
type jsonLoader struct {
	r   io.ReadSeeker
	c   io.Closer
	err error
}

var (
	ErrUnknownFields = errors.New("unknown fields in config")
	ErrNoDataSource  = errors.New("no data source for JSON loader")
)

func (j jsonLoader) Load(config any, ownConfig *confetti) (err error) {
	if j.err != nil {
		return j.err
	}

	if j.c != nil {
		defer j.c.Close() //nolint:errcheck // ok
	}

	if j.r == nil {
		return ErrNoDataSource
	}

	var errOnUnknown bool

	if ownConfig != nil {
		errOnUnknown = ownConfig.errOnUnknown
	}

	return loadJSON(j.r, config, errOnUnknown)
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
