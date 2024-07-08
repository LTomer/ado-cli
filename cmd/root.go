/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"LTomer/ado-cli/core/ado"
	"LTomer/ado-cli/core/config"
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ado-cli",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },

	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Setup global configurations based on flags here
		// For example, setting up logging based on log-level
		profile, _ := cmd.Flags().GetString("profile")
		// Implement your logic to set the profile based on the flag value
		if profile != "" {
			fmt.Printf("> Select %s profile\n", profile)
		}
		c, _ := config.Read(&profile)
		// ctx := context.WithValue(cmd.Context(), core.Config{}, *c)

		adoClient := ado.NewClient(c.URL, c.PAT)
		ctx := context.WithValue(cmd.Context(), ado.Ado{}, *adoClient)

		cmd.SetContext(ctx)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ado-cli.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.

	// // rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	rootCmd.PersistentFlags().String("profile", "", "Use a specific profile from your config file.")
}
