// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright 2024 Stacklok, Inc

package mindpack

import (
	"crypto/sha256"
	"fmt"
	"io"
	"io/fs"
	"os"
	"regexp"
	"strings"
)

var (
	DnsStyleNameRegex = regexp.MustCompile(`^[a-zA-Z0-9](?:[-_a-zA-Z0-9]{0,61}[a-zA-Z0-9])?$`)
)

const (
	PathProfiles     = "profiles"
	PathRuleTypes    = "rule_types"
	ManifestFileName = "manifest.json"
	SHA256           = "sha-256"
)

type Bundle struct {
	Manifest *Manifest
	Metadata *Metadata
	Files    *Files
	Source   fs.StatFS
}

func NewBundleFromDirectory(path string) (*Bundle, error) {
	bundle := &Bundle{
		Source: os.DirFS(path).(fs.StatFS),
	}
	if err := bundle.ReadSource(); err != nil {
		return nil, fmt.Errorf("reading bundle data from %q: %w", path, err)
	}

	return bundle, nil
}

// UpdateManifest updates the bundle manifest to reflect the bundle data source
func (b *Bundle) UpdateManifest() error {
	b.Manifest = &Manifest{
		Metadata: b.Metadata,
		Files:    b.Files,
	}
	return nil
}

// ReadSource loads the data from the mindpack source filesystem
func (b *Bundle) ReadSource() error {
	if b.Source == nil {
		return fmt.Errorf("unable to read source, mindpack filesystem not defined")
	}

	b.Manifest = &Manifest{
		Metadata: &Metadata{},
		Files: &Files{
			Profiles:  []*File{},
			RuleTypes: []*File{},
		},
	}

	b.Files = &Files{
		Profiles:  []*File{},
		RuleTypes: []*File{},
	}

	err := fs.WalkDir(b.Source, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("reading %q: %w", path, err)
		}
		if d.IsDir() {
			return nil
		}

		if !strings.HasPrefix(path, PathProfiles+"/") &&
			!strings.HasPrefix(path, PathRuleTypes+"/") &&
			!strings.HasPrefix(path, ManifestFileName) {
			return fmt.Errorf("found unexpected entry in mindpack source: %q", path)
		}

		f, err := b.Source.Open(path)
		if err != nil {
			return fmt.Errorf("opening %q", path)
		}
		defer f.Close()

		if path == ManifestFileName {
			man := &Manifest{}
			if err := man.Read(f); err != nil {
				return fmt.Errorf("parsing manifest: %w", err)
			}
			b.Manifest = man
		}

		h := sha256.New()

		if _, err := io.Copy(h, f); err != nil {
			return fmt.Errorf("hashing %q", path)
		}

		fentry := File{
			Name: d.Name(),
			Hashes: map[string]string{
				"sha-256": fmt.Sprintf("%x", h.Sum(nil)),
			},
		}

		switch {
		case strings.HasPrefix(path, PathProfiles):
			b.Files.Profiles = append(b.Files.Profiles, &fentry)
		case strings.HasPrefix(path, PathRuleTypes):
			b.Files.RuleTypes = append(b.Files.RuleTypes, &fentry)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("traversing bundle data source: %w", err)
	}
	return nil
}

// Verify checks the contents of the bundle against its manifest
func (b *Bundle) Verify() error {
	return nil
}
