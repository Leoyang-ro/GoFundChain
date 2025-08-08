package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

// Level 日志级别
type Level int

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
)

var levelNames = map[Level]string{
	DEBUG: "DEBUG",
	INFO:  "INFO",
	WARN:  "WARN",
	ERROR: "ERROR",
}

// Logger 日志记录器
type Logger struct {
	level  Level
	logger *log.Logger
}

var defaultLogger *Logger

// Init 初始化默认日志记录器
func Init(level Level, logFile string) error {
	var output *os.File
	var err error

	if logFile != "" {
		output, err = os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return fmt.Errorf("打开日志文件失败: %v", err)
		}
	} else {
		output = os.Stdout
	}

	defaultLogger = &Logger{
		level:  level,
		logger: log.New(output, "", log.LstdFlags),
	}

	return nil
}

// SetLevel 设置日志级别
func SetLevel(level Level) {
	if defaultLogger != nil {
		defaultLogger.level = level
	}
}

// log 内部日志记录方法
func (l *Logger) log(level Level, format string, args ...interface{}) {
	if level < l.level {
		return
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	message := fmt.Sprintf(format, args...)
	logMessage := fmt.Sprintf("[%s] %s: %s", timestamp, levelNames[level], message)

	l.logger.Println(logMessage)
}

// Debug 调试日志
func Debug(format string, args ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.log(DEBUG, format, args...)
	}
}

// Info 信息日志
func Info(format string, args ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.log(INFO, format, args...)
	}
}

// Warn 警告日志
func Warn(format string, args ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.log(WARN, format, args...)
	}
}

// Error 错误日志
func Error(format string, args ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.log(ERROR, format, args...)
	}
}

// DeviceLogger 设备操作日志记录器
type DeviceLogger struct {
	deviceAddr uint8
	bus        int
}

// NewDeviceLogger 创建设备日志记录器
func NewDeviceLogger(bus int, addr uint8) *DeviceLogger {
	return &DeviceLogger{
		deviceAddr: addr,
		bus:        bus,
	}
}

// LogRead 记录读取操作
func (dl *DeviceLogger) LogRead(reg uint8, value uint8, err error) {
	if err != nil {
		Error("设备 0x%02X (总线 %d) 读取寄存器 0x%02X 失败: %v", dl.deviceAddr, dl.bus, reg, err)
	} else {
		Debug("设备 0x%02X (总线 %d) 读取寄存器 0x%02X: 0x%02X", dl.deviceAddr, dl.bus, reg, value)
	}
}

// LogWrite 记录写入操作
func (dl *DeviceLogger) LogWrite(reg uint8, value uint8, err error) {
	if err != nil {
		Error("设备 0x%02X (总线 %d) 写入寄存器 0x%02X 失败: %v", dl.deviceAddr, dl.bus, reg, err)
	} else {
		Debug("设备 0x%02X (总线 %d) 写入寄存器 0x%02X: 0x%02X", dl.deviceAddr, dl.bus, reg, value)
	}
}

// LogScan 记录扫描操作
func (dl *DeviceLogger) LogScan(found bool) {
	if found {
		Info("扫描发现设备: 0x%02X (总线 %d)", dl.deviceAddr, dl.bus)
	} else {
		Debug("扫描地址 0x%02X (总线 %d): 无设备", dl.deviceAddr, dl.bus)
	}
}
