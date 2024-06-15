package cmd

import (
	"fmt"
	"os"

	"github.com/denniskniep/json-patch-by-path-demo/src/json"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "jpbp",
	Short: "Patch Json by selecting with JsonPath and edit via JsonPatch",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := runRootCommand(); err != nil {
			return err
		}
		return nil
	},
}

var jsonDoc string
var jsonPath string
var operation string
var value string

func Init() {
	rootCmd.Flags().StringVarP(&jsonDoc, "json", "j", "", "Json Doc (required)")
	rootCmd.MarkFlagRequired("json")

	rootCmd.Flags().StringVarP(&jsonPath, "jsonPath", "p", "", "Json Path (required)")
	rootCmd.MarkFlagRequired("jsonPath")

	rootCmd.Flags().StringVarP(&operation, "operation", "o", "", "Operation for Json Patch (required)")
	rootCmd.MarkFlagRequired("operation")

	rootCmd.Flags().StringVarP(&value, "value", "v", "", "Value (Only necessary for some operations)")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func runRootCommand() error {
	req := json.PatchRequest{
		JsonDoc:   jsonDoc,
		JsonPath:  jsonPath,
		Operation: operation,
		Value:     value,
	}

	result, err := json.Patch(req)
	if err != nil {
		return err
	}

	fmt.Println(result)
	return nil
}
