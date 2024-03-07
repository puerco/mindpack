// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright 2024 Stacklok, Inc

package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"sigs.k8s.io/release-utils/log"
	"sigs.k8s.io/release-utils/version"
)

const (
	appname = "mindpack"
)

func init() {
	addPack(rootCmd)
	addInit(rootCmd)
	rootCmd.AddCommand(version.WithFont("doom"))
}

type commandLineOptions struct {
	logLevel string
}

var commandLineOpts = commandLineOptions{
	logLevel: "debug",
}

var rootCmd = &cobra.Command{
	Short:             "manage minder bundles",
	Long:              `manage minder bundles`,
	Use:               "mindpack",
	SilenceUsage:      false,
	PersistentPreRunE: initLogging,
}

func initLogging(*cobra.Command, []string) error {
	return log.SetupGlobalLogger(commandLineOpts.logLevel)
}

// Execute builds the command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logrus.Fatal(err)
	}
}
