package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"sensorcli/config"

	"github.com/spf13/cobra"
)

var (
	configPath string
	// 配置参数
	setDefaultBus     int
	setDefaultTimeout int
	setLogLevel       string
	setOutputFormat   string
	setMockMode       bool
)

func init() {
	// 创建主配置命令
	configCmd := &cobra.Command{
		Use:   "config",
		Short: "管理配置文件",
		Long: `管理 SensorCLI 的配置文件。

示例:
  sensorcli config show
  sensorcli config set --default-bus 2
  sensorcli config reset`,
	}

	// 创建子命令
	showConfigCmd := &cobra.Command{
		Use:   "show",
		Short: "显示当前配置",
		RunE: func(cmd *cobra.Command, args []string) error {
			return showConfig()
		},
	}

	setConfigCmd := &cobra.Command{
		Use:   "set",
		Short: "设置配置项",
		RunE: func(cmd *cobra.Command, args []string) error {
			return setConfig()
		},
	}

	resetConfigCmd := &cobra.Command{
		Use:   "reset",
		Short: "重置为默认配置",
		RunE: func(cmd *cobra.Command, args []string) error {
			return resetConfig()
		},
	}

	// 添加子命令
	configCmd.AddCommand(showConfigCmd)
	configCmd.AddCommand(setConfigCmd)
	configCmd.AddCommand(resetConfigCmd)

	// 全局配置路径参数
	configCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "", "配置文件路径")

	// set 命令参数
	setConfigCmd.Flags().IntVarP(&setDefaultBus, "default-bus", "b", 0, "默认I2C总线号")
	setConfigCmd.Flags().IntVarP(&setDefaultTimeout, "default-timeout", "t", 0, "默认超时时间(毫秒)")
	setConfigCmd.Flags().StringVarP(&setLogLevel, "log-level", "l", "", "日志级别 (debug, info, warn, error)")
	setConfigCmd.Flags().StringVarP(&setOutputFormat, "output-format", "f", "", "默认输出格式 (json, csv, hex)")
	setConfigCmd.Flags().BoolVarP(&setMockMode, "mock-mode", "m", false, "是否启用模拟模式")

	// 添加到根命令
	rootCmd.AddCommand(configCmd)
}

func showConfig() error {
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("加载配置失败: %v", err)
	}

	fmt.Println("当前配置:")
	fmt.Printf("  默认总线: %d\n", cfg.DefaultBus)
	fmt.Printf("  默认超时: %d ms\n", cfg.DefaultTimeout)
	fmt.Printf("  日志级别: %s\n", cfg.LogLevel)
	fmt.Printf("  输出格式: %s\n", cfg.OutputFormat)
	fmt.Printf("  模拟模式: %t\n", cfg.MockMode)

	if configPath == "" {
		homeDir, _ := os.UserHomeDir()
		defaultPath := filepath.Join(homeDir, ".sensorcli", "config.json")
		fmt.Printf("\n配置文件路径: %s\n", defaultPath)
	} else {
		fmt.Printf("\n配置文件路径: %s\n", configPath)
	}

	return nil
}

func setConfig() error {
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("加载配置失败: %v", err)
	}

	// 更新配置
	if setDefaultBus > 0 {
		cfg.DefaultBus = setDefaultBus
	}
	if setDefaultTimeout > 0 {
		cfg.DefaultTimeout = setDefaultTimeout
	}
	if setLogLevel != "" {
		cfg.LogLevel = setLogLevel
	}
	if setOutputFormat != "" {
		cfg.OutputFormat = setOutputFormat
	}
	// MockMode 是布尔值，需要特殊处理
	// 这里我们直接使用 setMockMode 的值，因为它是通过命令行参数设置的
	if setMockMode {
		cfg.MockMode = setMockMode
	}

	// 保存配置
	if err := config.SaveConfig(cfg, configPath); err != nil {
		return fmt.Errorf("保存配置失败: %v", err)
	}

	fmt.Println("配置已更新")
	return showConfig()
}

func resetConfig() error {
	cfg := config.DefaultConfig()

	if err := config.SaveConfig(cfg, configPath); err != nil {
		return fmt.Errorf("重置配置失败: %v", err)
	}

	fmt.Println("配置已重置为默认值")
	return showConfig()
}
