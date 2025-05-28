package confetti

import (
	"encoding/json"
	"os"
)

// jsonLoader loads config from a JSON file at the given path.
type jsonLoader string

func (j jsonLoader) Load(config any) (err error) {
	f, err := os.Open(string(j))
	if err != nil {
		return
	}

	defer f.Close() //nolint:errcheck // ok

	dec := json.NewDecoder(f)

	return dec.Decode(config)
}
