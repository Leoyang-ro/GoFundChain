package i2c

// Device 定义I2C设备接口
type Device interface {
	// ReadRegister 读取寄存器值
	ReadRegister(reg uint8) (uint8, error)
	
	// WriteRegister 写入寄存器值
	WriteRegister(reg, value uint8) error
	
	// ReadBytes 读取多个字节
	ReadBytes(reg uint8, count int) ([]byte, error)
	
	// WriteBytes 写入多个字节
	WriteBytes(reg uint8, data []byte) error
	
	// Close 关闭设备
	Close() error
}

// Open 打开I2C设备的工厂函数
func Open(bus int, addr uint8) (Device, error) {
	// 根据平台选择实现
	return openPlatform(bus, addr)
}
