package main

import (
	"crypto/tls"
	"flag"
	"github.com/elazarl/goproxy"
	"log"
	"net/http"
	"net/url"
)

func main() {
	upstreamProxy := flag.String("proxy", "", "upstream proxy")
	addr := flag.String("addr", ":8080", "local addr, as http proxy")
	flag.Parse()

	proxyURL, err := url.Parse(*upstreamProxy)
	if err != nil {
		panic(err)
	}

	if proxyURL.Scheme == "http" {
		panic("don't waste time, it's already http proxy")
	}

	proxy := goproxy.NewProxyHttpServer()
	proxy.Tr = &http.Transport{
		Proxy:           http.ProxyURL(proxyURL),
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	//proxy.ConnectDial = proxy.NewConnectDialToProxy(*upstreamProxy)

	//proxy.OnRequest().DoFunc(
	//	func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	//		return r, nil
	//	})
	proxy.Verbose = true
	log.Fatal(http.ListenAndServe(*addr, proxy))
}
