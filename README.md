# 读取设备地址 0x48 的寄存器 0x01
sensorcli read --addr 0x48 --reg 0x01 --bus 1

# 向设备地址 0x48 的寄存器 0x02 写入 0x55
sensorcli write --addr 0x48 --reg 0x02 --value 0x55 --bus 1

| 功能     | 实现方式                                             |
| ------ | ------------------------------------------------ |
| 命令解析   | `cobra` 框架                                       |
| 参数处理   | `pflag` + cobra                                  |
| I2C 访问 | Linux 的 `/dev/i2c-*`，使用 `syscall` or `periph.io` |
| 可扩展性   | 可以加 SPI、UART、日志导出等模块                             |



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