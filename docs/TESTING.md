# 测试说明

## 单元测试

```shell
go test ./...
```

## 说明

- 根目录 `https2http` 包当前无 `_test.go` 文件。
- `cmd/proxychecker` 包含表驱动单元测试；新增逻辑请同步补充测试并保持 `go test ./...` 通过。
