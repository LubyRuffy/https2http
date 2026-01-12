# https2http

把https的代理变成http的代理。

## 背景

有特别多的代理软件不支持https，只支持http，比如在windows下面有一个最好用的全局代理软件Proxifier，它就不支持最简单的https代理服务器。

## 方案

为了支持其他一些代理软件，我们需要进行一个转换，upstream的是https，本地是http，这样就能够完成网络请求流程了。

```text
app <-> proxy(Proxifier) <-> https2http <-> https proxy
```

## 运行

```shell
https2http -proxy https://proxy.xxx.com -addr :8080
```

直接用gost也可以：

```shell
git clone https://github.com/ginuerzh/gost.git
cd gost/cmd/gost
go build
./gost.exe -L=:8080 -F=https://proxy.xxx.com
```

## 辅助工具

- 通过fofa提取代理：

```shell
proxychecker -query 'type="service" && protocol="http" && banner="ERR_INVALID_URL"' -expr 'response.Header("Server")=~"(?is)(bws|bfe)"' -target https://www.baidu.com -size 100
proxychecker -query 'port="3128" && title="ERROR: The requested URL could not be retrieved"' -expr 'response.Header("Server")=~"(?is)(bws|bfe)"' -target https://www.baidu.com -size 100
```

http代理测试，这种规则基本都在中国的代理：[port="9091" && banner="403 Forbidden" && banner="nginx/1.12.1"](https://fofa.info/result?qbase64=cG9ydD05MDkxICYmIGJhbm5lcj0iNDAzIEZvcmJpZGRlbiIgJiYgYmFubmVyPSJuZ2lueC8xLjEyLjEi)
```shell
proxychecker -query 'port=9091 && banner="403 Forbidden" && banner="nginx/1.12.1"' -expr 'response.Body()=~"(?is)百度一下"' -target http://www.baidu.com -size 1000
```

检查Mikrotik代理: [banner="Mikrotik HttpProxy"](https://fofa.info/result?qbase64=YmFubmVyPSJNaWtyb3RpayBIdHRwUHJveHki)

```shell
proxychecker -query 'banner="Mikrotik HttpProxy"' -expr 'response.Body()=~"(?is)百度一下"' -target http://www.baidu.com -size 1000
```

检查socks5代理: [banner="Authentication(0x00)"](https://fofa.info/result?qbase64=YmFubmVyPSJBdXRoZW50aWNhdGlvbigweDAwKSI%3D)

```shell
proxychecker -query 'banner="Authentication(0x00)"' -expr 'response.Body()=~"(?is)百度一下"' -target http://www.baidu.com -size 1000 -type socks5
```

检查socks5代理: [body="This is a proxy server. Does not respond to non-proxy requests."](https://fofa.info/result?qbase64=Ym9keT0iVGhpcyBpcyBhIHByb3h5IHNlcnZlci4gRG9lcyBub3QgcmVzcG9uZCB0byBub24tcHJveHkgcmVxdWVzdHMuIg%3D%3D)

```shell
# 应该是通用组件，在请求www.baidu.com的情况下会提示错误：dial tcp: address [2405:19c0:c303:423d:deef:db8e:1c8d:ecb0]:0: no suitable address found
proxychecker -query 'body="This is a proxy server. Does not respond to non-proxy requests."' -expr 'response.Body()=~"(?is)百度一下"' -target http://www.baidu.com -size 1000 -type http

proxychecker -query 'body="This is a proxy server. Does not respond to non-proxy requests."' -expr 'response.Body()=~"(?is)GeoNameID"' -target http://ip.bmh.im/h -size 1000 -type http
```

## 地理信息查询

proxychecker 支持使用 `-geo` 参数获取有效代理的国家代码：

```shell
# 使用 geo 参数获取代理的国家代码
proxychecker -query 'type="subdomain" && cert.is_valid=true && domain!="" && title="ERROR: The requested URL could not be retrieved"' -expr 'response.Header("Server")=="gws"' -target https://www.google.com -size 100 -geo
```

### 输出示例

```text
successful proxy: http://1.1.1.1:8080, country: CN, ip: 1.2.3.4
```

### 功能说明

- 当使用 `-geo` 参数时，proxychecker 会在发现有效代理后，通过该代理访问 `http://ip.bmh.im/c` 获取地理信息
- 显示国家代码和出口 IP，格式简洁易读
- 支持各种类型的代理（http/https/socks5）
- 可以结合其他参数使用，如 `-type`、`-timeout` 等

## 生成 Clash 配置文件

proxychecker 支持使用 `-clash` 参数自动生成 Clash 配置文件：

```shell
# 检测代理并生成 Clash 配置
proxychecker -query 'port="3128"' -expr 'response.Body()=~"(?is)百度"' -target https://www.baidu.com -size 100 -clash clash.yaml

# 结合地理信息，代理名称会包含国家代码
proxychecker -query 'port="3128"' -expr 'response.Body()=~"(?is)百度"' -target https://www.baidu.com -size 100 -geo -clash clash.yaml

# 自定义代理组名称
proxychecker -query 'port="3128"' -expr 'response.Body()=~"(?is)百度"' -target https://www.baidu.com -size 100 -clash clash.yaml -clashGroup "my-proxies"
```

### 生成的配置示例

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
proxy-groups:
  - name: proxy
    type: select
    proxies:
      - US-1
      - JP-v6-2
```

### 功能说明

- 使用 `-clash` 参数指定输出文件路径
- 使用 `-clashGroup` 参数自定义代理组名称（默认为 `proxy`）
- 代理命名规则：
  - 启用 `-geo` 时：`国家代码-序号`（如 `US-1`），IPv6 代理为 `国家代码-v6-序号`
  - 未启用 `-geo` 时：`proxy-序号`
- 自动识别代理类型（http/socks5）和 TLS 设置

更多详细说明请参考 [docs/proxychecker.md](docs/proxychecker.md)
