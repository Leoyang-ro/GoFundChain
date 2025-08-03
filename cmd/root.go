package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "sensorcli",
	Short: "I2C 调试命令行工具",
	Long:  `sensorcli is a command line tool for I2C debugging`,
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
