# https2http
把https的代理变成http的代理。

## 背景
有特别多的代理软件不支持https，只支持http，比如在windows下面有一个最好用的全局代理软件Proxifier，它就不支持最简单的https代理服务器。

## 方案
为了支持其他一些代理软件，我们需要进行一个转换，upstream的是https，本地是http，这样就能够完成网络请求流程了。

```
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

