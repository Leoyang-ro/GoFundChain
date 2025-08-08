package i2c

import (
	"fmt"
	"sync"
)

// MockDevice 模拟I2C设备实现
type MockDevice struct {
	config    *DeviceConfig
	registers map[uint8]uint8
	mu        sync.RWMutex
	closed    bool
}

// NewMockDevice 创建模拟I2C设备
func NewMockDevice(config *DeviceConfig) *MockDevice {
	return &MockDevice{
		config:    config,
		registers: make(map[uint8]uint8),
		closed:    false,
	}
}

// ReadRegister 读取寄存器值
func (dev *MockDevice) ReadRegister(reg uint8) (uint8, error) {
	if dev.closed {
		return 0, fmt.Errorf("设备已关闭")
	}

	dev.mu.RLock()
	defer dev.mu.RUnlock()

	value, exists := dev.registers[reg]
	if !exists {
		// 返回默认值0
		return 0, nil
	}

	return value, nil
}

// WriteRegister 写入寄存器值
func (dev *MockDevice) WriteRegister(reg, value uint8) error {
	if dev.closed {
		return fmt.Errorf("设备已关闭")
	}

	dev.mu.Lock()
	defer dev.mu.Unlock()

	dev.registers[reg] = value
	return nil
}

// ReadBytes 读取多个字节
func (dev *MockDevice) ReadBytes(reg uint8, count int) ([]byte, error) {
	if dev.closed {
		return nil, fmt.Errorf("设备已关闭")
	}

	if count <= 0 {
		return nil, fmt.Errorf("无效的读取字节数: %d", count)
	}

	dev.mu.RLock()
	defer dev.mu.RUnlock()

	data := make([]byte, count)
	for i := 0; i < count; i++ {
		value, exists := dev.registers[reg+uint8(i)]
		if !exists {
			data[i] = 0
		} else {
			data[i] = value
		}
	}

	return data, nil
}

// WriteBytes 写入多个字节
func (dev *MockDevice) WriteBytes(reg uint8, data []byte) error {
	if dev.closed {
		return fmt.Errorf("设备已关闭")
	}

	if len(data) == 0 {
		return fmt.Errorf("写入数据为空")
	}

	dev.mu.Lock()
	defer dev.mu.Unlock()

	for i, value := range data {
		dev.registers[reg+uint8(i)] = value
	}

	return nil
}

// Close 关闭设备
func (dev *MockDevice) Close() error {
	dev.mu.Lock()
	defer dev.mu.Unlock()

	dev.closed = true
	dev.registers = nil
	return nil
}

// GetAddress 获取设备地址
func (dev *MockDevice) GetAddress() uint8 {
	return dev.config.Address
}

// GetBus 获取总线号
func (dev *MockDevice) GetBus() int {
	return dev.config.Bus
}

// GetRegisters 获取所有寄存器值（用于调试）
func (dev *MockDevice) GetRegisters() map[uint8]uint8 {
	dev.mu.RLock()
	defer dev.mu.RUnlock()

	result := make(map[uint8]uint8)
	for k, v := range dev.registers {
		result[k] = v
	}
	return result
}

// IsClosed 检查设备是否已关闭
func (dev *MockDevice) IsClosed() bool {
	dev.mu.RLock()
	defer dev.mu.RUnlock()
	return dev.closed
}

// openPlatform 平台特定的打开函数
func openPlatform(config *DeviceConfig) (Device, error) {
	// 在Windows环境下使用模拟实现
	return NewMockDevice(config), nil
}
