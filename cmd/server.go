/*
Copyright (c) 2025 Tobias Schäfer. All rights reserved.
Licensed under the MIT license, see LICENSE in the project root for details.
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tschaefer/rpinfo/server"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the Rest API server",
	Long:  "Start the Rest API server",
	RunE:  RunServerCmd,
}

func init() {
	serverCmd.Flags().String("port", "8080", "Port to run the server on")
	serverCmd.Flags().String("host", "localhost", "Host to run the server on")
	serverCmd.Flags().Bool("auth", false, "Enable authentication")
	serverCmd.Flags().String("token", "", "Bearer Token for authentication")
	serverCmd.Flags().Bool("metrics", false, "Enable Prometheus metrics")
	serverCmd.Flags().Bool("redoc", false, "Enable ReDoc API documentation")
	serverCmd.Flags().String("log-format", "text", "Log format (text, structured, json)")
	serverCmd.Flags().String("log-level", "info", "Log level (debug, info, warn, error)")

	rootCmd.AddCommand(serverCmd)
}

func RunServerCmd(cmd *cobra.Command, args []string) error {
	var config server.Config

	config.Port, _ = cmd.Flags().GetString("port")
	config.Host, _ = cmd.Flags().GetString("host")
	config.Auth, _ = cmd.Flags().GetBool("auth")
	config.Token, _ = cmd.Flags().GetString("token")
	config.Metrics, _ = cmd.Flags().GetBool("metrics")
	config.Redoc, _ = cmd.Flags().GetBool("redoc")
	config.LogFormat, _ = cmd.Flags().GetString("log-format")
	config.LogLevel, _ = cmd.Flags().GetString("log-level")

	server.Run(config)

	return nil
}
