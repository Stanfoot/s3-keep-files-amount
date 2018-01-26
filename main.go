package main

import (
	"github.com/spf13/cobra"
)


func main() {
	rootCmd := &cobra.Command{Use: "s3-keep-files-amount"}
	rootCmd.AddCommand(Keep)
	rootCmd.Execute()
}
