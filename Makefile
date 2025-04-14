# 定义变量
BINARY_NAME=uid
BUILD_DIR=build
MAIN_PATH=./cmd/uid
VERSION=$(shell git describe --tags --always --dirty)
LDFLAGS=-ldflags "-X main.Version=$(VERSION)"

# 定义所有伪目标
.PHONY: all build clean test bench lint help

# 默认目标
all: lint test build

# 构建目标
build:
	@echo "构建 $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "构建完成: $(BUILD_DIR)/$(BINARY_NAME)"

# 清理构建产物
clean:
	@echo "清理构建产物..."
	@rm -rf $(BUILD_DIR)
	@echo "清理完成"

# 运行测试
test:
	@echo "运行测试..."
	go test -v ./...

# 运行基准测试
bench:
	@echo "运行基准测试..."
	go test -bench=. -benchmem ./crypto

# 代码质量检查
lint:
	@echo "运行代码质量检查..."
	@if command -v golangci-lint >/dev/null; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint 未安装，跳过检查"; \
	fi

# 显示帮助信息
help:
	@echo "可用的命令:"
	@echo "  make all      - 运行 lint、test 和 build"
	@echo "  make build    - 构建程序"
	@echo "  make clean    - 清理构建产物"
	@echo "  make test     - 运行测试"
	@echo "  make bench    - 运行基准测试"
	@echo "  make lint     - 运行代码质量检查"
	@echo "  make help     - 显示帮助信息"
