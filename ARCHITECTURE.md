# 架构说明

## 系统整体

本仓库在同一 Go module 下提供两个命令行程序：`https2http`（根目录入口）与 `proxychecker`（`cmd/proxychecker`）。二者无运行时耦合，可分别交叉编译与分发。

## 核心模块

- **https2http**：基于 `goproxy` 的本地 HTTP 代理，把上游 HTTPS 代理对接到仅支持 HTTP 代理的客户端链路中。
- **proxychecker**：代理发现、探测与结果导出（Clash 配置、地理信息等），依赖 FOFA 等外部服务；细节见 `docs/proxychecker.md`。

## 发布与制品流

GoReleaser 根据 `.goreleaser.yaml` 执行 `go mod tidy` 前置钩子，在 `CGO_ENABLED=0` 下为指定 `GOOS`/`GOARCH` 组合构建两个二进制，按平台打包为 `tar.gz` 或 `zip`，并生成聚合校验和文件，供下载校验与发行说明使用。
