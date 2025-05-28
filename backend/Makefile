.PHONY: run build swagger clean test dev

# 默认目标
all: build

# 运行项目
run:
	go run main.go

# 开发模式（热更新）
dev:
	~/go/bin/gin --port 9000 run main.go

# 构建项目
build:
	go build -o bin/app main.go

# 生成 swagger 文档
swagger:
	~/go/bin/swag init -g main.go -o docs

# 清理构建文件
clean:
	rm -rf bin/
	rm -rf docs/
	rm -rf tmp/

# 运行测试
test:
	go test -v ./...

# 安装依赖
deps:
	go mod download

# 更新依赖
update-deps:
	go get -u ./...

# 代码格式化
fmt:
	go fmt ./...

# 代码检查
lint:
	golangci-lint run

# 帮助信息
help:
	@echo "可用的命令:"
	@echo "  make run        - 运行项目"
	@echo "  make dev        - 开发模式（热更新）"
	@echo "  make build      - 构建项目"
	@echo "  make swagger    - 生成 swagger 文档"
	@echo "  make clean      - 清理构建文件"
	@echo "  make test       - 运行测试"
	@echo "  make deps       - 安装依赖"
	@echo "  make update-deps - 更新依赖"
	@echo "  make fmt        - 格式化代码"
	@echo "  make lint       - 运行代码检查" 