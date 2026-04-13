# Changelog

## Unreleased

### Changed

- GoReleaser 构建增加 `-buildvcs=false`；在 `AGENTS.md` 补充预编译体积与可选优化说明。

## [v0.1.0] - 2026-04-13

### Added

- 增加 GoReleaser 配置（`.goreleaser.yaml`）：为 macOS arm64、Linux amd64、Windows amd64 预编译 `https2http` 与 `proxychecker`，生成压缩包与校验和文件。
