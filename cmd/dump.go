package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"sensorcli/i2c"
)

var (
	dumpAddr uint8
	dumpReg  uint8
	dumpBus  int
	dumpCount int
	dumpFormat string
	dumpOutput string
)

var dumpCmd = &cobra.Command{
	Use:   "dump",
	Short: "导出寄存器数据",
	Long: `导出I2C设备的寄存器数据到文件。

示例:
  sensorcli dump --addr 0x48 --reg 0x00 --count 16 --format json --output data.json
  sensorcli dump --addr 0x48 --reg 0x00 --count 16 --format csv --output data.csv`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return dumpRegisters()
	},
}

func init() {
	rootCmd.AddCommand(dumpCmd)
	
	// 添加参数
	dumpCmd.Flags().Uint8VarP(&dumpAddr, "addr", "a", 0, "I2C设备地址 (十六进制)")
	dumpCmd.Flags().Uint8VarP(&dumpReg, "reg", "r", 0, "起始寄存器地址 (十六进制)")
	dumpCmd.Flags().IntVarP(&dumpBus, "bus", "b", 1, "I2C总线号")
	dumpCmd.Flags().IntVarP(&dumpCount, "count", "c", 16, "读取字节数")
	dumpCmd.Flags().StringVarP(&dumpFormat, "format", "f", "json", "输出格式 (json, csv, hex)")
	dumpCmd.Flags().StringVarP(&dumpOutput, "output", "o", "", "输出文件路径")
	
	// 设置必需参数
	dumpCmd.MarkFlagRequired("addr")
	dumpCmd.MarkFlagRequired("reg")
}

type RegisterData struct {
	DeviceAddr uint8             `json:"device_addr"`
	StartReg   uint8             `json:"start_register"`
	Timestamp  string            `json:"timestamp"`
	Data       map[string]uint8  `json:"data"`
}

func dumpRegisters() error {
	// 打开I2C设备
	device, err := i2c.Open(dumpBus, dumpAddr)
	if err != nil {
		return fmt.Errorf("打开I2C设备失败: %v", err)
	}
	defer device.Close()

	// 读取数据
	data, err := device.ReadBytes(dumpReg, dumpCount)
	if err != nil {
		return fmt.Errorf("读取数据失败: %v", err)
	}

	// 准备输出数据
	regData := RegisterData{
		DeviceAddr: dumpAddr,
		StartReg:   dumpReg,
		Timestamp:  time.Now().Format(time.RFC3339),
		Data:       make(map[string]uint8),
	}

	for i, value := range data {
		regAddr := dumpReg + uint8(i)
		regData.Data[fmt.Sprintf("0x%02X", regAddr)] = value
	}

	// 根据格式输出
	switch dumpFormat {
	case "json":
		return outputJSON(regData)
	case "csv":
		return outputCSV(regData)
	case "hex":
		return outputHex(regData)
	default:
		return fmt.Errorf("不支持的输出格式: %s", dumpFormat)
	}
}

func outputJSON(data RegisterData) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("JSON编码失败: %v", err)
	}

	if dumpOutput != "" {
		return os.WriteFile(dumpOutput, jsonData, 0644)
	} else {
		fmt.Println(string(jsonData))
		return nil
	}
}

func outputCSV(data RegisterData) error {
	var output string
	output += "Register,Value,Decimal\n"
	
	for reg, value := range data.Data {
		output += fmt.Sprintf("%s,0x%02X,%d\n", reg, value, value)
	}

	if dumpOutput != "" {
		return os.WriteFile(dumpOutput, []byte(output), 0644)
	} else {
		fmt.Print(output)
		return nil
	}
}

func outputHex(data RegisterData) error {
	var output string
	output += fmt.Sprintf("设备地址: 0x%02X\n", data.DeviceAddr)
	output += fmt.Sprintf("起始寄存器: 0x%02X\n", data.StartReg)
	output += fmt.Sprintf("时间戳: %s\n", data.Timestamp)
	output += "数据:\n"
	
	for reg, value := range data.Data {
		output += fmt.Sprintf("  %s: 0x%02X (%d)\n", reg, value, value)
	}

	if dumpOutput != "" {
		return os.WriteFile(dumpOutput, []byte(output), 0644)
	} else {
		fmt.Print(output)
		return nil
	}
}
