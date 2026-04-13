# 开发者说明（AGENTS）

## 构建与运行

- 主程序：`go run . -proxy <上游HTTPS代理> -addr :8080`
- 辅助工具：`go run ./cmd/proxychecker <参数>`（参数说明见 `docs/proxychecker.md`）

## 测试

```shell
go test ./...
```

更细的测试说明见 `docs/TESTING.md`。

## GoReleaser 发布

- 校验配置：`goreleaser check`
- 本地快照（产物输出到 `dist/`，默认跳过校验与上传）：`goreleaser release --snapshot --clean --skip=publish`
- 正式 `goreleaser release` 前需有对应版本的 git tag，并按需配置 `GITHUB_TOKEN` 等发布管道

预编译目标平台由 `.goreleaser.yaml` 中 `builds` 的 `goos` / `goarch` 与 `ignore` 矩阵限定。

## 预编译体积与进一步优化

当前发布构建已采用常见「体积友好」组合：`CGO_ENABLED=0`、`-trimpath`、`-buildvcs=false`、`-ldflags="-s -w"`（去掉符号表与 DWARF）。在此之上：

- **依赖与代码路径**：体积主要来自标准库与第三方依赖（如 `goproxy`、FOFA 客户端等）；去掉未用功能、拆分发行物，通常比再抠链接参数更有效。
- **编译器取舍**：若可接受明显性能下降，可在对应 `build` 下增加 `gcflags: ["all=-l"]`（关闭内联）。在本仓库当前依赖下，`proxychecker` 与根目录程序实测可再缩小约 **9%~10%**，但代理/并发场景下吞吐会受影响，默认发布不启用。
- **二次压缩**：`UPX` 等可再显著缩小磁盘占用，但可能触发误报、并与部分平台签名流程不兼容，默认不采用。
- **发布包**：GoReleaser 对 `tar.gz` / `zip` 已使用较高压缩级别；下载体积主要仍由可执行文件决定。

## 约定

- Go 工具链版本以 `go.mod` 的 `go` 指令为准。
- 变更行为或发布方式时同步更新 `README.md`、`CHANGELOG.md` 与本文件。
