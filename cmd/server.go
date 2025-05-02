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

	rootCmd.AddCommand(serverCmd)
}

func RunServerCmd(cmd *cobra.Command, args []string) error {
	port, _ := cmd.Flags().GetString("port")
	host, _ := cmd.Flags().GetString("host")
	auth, _ := cmd.Flags().GetBool("auth")
	token, _ := cmd.Flags().GetString("token")

	server.Run(port, host, auth, token)

	return nil
}
