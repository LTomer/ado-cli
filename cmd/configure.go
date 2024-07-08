/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"LTomer/ado-cli/core/config"
	"errors"
	"log"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

func getValueFromUser(label string, errorMessage string) string {
	validate := func(input string) error {
		if len(errorMessage) > 0 && len(input) <= 0 {
			return errors.New(errorMessage)
		}
		return nil
	}

	templates := &promptui.PromptTemplates{
		Prompt:  "{{ . }} ",
		Valid:   "{{ . | green }} ",
		Invalid: "{{ . | red }} ",
		Success: "{{ . | bold }} ",
	}

	prompt := promptui.Prompt{
		Label:     label,
		Templates: templates,
		Validate:  validate,
	}

	result, err := prompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed: %v", err)
	}

	return result
}

// configureCmd represents the configure command
var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		orgLabel := "Enter the Organization name (skip for enter URL):"
		org := getValueFromUser(orgLabel, "")

		url := "https://dev.azure.com/" + org
		if len(org) == 0 {
			urlErrorMessage := "Value could not be empty."
			urlLabel := "Enter the URL of the Azure DevOps platform:"
			url = getValueFromUser(urlLabel, urlErrorMessage)
		}

		patErrorMessage := "Value could not be empty."
		patLabel := "Enter Personal Access Token (PAT):"
		pat := getValueFromUser(patLabel, patErrorMessage)

		configToWrite := &config.Config{
			PAT: pat,
			URL: url,
		}

		if err := config.Write(configToWrite); err != nil {
			log.Fatalf("Error writing config: %v", err)
		}

		// fmt.Printf("Input: %s\n", url, pat)

		// reader := bufio.NewReader(os.Stdin)
		// fmt.Print("Enter Defult Organization on Azure DevOps: ")
		// org, _ := reader.ReadString('\n')
		// fmt.Print("Enter Azure DevOps PAT: ")
		// pat, _ := reader.ReadString('\n')

		// // Example config to write
		// configToWrite := &core.Config{
		// 	PAT:          pat,
		// 	Organization: org,
		// }

		// // Write configuration to file
		// if err := core.Write(configToWrite); err != nil {
		// 	log.Fatalf("Error writing config: %v", err)
		// }
	},
}

func init() {
	rootCmd.AddCommand(configureCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configureCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configureCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
