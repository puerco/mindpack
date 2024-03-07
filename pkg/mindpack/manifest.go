// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright 2024 Stacklok, Inc

package mindpack

import (
	"encoding/json"
	"fmt"
	"io"
)

// Write writes the bundle manifest to a file
func (m *Manifest) Write(w io.Writer) error {
	e := json.NewEncoder(w)
	e.SetIndent("", "  ")
	if err := e.Encode(&m); err != nil {
		return fmt.Errorf("encoding bundle manifest: %w", err)
	}

	return nil
}

// Read loads the manifest data by parsing json data from reader r
func (m *Manifest) Read(r io.Reader) error {
	dec := json.NewDecoder(r)
	if err := dec.Decode(m); err != nil {
		return fmt.Errorf("decoding manifest: %w", err)
	}
	return nil
}
