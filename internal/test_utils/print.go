package testutils

import (
	"encoding/json"
	"fmt"
	"io"
)

func PrettyPrint(in interface{}, out io.Writer) error {
	enc := json.NewEncoder(out)
	enc.SetIndent("", "    ")

	err := enc.Encode(in)
	if err != nil {
		return fmt.Errorf("encoding: %w", err)
	}

	return nil
}
