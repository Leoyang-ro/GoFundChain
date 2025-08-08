package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const Version = "1.0.0"

var rootCmd = &cobra.Command{
	Use:   "sensorcli",
	Short: "I2C 调试命令行工具",
	Long: `SensorCLI 是一个跨平台的 I2C 传感器调试工具。

支持功能:
  - I2C 设备扫描
  - 寄存器读写操作
  - 数据导出 (JSON/CSV/HEX)
  - 跨平台支持 (Windows/Linux/macOS)

示例:
  sensorcli scan --bus 1
  sensorcli read --addr 0x48 --reg 0x01 --bus 1
  sensorcli write --addr 0x48 --reg 0x02 --value 0x55 --bus 1
  sensorcli dump --addr 0x48 --reg 0x00 --count 16 --format json`,
	Version: Version,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// 全局初始化逻辑
		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "错误: %v\n", err)
		os.Exit(1)
	}
}
