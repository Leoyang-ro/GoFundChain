# 读取设备地址 0x48 的寄存器 0x01
sensorcli read --addr 0x48 --reg 0x01 --bus 1

# 向设备地址 0x48 的寄存器 0x02 写入 0x55
sensorcli write --addr 0x48 --reg 0x02 --value 0x55 --bus 1

| 功能     | 实现方式                                             |
| ------ | ------------------------------------------------ |
| 命令解析   | `cobra` 框架                                       |
| 参数处理   | `pflag` + cobra                                    |
| I2C 访问 | Linux 的 `/dev/i2c-*`，使用 `syscall` or `periph.io` |
| 可扩展性   | 可以加 SPI、UART、日志导出等模块                     |

```sh
GOPATH D:\goenv
GOROOT C:\Program Files\Go
%GOROOT%\bin %GOPATH%\bin
```

```sh
sensorcli/
├── cmd/
│   ├── root.go        # 主命令定义
│   ├── read.go        # i2c 读取命令
│   └── write.go       # i2c 写入命令
├── i2c/
│   └── linux.go       # 实现 I2C 通信（基于 Linux /dev/i2c-*）
│   └── sensorhub.go   # 实现 I2C 通信（基于 sensorhub /dev/i2c-*）
│   └── RTOS.go        # 实现 I2C 通信（基于 RTOS /dev/i2c-*）
├── main.go            # 程序入口
├── go.mod
└── README.md


cmd/read.go
package cmd

import (
    "fmt"
    "strconv"

    "sensorcli/i2c"
    "github.com/spf13/cobra"
)

var (
    addr uint
    reg  uint
    bus  int
)

var readCmd = &cobra.Command{
    Use:   "read",
    Short: "读取I2C寄存器的值",
    Run: func(cmd *cobra.Command, args []string) {
        data, err := i2c.ReadByte(bus, uint8(addr), uint8(reg))
        if err != nil {
            fmt.Println("读取失败：", err)
            return
        }
        fmt.Printf("读取值: 0x%X\n", data)
    },
}

func init() {
    rootCmd.AddCommand(readCmd)
    readCmd.Flags().UintVar(&addr, "addr", 0x00, "I2C设备地址")
    readCmd.Flags().UintVar(&reg, "reg", 0x00, "寄存器地址")
    readCmd.Flags().IntVar(&bus, "bus", 1, "I2C总线编号")
}

i2c/linux.go

package i2c

import (
    "periph.io/x/conn/v3/i2c/i2creg"
    "periph.io/x/conn/v3/i2c"
    "periph.io/x/host/v3"
)

func ReadByte(busNum int, addr uint8, reg uint8) (byte, error) {
    _, err := host.Init()
    if err != nil {
        return 0, err
    }

    bus, err := i2creg.Open(fmt.Sprintf("/dev/i2c-%d", busNum))
    if err != nil {
        return 0, err
    }
    defer bus.Close()

    dev := &i2c.Dev{Addr: addr, Bus: bus}
    write := []byte{reg}
    read := make([]byte, 1)

    if err := dev.Tx(write, read); err != nil {
        return 0, err
    }
    return read[0], nil
}
```