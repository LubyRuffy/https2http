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

## 约定

- Go 工具链版本以 `go.mod` 的 `go` 指令为准。
- 变更行为或发布方式时同步更新 `README.md`、`CHANGELOG.md` 与本文件。
