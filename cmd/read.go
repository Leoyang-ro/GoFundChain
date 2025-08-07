package cmd

import (
	"fmt"

	"sensorcli/i2c"

	"github.com/spf13/cobra"
)

var (
	readAddr  uint8
	readReg   uint8
	readBus   int
	readCount int
)

var readCmd = &cobra.Command{
	Use:   "read",
	Short: "读取I2C设备寄存器",
	Long: `读取指定I2C设备的寄存器值。

示例:
  sensorcli read --addr 0x48 --reg 0x01 --bus 1
  sensorcli read --addr 0x48 --reg 0x01 --count 4 --bus 1`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return readRegister()
	},
}

func init() {
	rootCmd.AddCommand(readCmd)

	// 添加参数
	readCmd.Flags().Uint8VarP(&readAddr, "addr", "a", 0, "I2C设备地址 (十六进制)")
	readCmd.Flags().Uint8VarP(&readReg, "reg", "r", 0, "寄存器地址 (十六进制)")
	readCmd.Flags().IntVarP(&readBus, "bus", "b", 1, "I2C总线号")
	readCmd.Flags().IntVarP(&readCount, "count", "c", 1, "读取字节数")

	// 设置必需参数
	readCmd.MarkFlagRequired("addr")
	readCmd.MarkFlagRequired("reg")
}

func readRegister() error {
	// 打开I2C设备
	device, err := i2c.Open(readBus, readAddr)
	if err != nil {
		return fmt.Errorf("打开I2C设备失败: %v", err)
	}
	defer device.Close()

	if readCount == 1 {
		// 读取单个寄存器
		value, err := device.ReadRegister(readReg)
		if err != nil {
			return fmt.Errorf("读取寄存器失败: %v", err)
		}

		fmt.Printf("设备 0x%02X 寄存器 0x%02X 的值: 0x%02X (%d)\n",
			readAddr, readReg, value, value)
	} else {
		// 读取多个字节
		data, err := device.ReadBytes(readReg, readCount)
		if err != nil {
			return fmt.Errorf("读取数据失败: %v", err)
		}

		fmt.Printf("设备 0x%02X 寄存器 0x%02X 的 %d 字节数据:\n",
			readAddr, readReg, readCount)

		for i, value := range data {
			fmt.Printf("  0x%02X: 0x%02X (%d)\n", readReg+uint8(i), value, value)
		}
	}

	return nil
}
