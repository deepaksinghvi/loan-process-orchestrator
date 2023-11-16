/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/deepaksinghvi/loan-process-orchestrator/worker"

	"github.com/spf13/cobra"
)

// workerappCmd represents the workerapp command
var workerappCmd = &cobra.Command{
	Use:   "workerapp",
	Short: "Loan Origination Worker App",
	Long:  `Temporal Workflow Engine Based Loan Origination Worker App.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("Starting Loan Origination Worker")
		worker.StartWorrker()
	},
}

func init() {
	rootCmd.AddCommand(workerappCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// workerappCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// workerappCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
