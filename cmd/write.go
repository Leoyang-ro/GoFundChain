package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"sensorcli/i2c"
)

var (
	writeAddr  uint8
	writeReg   uint8
	writeValue uint8
	writeBus   int
	writeData  []string
)

var writeCmd = &cobra.Command{
	Use:   "write",
	Short: "写入I2C设备寄存器",
	Long: `写入数据到指定I2C设备的寄存器。

示例:
  sensorcli write --addr 0x48 --reg 0x02 --value 0x55 --bus 1
  sensorcli write --addr 0x48 --reg 0x02 --data 0x55,0x66,0x77 --bus 1`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return writeRegister()
	},
}

func init() {
	rootCmd.AddCommand(writeCmd)
	
	// 添加参数
	writeCmd.Flags().Uint8VarP(&writeAddr, "addr", "a", 0, "I2C设备地址 (十六进制)")
	writeCmd.Flags().Uint8VarP(&writeReg, "reg", "r", 0, "寄存器地址 (十六进制)")
	writeCmd.Flags().Uint8VarP(&writeValue, "value", "v", 0, "写入值 (十六进制)")
	writeCmd.Flags().IntVarP(&writeBus, "bus", "b", 1, "I2C总线号")
	writeCmd.Flags().StringSliceVarP(&writeData, "data", "d", nil, "写入的字节数据 (逗号分隔的十六进制值)")
	
	// 设置必需参数
	writeCmd.MarkFlagRequired("addr")
	writeCmd.MarkFlagRequired("reg")
}

func writeRegister() error {
	// 打开I2C设备
	device, err := i2c.Open(writeBus, writeAddr)
	if err != nil {
		return fmt.Errorf("打开I2C设备失败: %v", err)
	}
	defer device.Close()

	if len(writeData) == 0 {
		// 写入单个值
		err := device.WriteRegister(writeReg, writeValue)
		if err != nil {
			return fmt.Errorf("写入寄存器失败: %v", err)
		}
		
		fmt.Printf("已写入设备 0x%02X 寄存器 0x%02X: 0x%02X (%d)\n", 
			writeAddr, writeReg, writeValue, writeValue)
	} else {
		// 写入多个字节
		data := make([]byte, len(writeData))
		for i, hexStr := range writeData {
			var value uint8
			_, err := fmt.Sscanf(hexStr, "%x", &value)
			if err != nil {
				return fmt.Errorf("解析数据失败 %s: %v", hexStr, err)
			}
			data[i] = value
		}
		
		err := device.WriteBytes(writeReg, data)
		if err != nil {
			return fmt.Errorf("写入数据失败: %v", err)
		}
		
		fmt.Printf("已写入设备 0x%02X 寄存器 0x%02X 的 %d 字节数据:\n", 
			writeAddr, writeReg, len(data))
		
		for i, value := range data {
			fmt.Printf("  0x%02X: 0x%02X (%d)\n", writeReg+uint8(i), value, value)
		}
	}

	return nil
}
