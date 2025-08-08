package i2c

import (
	"context"
	"fmt"
	"time"
)

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

	// GetAddress 获取设备地址
	GetAddress() uint8

	// GetBus 获取总线号
	GetBus() int
}

// DeviceConfig 设备配置
type DeviceConfig struct {
	Bus      int
	Address  uint8
	Timeout  time.Duration
	Retries  int
	MockMode bool
}

// DefaultConfig 默认设备配置
func DefaultConfig() *DeviceConfig {
	return &DeviceConfig{
		Bus:      1,
		Address:  0,
		Timeout:  1 * time.Second,
		Retries:  3,
		MockMode: true, // Windows 下默认使用模拟模式
	}
}

// Open 打开I2C设备的工厂函数
func Open(bus int, addr uint8) (Device, error) {
	config := DefaultConfig()
	config.Bus = bus
	config.Address = addr
	return OpenWithConfig(config)
}

// OpenWithConfig 使用配置打开I2C设备
func OpenWithConfig(config *DeviceConfig) (Device, error) {
	if config == nil {
		config = DefaultConfig()
	}

	// 验证参数
	if config.Address < 0x03 || config.Address > 0x77 {
		return nil, fmt.Errorf("无效的I2C地址: 0x%02X (有效范围: 0x03-0x77)", config.Address)
	}

	if config.Bus < 0 {
		return nil, fmt.Errorf("无效的总线号: %d", config.Bus)
	}

	// 根据平台选择实现
	return openPlatform(config)
}

// WithTimeout 带超时的操作包装器
func WithTimeout(ctx context.Context, timeout time.Duration, fn func() error) error {
	if timeout <= 0 {
		timeout = 1 * time.Second
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	done := make(chan error, 1)
	go func() {
		done <- fn()
	}()

	select {
	case err := <-done:
		return err
	case <-ctx.Done():
		return fmt.Errorf("操作超时: %v", ctx.Err())
	}
}

// WithRetry 带重试的操作包装器
func WithRetry(maxRetries int, fn func() error) error {
	var lastErr error

	for i := 0; i <= maxRetries; i++ {
		if err := fn(); err == nil {
			return nil
		} else {
			lastErr = err
			if i < maxRetries {
				// 指数退避
				backoff := time.Duration(1<<uint(i)) * 10 * time.Millisecond
				time.Sleep(backoff)
			}
		}
	}

	return fmt.Errorf("操作失败，已重试 %d 次: %v", maxRetries, lastErr)
}
