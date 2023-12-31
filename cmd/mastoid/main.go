package main

import (
	"fmt"
	clay "github.com/go-go-golems/clay/pkg"
	"github.com/go-go-golems/glazed/pkg/cli"
	"github.com/go-go-golems/glazed/pkg/help"
	cmds "github.com/go-go-golems/mastoid/cmd/mastoid/cmds"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "mastoid",
	Short: "mastoid is a CLI app to interact with Mastodon",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// reinitialize the logger because we can now parse --log-level and co
		// from the command line flag
		err := clay.InitLogger()
		cobra.CheckErr(err)
	},
}

func initRootCmd() (*help.HelpSystem, error) {
	helpSystem := help.NewHelpSystem()
	helpSystem.SetupCobraRootCommand(rootCmd)

	err := clay.InitViper("mastoid", rootCmd)
	if err != nil {
		return nil, err
	}
	err = clay.InitLogger()
	if err != nil {
		return nil, err
	}

	return helpSystem, nil
}

func main() {
	_, err := initRootCmd()
	cobra.CheckErr(err)

	cmds.RenderCmd.Flags().StringP("status-id", "s", "", "Status ID")
	cmds.RenderCmd.Flags().BoolP("verbose", "v", false, "Verbose output")
	cmds.RenderCmd.Flags().String("output", "markdown", "Output format (html, text, markdown, json)")
	cmds.RenderCmd.Flags().Bool("with-header", true, "Print header")
	rootCmd.AddCommand(cmds.RenderCmd)

	cmds.RegisterCmd.Flags().StringP("client-name", "n", "mastoid", "Client name")
	cmds.RegisterCmd.Flags().StringP("redirect-uris", "r", "urn:ietf:wg:oauth:2.0:oob", "Redirect URIs")
	cmds.RegisterCmd.Flags().StringP("scopes", "s", "read", "Scopes")
	cmds.RegisterCmd.Flags().StringP("website", "w", "", "Website")
	cmds.RegisterCmd.Flags().StringP("server", "v", "https://hachyderm.io", "Mastodon instance")

	rootCmd.AddCommand(cmds.RegisterCmd)

	rootCmd.AddCommand(cmds.AuthorizeCmd)

	rootCmd.AddCommand(cmds.VerifyCmd)

	threadCmd, err := cmds.NewThreadCommand()
	cobra.CheckErr(err)
	command, err := cli.BuildCobraCommandFromGlazeCommand(threadCmd)
	cobra.CheckErr(err)
	rootCmd.AddCommand(command)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
