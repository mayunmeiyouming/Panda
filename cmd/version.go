package cmd

import (
  "fmt"

  "github.com/spf13/cobra"
)

func init() {
  rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
  Use:   "version",
  Short: "Print the version number of Panda",
  Long:  `All software has versions. This is Panda's`,
  Run: func(cmd *cobra.Command, args []string) {
    fmt.Println("Panda v0.0.1")
  },
}