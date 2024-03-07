// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright 2024 Stacklok, Inc

package build

import (
	"archive/tar"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/puerco/mindpack/pkg/mindpack"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Packer struct{}

func NewPacker() *Packer {
	return &Packer{}
}

type InitOptions struct {
	*mindpack.Metadata
	Path string
}

// Validate cheks the initializer options
func (io *InitOptions) Validate() error {
	var errs = []error{}
	if io.Name == "" {
		errs = append(errs, fmt.Errorf("name is required to initialize a mindpack"))
	} else if !mindpack.DnsStyleNameRegex.MatchString(io.Name) {
		errs = append(errs, fmt.Errorf("%q is not a valid mindpack name", io.Name))
	}

	if io.Namespace != "" && !mindpack.DnsStyleNameRegex.MatchString(io.Namespace) {
		errs = append(errs, fmt.Errorf("%q is not valida namespace", io.Namespace))
	}

	// FIXME(puerco): Check semver

	// Check path
	sdata, err := os.Stat(io.Path)
	if err != nil {
		errs = append(errs, fmt.Errorf("opening path: %w", err))
	} else {
		if !sdata.IsDir() {
			errs = append(errs, fmt.Errorf("path is not a directory"))
		}
	}

	return errors.Join(errs...)
}

// Initialize creates a new bundle structure in a directory with minder data
func (p *Packer) Init(opts *InitOptions) error {
	if opts.Metadata.Name == "" {
		return fmt.Errorf("unable to initialize new bundle, no name defined")
	}

	bundle, err := mindpack.NewBundleFromDirectory(opts.Path)
	if err != nil {
		return fmt.Errorf("reading source data: %w", err)
	}

	bundle.Metadata = opts.Metadata

	if err := bundle.UpdateManifest(); err != nil {
		return fmt.Errorf("updating new bundle manifest: %w", err)
	}

	bundle.Metadata.Date = timestamppb.Now()

	f, err := os.Create(filepath.Join(opts.Path, mindpack.ManifestFileName))
	if err != nil {
		return fmt.Errorf("opening manifest file: %w", err)
	}

	if err := bundle.Manifest.Write(f); err != nil {
		return fmt.Errorf("writing manifest data: %w", err)
	}

	return nil
}

func (p *Packer) WriteToFile(bundle *mindpack.Bundle, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("opening file: %w", err)
	}
	defer f.Close()

	if err := p.Write(bundle, f); err != nil {
		return fmt.Errorf("writing bundle to file: %w", err)
	}
	return nil
}

// Write writes a bundle archive to writer w
func (p *Packer) Write(bundle *mindpack.Bundle, w io.Writer) error {
	tarWriter := tar.NewWriter(w)
	defer tarWriter.Close()

	if bundle.Source == nil {
		return fmt.Errorf("unable to pack bundle, data source not defined")
	}

	err := fs.WalkDir(bundle.Source, ".", func(path string, _ fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("reading %q: %w", path, err)
		}

		stat, err := fs.Stat(bundle.Source, path)
		if err != nil {
			return fmt.Errorf("reading file info: %w", err)
		}
		if stat.IsDir() {
			return nil
		}

		f, err := bundle.Source.Open(path)
		if err != nil {
			return fmt.Errorf("opening %q", path)
		}
		defer f.Close()

		header := &tar.Header{
			Name:    path,
			Size:    stat.Size(),
			Mode:    int64(stat.Mode()),
			ModTime: stat.ModTime(),
		}

		if err := tarWriter.WriteHeader(header); err != nil {
			return fmt.Errorf("writing header for %q: %w", path, err)
		}

		if _, err := io.Copy(tarWriter, f); err != nil {
			return fmt.Errorf("writing data from %q to archive: %w", path, err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("walking bundle data source: %w", err)
	}

	return nil
}
