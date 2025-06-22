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
	serverCmd.Flags().StringP("port", "p", "8080", "Port to run the server on")
	serverCmd.Flags().StringP("host", "H", "localhost", "Host to run the server on")
	serverCmd.Flags().BoolP("auth", "a", false, "Enable authentication")
	serverCmd.Flags().StringP("token", "t", "", "Bearer Token for authentication")
	serverCmd.Flags().BoolP("metrics", "m", false, "Enable Prometheus metrics")
	serverCmd.Flags().BoolP("redoc", "r", false, "Enable ReDoc API documentation")

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

	server.Run(config)

	return nil
}
