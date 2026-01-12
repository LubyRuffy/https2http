# ProxyChecker 工具文档

ProxyChecker 是一个用于检测和验证代理服务器的命令行工具。它可以从 FOFA 搜索引擎获取代理列表，并验证这些代理是否能够正常工作。

## 安装

```shell
cd cmd/proxychecker
go build
```

## 命令行参数

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `-query` | `type="subdomain" && cert.is_valid=true && domain!="" && title="ERROR: The requested URL could not be retrieved"` | FOFA 查询语句 |
| `-expr` | `response.Header("Server")=="gws"` | 验证代理的表达式 |
| `-target` | `https://www.google.com` | 用于测试代理的目标 URL |
| `-testProxy` | - | 直接测试单个代理（不使用 FOFA 搜索） |
| `-method` | `GET` | HTTP 请求方法 |
| `-type` | `auto` | 代理类型：`socks5`/`http`/`https`/`auto` |
| `-timeout` | `10` | 请求超时时间（秒） |
| `-workers` | `20` | 并发工作线程数 |
| `-size` | `1000` | FOFA 搜索结果数量 |
| `-debug` | `false` | 启用调试模式 |
| `-geo` | `false` | 获取有效代理的地理位置信息 |
| `-clash` | - | 输出 Clash 配置文件路径（如：`clash.yaml`） |
| `-clashGroup` | `proxy` | Clash 代理组名称 |

## 功能说明

### 基本代理检测

通过 FOFA 搜索代理服务器，并验证它们是否能够正常访问指定目标：

```shell
proxychecker -query 'port="3128"' -expr 'response.Header("Server")=~"(?is)(nginx)"' -target https://www.baidu.com -size 100
```

### 地理信息查询

使用 `-geo` 参数可以获取有效代理的国家代码和出口 IP：

```shell
proxychecker -query 'port="3128"' -expr 'response.Body()=~"(?is)百度"' -target https://www.baidu.com -size 100 -geo
```

输出示例：
```
successful proxy: http://1.1.1.1:8080, country: CN, ip: 1.2.3.4
successful proxy: https://2.2.2.2:443, country: US, ip: 5.6.7.8 [IPv6]
```

### 生成 Clash 配置文件

使用 `-clash` 参数可以自动将有效的代理生成 Clash 配置文件：

```shell
# 基本用法
proxychecker -testProxy https://proxy.example.com:443 -expr 'response.Header("Server")=="gws"' -target https://www.google.com -clash clash.yaml

# 结合地理信息，代理名称会包含国家代码
proxychecker -query 'port="3128"' -expr 'response.Body()=~"(?is)百度"' -target https://www.baidu.com -size 100 -geo -clash proxies.yaml

# 自定义代理组名称
proxychecker -query 'port="3128"' -expr 'response.Body()=~"(?is)百度"' -target https://www.baidu.com -size 100 -clash clash.yaml -clashGroup "my-proxies"
```

生成的 Clash 配置文件示例：

```yaml
proxies:
  - name: US-1
    type: http
    server: proxy1.example.com
    port: 443
    tls: true
  - name: JP-v6-2
    type: http
    server: proxy2.example.com
    port: 8080
  - name: proxy-3
    type: socks5
    server: proxy3.example.com
    port: 1080
proxy-groups:
  - name: proxy
    type: select
    proxies:
      - US-1
      - JP-v6-2
      - proxy-3
```

#### Clash 配置说明

- **代理命名规则**：
  - 如果启用了 `-geo` 参数且获取到地理信息，名称格式为 `国家代码-序号`（如 `US-1`）
  - 如果是 IPv6 代理，名称格式为 `国家代码-v6-序号`（如 `JP-v6-2`）
  - 如果没有地理信息，名称格式为 `proxy-序号`（如 `proxy-3`）

- **代理类型识别**：
  - `http://` 或 `https://` 前缀的代理会被识别为 `http` 类型
  - `socks5://` 前缀的代理会被识别为 `socks5` 类型
  - `https://` 前缀的代理会自动启用 `tls: true`

- **代理组**：
  - 默认生成一个名为 `proxy` 的选择组（可通过 `-clashGroup` 参数自定义）
  - 代理组类型为 `select`，包含所有有效代理

## 使用示例

### 检测 Squid 代理

```shell
proxychecker -query 'port="3128" && title="ERROR: The requested URL could not be retrieved"' -expr 'response.Header("Server")=~"(?is)(bws|bfe)"' -target https://www.baidu.com -size 100
```

### 检测 Mikrotik 代理

```shell
proxychecker -query 'banner="Mikrotik HttpProxy"' -expr 'response.Body()=~"(?is)百度一下"' -target http://www.baidu.com -size 1000
```

### 检测 SOCKS5 代理

```shell
proxychecker -query 'banner="Authentication(0x00)"' -expr 'response.Body()=~"(?is)百度一下"' -target http://www.baidu.com -size 1000 -type socks5
```

### 测试单个代理并生成配置

```shell
proxychecker -testProxy https://proxy.example.com:443 -expr 'response.Header("Server")=="gws"' -target https://www.google.com -clash single-proxy.yaml
```

### 批量检测并生成 Clash 配置

```shell
proxychecker -query 'port="3128"' -expr 'response.Body()=~"(?is)百度"' -target https://www.baidu.com -size 500 -workers 50 -geo -clash all-proxies.yaml -clashGroup "fofa-proxies"
```
