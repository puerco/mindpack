// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright 2024 Stacklok, Inc

package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/puerco/mindpack/pkg/build"
	"github.com/puerco/mindpack/pkg/mindpack"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"sigs.k8s.io/release-utils/util"
)

type InitOptions struct {
	SourceDir string
	Name      string
	Namespace string
	Version   string
}

func (io *InitOptions) AddFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVarP(
		&io.SourceDir,
		"source",
		"s",
		"",
		"source directory of the mindpack bundle",
	)

	cmd.PersistentFlags().StringVarP(
		&io.Name,
		"name",
		"n",
		"",
		"name of the bundle",
	)

	cmd.PersistentFlags().StringVar(
		&io.Namespace,
		"ns",
		"",
		"namespace of the bundle",
	)

	cmd.PersistentFlags().StringVarP(
		&io.Version,
		"version",
		"v",
		"v0.0.1",
		"initial version for the new bundle",
	)
}

func (io *InitOptions) Validate() error {
	if io.SourceDir == "" {
		return fmt.Errorf("source directory is required")
	}
	var errs = []error{}
	s, serr := os.Stat(io.SourceDir)
	if serr == nil {
		if !s.IsDir() {
			errs = append(errs, fmt.Errorf("source path is not a directory"))
		}
	} else {
		errs = append(errs, serr)
	}

	if io.Name == "" {
		errs = append(errs, fmt.Errorf("bundle name is required"))
	}

	if util.Exists(filepath.Join(io.SourceDir, mindpack.ManifestFileName)) {
		errs = append(errs, fmt.Errorf("manifest found in %s", filepath.Join(io.SourceDir, mindpack.ManifestFileName)))
	}

	return errors.Join(errs...)
}

func addInit(parentCmd *cobra.Command) {
	opts := InitOptions{}
	mergeCmd := &cobra.Command{
		Short:             "initializes a mindpack source directory",
		Use:               "init [flags]",
		Example:           fmt.Sprintf("%s init --source=bundle-data/ --name=bundle --version=v0.1.0", appname),
		SilenceUsage:      false,
		SilenceErrors:     true,
		PersistentPreRunE: initLogging,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			packer := build.Packer{}

			if err := packer.Init(&build.InitOptions{
				Metadata: &mindpack.Metadata{
					Name:      opts.Name,
					Namespace: opts.Namespace,
					Version:   opts.Version,
				},
				Path: opts.SourceDir,
			}); err != nil {
				return fmt.Errorf("initializing bundle: %w", err)
			}
			logrus.Infof(
				"wrote new manifest in %q",
				filepath.Join(opts.SourceDir, mindpack.ManifestFileName),
			)
			return nil
		},
	}
	opts.AddFlags(mergeCmd)
	parentCmd.AddCommand(mergeCmd)
}
