# SensorCLI - 跨平台传感器调试套件

## 📋 项目简介

SensorCLI 是一个专为工程师和开发者设计的跨平台传感器调试工具，支持 I2C、SPI、UART 等多种通信协议。

### 🎯 技术架构

```sh
用 C 写核心驱动（如 I2C/SPI 通信）
Go 实现 CLI 调试工具 & 后端服务
Python 做 GUI + 数据可视化（支持 JSON/CSV 导出）
```

## 🚀 快速开始

### 安装

```bash
# 克隆项目
git clone <repository-url>
cd sensorcli-ro

# 编译
go build -o sensorcli

# 运行
./sensorcli --help
```

### 基本使用

#### 扫描I2C设备
```bash
sensorcli scan --bus 1
```

#### 读取设备寄存器
```bash
# 读取单个寄存器
sensorcli read --addr 0x48 --reg 0x01 --bus 1

# 读取多个字节
sensorcli read --addr 0x48 --reg 0x01 --count 4 --bus 1
```

#### 写入设备寄存器
```bash
# 写入单个值
sensorcli write --addr 0x48 --reg 0x02 --value 0x55 --bus 1

# 写入多个字节
sensorcli write --addr 0x48 --reg 0x02 --data 0x55,0x66,0x77 --bus 1
```

#### 导出寄存器数据
```bash
# 导出为JSON格式
sensorcli dump --addr 0x48 --reg 0x00 --count 16 --format json --output data.json

# 导出为CSV格式
sensorcli dump --addr 0x48 --reg 0x00 --count 16 --format csv --output data.csv

# 导出为十六进制格式
sensorcli dump --addr 0x48 --reg 0x00 --count 16 --format hex --output data.txt
```

## 📦 功能特性

| 功能 | 实现方式 | 状态 |
|------|----------|------|
| 命令解析 | `cobra` 框架 | ✅ 已完成 |
| 参数处理 | `pflag` + cobra | ✅ 已完成 |
| I2C 读取 | 跨平台 I2C 接口 | ✅ 已完成 |
| I2C 写入 | 跨平台 I2C 接口 | ✅ 已完成 |
| 设备扫描 | 自动检测 I2C 设备 | ✅ 已完成 |
| 数据导出 | JSON/CSV/HEX 格式 | ✅ 已完成 |
| 模拟模式 | Windows 开发环境支持 | ✅ 已完成 |

## 🏗️ 项目结构

```
sensorcli/
├── cmd/
│   ├── root.go        # 主命令定义
│   ├── read.go        # I2C 读取命令
│   ├── write.go       # I2C 写入命令
│   ├── scan.go        # I2C 设备扫描
│   └── dump.go        # 数据导出命令
├── i2c/
│   ├── interface.go   # I2C 设备接口定义
│   └── mock.go        # 模拟 I2C 实现
├── main.go            # 程序入口
├── go.mod             # Go 模块依赖
└── README.md          # 项目文档
```

## 🔧 开发环境

### 系统要求
- Go 1.24.2 或更高版本
- Windows/Linux/macOS 支持

### 依赖管理
```bash
# 安装依赖
go mod tidy

# 运行测试
go test ./...
```

## 📝 命令参考

### 全局选项
- `--help, -h`: 显示帮助信息
- `--version`: 显示版本信息

### read 命令
读取 I2C 设备寄存器值

**选项:**
- `--addr, -a`: I2C 设备地址 (必需)
- `--reg, -r`: 寄存器地址 (必需)
- `--bus, -b`: I2C 总线号 (默认: 1)
- `--count, -c`: 读取字节数 (默认: 1)

**示例:**
```bash
sensorcli read --addr 0x48 --reg 0x01 --bus 1
sensorcli read --addr 0x48 --reg 0x01 --count 4 --bus 1
```

### write 命令
写入数据到 I2C 设备寄存器

**选项:**
- `--addr, -a`: I2C 设备地址 (必需)
- `--reg, -r`: 寄存器地址 (必需)
- `--value, -v`: 写入值 (十六进制)
- `--data, -d`: 写入的字节数据 (逗号分隔的十六进制值)
- `--bus, -b`: I2C 总线号 (默认: 1)

**示例:**
```bash
sensorcli write --addr 0x48 --reg 0x02 --value 0x55 --bus 1
sensorcli write --addr 0x48 --reg 0x02 --data 0x55,0x66,0x77 --bus 1
```

### scan 命令
扫描 I2C 总线上的设备

**选项:**
- `--bus, -b`: I2C 总线号 (默认: 1)

**示例:**
```bash
sensorcli scan --bus 1
```

### dump 命令
导出 I2C 设备寄存器数据

**选项:**
- `--addr, -a`: I2C 设备地址 (必需)
- `--reg, -r`: 起始寄存器地址 (必需)
- `--count, -c`: 读取字节数 (默认: 16)
- `--format, -f`: 输出格式 (json, csv, hex) (默认: json)
- `--output, -o`: 输出文件路径
- `--bus, -b`: I2C 总线号 (默认: 1)

**示例:**
```bash
sensorcli dump --addr 0x48 --reg 0x00 --count 16 --format json --output data.json
sensorcli dump --addr 0x48 --reg 0x00 --count 16 --format csv --output data.csv
```

## 🔮 未来计划

- [ ] SPI 通信支持
- [ ] UART 通信支持
- [ ] 实时数据监控
- [ ] GUI 界面开发
- [ ] 更多传感器驱动
- [ ] 自动化测试套件

## 🤝 贡献指南

欢迎提交 Issue 和 Pull Request！

## 📄 许可证

本项目采用 MIT 许可证。
