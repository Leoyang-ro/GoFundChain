package i2c

import (
	"testing"
	"time"
)

func TestMockDevice(t *testing.T) {
	config := &DeviceConfig{
		Bus:      1,
		Address:  0x48,
		Timeout:  1 * time.Second,
		Retries:  3,
		MockMode: true,
	}

	device := NewMockDevice(config)

	// 测试基本属性
	if device.GetAddress() != 0x48 {
		t.Errorf("期望地址 0x48，实际 %02X", device.GetAddress())
	}

	if device.GetBus() != 1 {
		t.Errorf("期望总线 1，实际 %d", device.GetBus())
	}

	// 测试写入和读取
	testReg := uint8(0x10)
	testValue := uint8(0x55)

	err := device.WriteRegister(testReg, testValue)
	if err != nil {
		t.Errorf("写入寄存器失败: %v", err)
	}

	value, err := device.ReadRegister(testReg)
	if err != nil {
		t.Errorf("读取寄存器失败: %v", err)
	}

	if value != testValue {
		t.Errorf("期望值 0x%02X，实际 0x%02X", testValue, value)
	}

	// 测试多字节操作
	testData := []byte{0x11, 0x22, 0x33, 0x44}
	startReg := uint8(0x20)

	err = device.WriteBytes(startReg, testData)
	if err != nil {
		t.Errorf("写入多字节失败: %v", err)
	}

	readData, err := device.ReadBytes(startReg, len(testData))
	if err != nil {
		t.Errorf("读取多字节失败: %v", err)
	}

	if len(readData) != len(testData) {
		t.Errorf("期望长度 %d，实际 %d", len(testData), len(readData))
	}

	for i, expected := range testData {
		if readData[i] != expected {
			t.Errorf("位置 %d: 期望 0x%02X，实际 0x%02X", i, expected, readData[i])
		}
	}

	// 测试关闭设备
	err = device.Close()
	if err != nil {
		t.Errorf("关闭设备失败: %v", err)
	}

	if !device.IsClosed() {
		t.Error("设备应该已关闭")
	}

	// 测试关闭后操作
	_, err = device.ReadRegister(testReg)
	if err == nil {
		t.Error("关闭后读取应该失败")
	}

	err = device.WriteRegister(testReg, testValue)
	if err == nil {
		t.Error("关闭后写入应该失败")
	}
}

func TestMockDeviceConcurrency(t *testing.T) {
	config := &DeviceConfig{
		Bus:      1,
		Address:  0x48,
		Timeout:  1 * time.Second,
		Retries:  3,
		MockMode: true,
	}

	device := NewMockDevice(config)

	// 并发读写测试
	done := make(chan bool, 10)
	for i := 0; i < 5; i++ {
		go func(id int) {
			reg := uint8(0x10 + id)
			value := uint8(0x10 + id)

			err := device.WriteRegister(reg, value)
			if err != nil {
				t.Errorf("并发写入失败: %v", err)
			}

			readValue, err := device.ReadRegister(reg)
			if err != nil {
				t.Errorf("并发读取失败: %v", err)
			}

			if readValue != value {
				t.Errorf("并发测试值不匹配: 期望 0x%02X，实际 0x%02X", value, readValue)
			}

			done <- true
		}(i)
	}

	// 等待所有goroutine完成
	for i := 0; i < 5; i++ {
		<-done
	}

	device.Close()
}

func TestDeviceConfigValidation(t *testing.T) {
	// 测试有效配置
	validConfig := &DeviceConfig{
		Bus:      1,
		Address:  0x48,
		Timeout:  1 * time.Second,
		Retries:  3,
		MockMode: true,
	}

	device, err := OpenWithConfig(validConfig)
	if err != nil {
		t.Errorf("有效配置应该成功: %v", err)
	}
	if device == nil {
		t.Error("应该返回设备实例")
	}

	// 测试无效地址
	invalidAddrConfig := &DeviceConfig{
		Bus:      1,
		Address:  0x02, // 无效地址
		Timeout:  1 * time.Second,
		Retries:  3,
		MockMode: true,
	}

	_, err = OpenWithConfig(invalidAddrConfig)
	if err == nil {
		t.Error("无效地址应该失败")
	}

	// 测试无效总线
	invalidBusConfig := &DeviceConfig{
		Bus:      -1, // 无效总线
		Address:  0x48,
		Timeout:  1 * time.Second,
		Retries:  3,
		MockMode: true,
	}

	_, err = OpenWithConfig(invalidBusConfig)
	if err == nil {
		t.Error("无效总线应该失败")
	}
}
