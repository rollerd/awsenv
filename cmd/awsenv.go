/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// awsenvCmd represents the awsenv command
var awsenvCmd = &cobra.Command{
	Use:   "awsenv",
	Short: "Update AWS credentials to use specified profile",
	Long:  `Update the 'default' AWS profile with the credentials from the specified profile (dev,prod,etc.)`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		newEnv := args[0]
		updateEnv(newEnv)
	},
}

func init() {
	rootCmd.AddCommand(awsenvCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// awsenvCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// awsenvCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
