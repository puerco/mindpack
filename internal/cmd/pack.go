// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright 2024 Stacklok, Inc

package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/puerco/mindpack/pkg/build"
	"github.com/puerco/mindpack/pkg/mindpack"
	"github.com/spf13/cobra"
)

type PackOptions struct {
	SourceDir string
	File      string
}

func (po *PackOptions) AddFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVarP(
		&po.SourceDir,
		"source",
		"s",
		"",
		"source directory of mindpack bundle",
	)

	cmd.PersistentFlags().StringVarP(
		&po.File,
		"file",
		"f",
		"",
		"path to write the bundle",
	)

	cmd.MarkFlagRequired("file")
	cmd.MarkFlagRequired("source")
}

func (po *PackOptions) Validate() error {
	if po.File == "" || po.SourceDir == "" {
		return fmt.Errorf("source directory and file are required")
	}
	var errs = []error{}
	s, serr := os.Stat(po.SourceDir)
	if serr == nil {
		if !s.IsDir() {
			errs = append(errs, fmt.Errorf("source path is not a directory"))
		}
	} else {
		errs = append(errs, serr)
	}
	return errors.Join(errs...)
}

func addPack(parentCmd *cobra.Command) {
	opts := PackOptions{}
	mergeCmd := &cobra.Command{
		Short:             "writes a mindpack bundle to a distributable archive",
		Use:               "pack [flags]",
		Example:           fmt.Sprintf("%s pack --source=file --file=mypack.mpk", appname),
		SilenceUsage:      false,
		SilenceErrors:     true,
		PersistentPreRunE: initLogging,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			packer := build.Packer{}

			bundle, err := mindpack.NewBundleFromDirectory(opts.SourceDir)
			if err != nil {
				return err
			}
			return packer.WriteToFile(bundle, opts.File)
		},
	}
	opts.AddFlags(mergeCmd)
	parentCmd.AddCommand(mergeCmd)
}
