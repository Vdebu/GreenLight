# 导入环境变量
include .envrc

# 检测操作系统并设置变量
ifeq ($(OS),Windows_NT)
    CONFIRM_CMD = powershell -Command "Write-Host 'Are you sure? [y/N]' -NoNewline; $$ans = Read-Host; if ($$ans.ToLower() -eq 'y') { exit 0 } else { exit 1 }"
else
    CONFIRM_CMD = read -p "Are you sure? [y/N] " ans && [ "$${ans}" = "y" ]
endif

# 声明所有伪目标
.PHONY: help confirm run/api db/psql db/migration/new db/migrations/up audit vendor build/api

## help: 展示帮助信息
help:
	@echo Usage:
	@echo   make help               - 显示帮助信息
	@echo   make run/api            - 启动 API
	@echo   make db/psql            - 连接数据库
	@echo   make db/migration/new   - 创建新迁移文件（需指定 name=迁移名称）
	@echo   make db/migrations/up   - 执行数据库迁移

## confirm: 确认操作（支持 y/Y）
confirm:
	@$(CONFIRM_CMD)

## run/api: 启动 API 并导入相应的环境变量
run/api:
	go run ./cmd/api -db-dsn=${GREENLIGHT_DB_DSN}

## db/psql: 连接数据库
db/psql:
	psql "${GREENLIGHT_DB_DSN}"

## db/migration/new: 创建新迁移文件（用法：make db/migration/new name=<迁移名称>）
db/migration/new:
	@echo Creating migration files for ${name}
	migrate create -seq -ext=.sql -dir=./migrations ${name}

## db/migrations/up: 执行数据库迁移（依赖 confirm 目标）
db/migrations/up: confirm
	@echo Running up migrations
	migrate -path ./migrations -database "${GREENLIGHT_DB_DSN}" up

## 整理代码结构与启动相关测试
audit:
	@echo Formatting code...
	go fmt ./...
	@echo Vetting code...
	go vet ./...
	staticcheck ./...
	@echo Running tests...
	go test -race -vet=off ./...

## 整理依赖并将项目中用到的依赖包备份
vendor:
	@echo Tidying and verifying module dependencies
	go mod tidy
	go mod verify
	@echo Vendoring dependencies...
	go mod vendor

## 存储当前时间
current_time=$(shell powershell -Command "Get-Date -Format o")
## 通过git获取相关版本信息
git_description=$(shell powershell -Command "git describe --always --dirty --tags --long")
## 存储交叉编译方法方便复用
linker_flags='-s -X main.buildTime=${current_time} -X main.version=${git_description}'

## 交叉编译二进制文件
build/api:
	@echo Building cmd/api
	go build -ldflags=${linker_flags} -o=./bin/api.exe ./cmd/api
	set GOOS=linux
	set GOARCH=amd64
	go build -ldflags=${linker_flags} -o=./bin/linux_amd64/api ./cmd/api