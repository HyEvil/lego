package cmd

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"
)


var rootCmd = &cobra.Command{
	Use:   "zeus [sub]",
	Short: "A helper tool for zeus service development",
	Long:  `A helper tool for zeus service development`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
