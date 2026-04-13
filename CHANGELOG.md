# Changelog

## Unreleased

### Fixed

- 修复 CONNECT 隧道（HTTPS 流量）未经过上游代理的问题：启用 `proxy.ConnectDial = proxy.NewConnectDialToProxy(upstreamProxy)`，确保所有 HTTPS 目标的 CONNECT 请求都通过上游代理转发，而非直连目标服务器。此前在无直接出网能力的机器上会出现连接超时、curl 收到 502 的现象。([#1](https://github.com/LubyRuffy/https2http/issues/1))

### Changed

- GoReleaser 构建增加 `-buildvcs=false`；在 `AGENTS.md` 补充预编译体积与可选优化说明。

## [v0.1.0] - 2026-04-13

### Added

- 增加 GoReleaser 配置（`.goreleaser.yaml`）：为 macOS arm64、Linux amd64、Windows amd64 预编译 `https2http` 与 `proxychecker`，生成压缩包与校验和文件。
