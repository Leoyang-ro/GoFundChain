# SensorCLI Makefile

# 变量定义
BINARY_NAME=sensorcli
VERSION=1.0.0
BUILD_DIR=build
DIST_DIR=dist

# Go 相关变量
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# 构建标志
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(shell date -u '+%Y-%m-%d_%H:%M:%S')"

# 默认目标
.PHONY: all
all: clean build

# 构建项目
.PHONY: build
build:
	@echo "构建 $(BINARY_NAME)..."
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME) .

# 构建 Windows 版本
.PHONY: build-windows
build-windows:
	@echo "构建 Windows 版本..."
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME).exe .

# 构建 Linux 版本
.PHONY: build-linux
build-linux:
	@echo "构建 Linux 版本..."
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME) .

# 构建 macOS 版本
.PHONY: build-macos
build-macos:
	@echo "构建 macOS 版本..."
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME) .

# 构建所有平台
.PHONY: build-all
build-all: build-windows build-linux build-macos
	@echo "所有平台构建完成"

# 清理构建文件
.PHONY: clean
clean:
	@echo "清理构建文件..."
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_NAME).exe
	rm -rf $(BUILD_DIR)
	rm -rf $(DIST_DIR)

# 运行测试
.PHONY: test
test:
	@echo "运行测试..."
	$(GOTEST) -v ./...

# 运行测试并生成覆盖率报告
.PHONY: test-coverage
test-coverage:
	@echo "运行测试并生成覆盖率报告..."
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "覆盖率报告已生成: coverage.html"

# 安装依赖
.PHONY: deps
deps:
	@echo "安装依赖..."
	$(GOMOD) download
	$(GOMOD) tidy

# 格式化代码
.PHONY: fmt
fmt:
	@echo "格式化代码..."
	$(GOCMD) fmt ./...

# 代码检查
.PHONY: lint
lint:
	@echo "运行代码检查..."
	$(GOCMD) vet ./...

# 安装到系统
.PHONY: install
install: build
	@echo "安装到系统..."
	cp $(BINARY_NAME) /usr/local/bin/ || cp $(BINARY_NAME).exe /usr/local/bin/

# 创建发布包
.PHONY: release
release: clean build-all
	@echo "创建发布包..."
	mkdir -p $(DIST_DIR)
	cp $(BINARY_NAME) $(DIST_DIR)/$(BINARY_NAME)-linux-amd64
	cp $(BINARY_NAME).exe $(DIST_DIR)/$(BINARY_NAME)-windows-amd64.exe
	cp README.md $(DIST_DIR)/
	@echo "发布包已创建在 $(DIST_DIR) 目录"

# 运行示例
.PHONY: example
example: build
	@echo "运行示例..."
	./$(BINARY_NAME) --help
	./$(BINARY_NAME) scan --bus 1
	./$(BINARY_NAME) config show

# 开发模式（监听文件变化并重新构建）
.PHONY: dev
dev:
	@echo "开发模式 - 监听文件变化..."
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "请安装 air: go install github.com/cosmtrek/air@latest"; \
		echo "或者手动运行: make build && ./$(BINARY_NAME)"; \
	fi

# 帮助信息
.PHONY: help
help:
	@echo "SensorCLI Makefile 帮助:"
	@echo ""
	@echo "构建相关:"
	@echo "  make build        - 构建当前平台版本"
	@echo "  make build-windows - 构建 Windows 版本"
	@echo "  make build-linux   - 构建 Linux 版本"
	@echo "  make build-macos   - 构建 macOS 版本"
	@echo "  make build-all     - 构建所有平台版本"
	@echo ""
	@echo "测试相关:"
	@echo "  make test          - 运行测试"
	@echo "  make test-coverage - 运行测试并生成覆盖率报告"
	@echo ""
	@echo "代码质量:"
	@echo "  make fmt           - 格式化代码"
	@echo "  make lint          - 代码检查"
	@echo "  make deps          - 安装依赖"
	@echo ""
	@echo "其他:"
	@echo "  make clean         - 清理构建文件"
	@echo "  make install       - 安装到系统"
	@echo "  make release       - 创建发布包"
	@echo "  make example       - 运行示例"
	@echo "  make dev           - 开发模式"
	@echo "  make help          - 显示此帮助信息"
