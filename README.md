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

## 辅助工具
- 通过fofa提取代理：
```shell
proxychecker -query 'type="service" && protocol="http" && banner="ERR_INVALID_URL"' -expr 'response.Header("Server")=~"(?is)(bws|bfe)"' -target https://www.baidu.com -size 100
proxychecker -query 'port="3128" && title="ERROR: The requested URL could not be retrieved"' -expr 'response.Header("Server")=~"(?is)(bws|bfe)"' -target https://www.baidu.com -size 100
```