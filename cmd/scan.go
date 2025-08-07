package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"sensorcli/i2c"
)

var (
	scanBus int
)

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "扫描I2C总线上的设备",
	Long: `扫描指定I2C总线上的所有设备。

示例:
  sensorcli scan --bus 1`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return scanDevices()
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
	
	// 添加参数
	scanCmd.Flags().IntVarP(&scanBus, "bus", "b", 1, "I2C总线号")
}

func scanDevices() error {
	fmt.Printf("扫描I2C总线 %d 上的设备...\n", scanBus)
	
	foundDevices := 0
	
	// 扫描所有可能的I2C地址 (0x03-0x77)
	for addr := uint8(0x03); addr <= 0x77; addr++ {
		// 跳过保留地址
		if addr >= 0x78 && addr <= 0x7F {
			continue
		}
		
		device, err := i2c.Open(scanBus, addr)
		if err != nil {
			continue
		}
		
		// 尝试读取一个寄存器来检测设备是否存在
		_, err = device.ReadRegister(0x00)
		if err == nil {
			fmt.Printf("发现设备: 0x%02X\n", addr)
			foundDevices++
		}
		
		device.Close()
	}
	
	if foundDevices == 0 {
		fmt.Println("未发现任何I2C设备")
	} else {
		fmt.Printf("共发现 %d 个I2C设备\n", foundDevices)
	}
	
	return nil
}
