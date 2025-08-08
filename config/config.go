package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config 全局配置结构
type Config struct {
	DefaultBus     int    `json:"default_bus"`
	DefaultTimeout int    `json:"default_timeout"`
	LogLevel       string `json:"log_level"`
	OutputFormat   string `json:"output_format"`
	MockMode       bool   `json:"mock_mode"`
}

// DefaultConfig 默认配置
func DefaultConfig() *Config {
	return &Config{
		DefaultBus:     1,
		DefaultTimeout: 1000, // 毫秒
		LogLevel:       "info",
		OutputFormat:   "json",
		MockMode:       true, // Windows 下默认使用模拟模式
	}
}

// LoadConfig 加载配置文件
func LoadConfig(configPath string) (*Config, error) {
	config := DefaultConfig()

	if configPath == "" {
		// 尝试从默认位置加载
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return config, nil
		}
		configPath = filepath.Join(homeDir, ".sensorcli", "config.json")
	}

	// 检查配置文件是否存在
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// 创建默认配置文件
		if err := SaveConfig(config, configPath); err != nil {
			return config, fmt.Errorf("创建默认配置文件失败: %v", err)
		}
		return config, nil
	}

	// 读取配置文件
	data, err := os.ReadFile(configPath)
	if err != nil {
		return config, fmt.Errorf("读取配置文件失败: %v", err)
	}

	if err := json.Unmarshal(data, config); err != nil {
		return config, fmt.Errorf("解析配置文件失败: %v", err)
	}

	return config, nil
}

// SaveConfig 保存配置文件
func SaveConfig(config *Config, configPath string) error {
	// 确保目录存在
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建配置目录失败: %v", err)
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化配置失败: %v", err)
	}

	return os.WriteFile(configPath, data, 0644)
}
