package main

import (
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseProxyURL(t *testing.T) {
	tests := []struct {
		name          string
		proxyURL      string
		wantServer    string
		wantPort      int
		wantTLS       bool
		wantProxyType string
		wantErr       bool
	}{
		{
			name:          "https with default port",
			proxyURL:      "https://example.com",
			wantServer:    "example.com",
			wantPort:      443,
			wantTLS:       true,
			wantProxyType: "http",
			wantErr:       false,
		},
		{
			name:          "https with custom port",
			proxyURL:      "https://example.com:8443",
			wantServer:    "example.com",
			wantPort:      8443,
			wantTLS:       true,
			wantProxyType: "http",
			wantErr:       false,
		},
		{
			name:          "http with default port",
			proxyURL:      "http://example.com",
			wantServer:    "example.com",
			wantPort:      80,
			wantTLS:       false,
			wantProxyType: "http",
			wantErr:       false,
		},
		{
			name:          "http with custom port",
			proxyURL:      "http://example.com:8080",
			wantServer:    "example.com",
			wantPort:      8080,
			wantTLS:       false,
			wantProxyType: "http",
			wantErr:       false,
		},
		{
			name:          "socks5 with default port",
			proxyURL:      "socks5://example.com",
			wantServer:    "example.com",
			wantPort:      1080,
			wantTLS:       false,
			wantProxyType: "socks5",
			wantErr:       false,
		},
		{
			name:          "socks5 with custom port",
			proxyURL:      "socks5://example.com:1081",
			wantServer:    "example.com",
			wantPort:      1081,
			wantTLS:       false,
			wantProxyType: "socks5",
			wantErr:       false,
		},
		{
			name:          "IPv6 address",
			proxyURL:      "https://[2001:db8::1]:443",
			wantServer:    "2001:db8::1",
			wantPort:      443,
			wantTLS:       true,
			wantProxyType: "http",
			wantErr:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, port, tls, proxyType, err := ParseProxyURL(tt.proxyURL)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.wantServer, server)
			assert.Equal(t, tt.wantPort, port)
			assert.Equal(t, tt.wantTLS, tls)
			assert.Equal(t, tt.wantProxyType, proxyType)
		})
	}
}

func TestGenerateProxyName(t *testing.T) {
	tests := []struct {
		name     string
		index    int
		geoInfo  *GeoInfo
		expected string
	}{
		{
			name:     "without geo info",
			index:    1,
			geoInfo:  nil,
			expected: "proxy-1",
		},
		{
			name:     "with geo info IPv4",
			index:    2,
			geoInfo:  &GeoInfo{Country: "US", IP: "1.2.3.4"},
			expected: "US-2",
		},
		{
			name:     "with geo info IPv6",
			index:    3,
			geoInfo:  &GeoInfo{Country: "JP", IP: "2001:db8::1"},
			expected: "JP-v6-3",
		},
		{
			name:     "with empty country",
			index:    4,
			geoInfo:  &GeoInfo{Country: "", IP: "1.2.3.4"},
			expected: "proxy-4",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GenerateProxyName(tt.index, tt.geoInfo)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGeoInfo_IsIPv6(t *testing.T) {
	tests := []struct {
		name     string
		geoInfo  GeoInfo
		expected bool
	}{
		{
			name:     "IPv4 address",
			geoInfo:  GeoInfo{IP: "192.168.1.1"},
			expected: false,
		},
		{
			name:     "IPv6 address",
			geoInfo:  GeoInfo{IP: "2001:db8::1"},
			expected: true,
		},
		{
			name:     "invalid IP",
			geoInfo:  GeoInfo{IP: "invalid"},
			expected: false,
		},
		{
			name:     "empty IP",
			geoInfo:  GeoInfo{IP: ""},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.geoInfo.IsIPv6())
		})
	}
}

func TestValidProxyCollector(t *testing.T) {
	collector := &ValidProxyCollector{}

	// 测试初始状态为空
	assert.Empty(t, collector.GetAll())

	// 添加代理
	collector.Add("http://proxy1.com:8080", nil)
	collector.Add("https://proxy2.com:443", &GeoInfo{Country: "US", IP: "1.2.3.4"})

	// 获取所有代理
	proxies := collector.GetAll()
	assert.Len(t, proxies, 2)
	assert.Equal(t, "http://proxy1.com:8080", proxies[0].Host)
	assert.Nil(t, proxies[0].GeoInfo)
	assert.Equal(t, "https://proxy2.com:443", proxies[1].Host)
	assert.NotNil(t, proxies[1].GeoInfo)
	assert.Equal(t, "US", proxies[1].GeoInfo.Country)
}

func TestGenerateClashConfig(t *testing.T) {
	tests := []struct {
		name      string
		proxies   []ValidProxy
		groupName string
		wantErr   bool
		validate  func(t *testing.T, config *ClashConfig)
	}{
		{
			name:      "empty proxies",
			proxies:   []ValidProxy{},
			groupName: "test",
			wantErr:   true,
		},
		{
			name: "single http proxy",
			proxies: []ValidProxy{
				{Host: "http://example.com:8080", GeoInfo: nil},
			},
			groupName: "proxy",
			wantErr:   false,
			validate: func(t *testing.T, config *ClashConfig) {
				assert.Len(t, config.Proxies, 1)
				assert.Equal(t, "proxy-1", config.Proxies[0].Name)
				assert.Equal(t, "http", config.Proxies[0].Type)
				assert.Equal(t, "example.com", config.Proxies[0].Server)
				assert.Equal(t, 8080, config.Proxies[0].Port)
				assert.False(t, config.Proxies[0].TLS)

				assert.Len(t, config.ProxyGroups, 1)
				assert.Equal(t, "proxy", config.ProxyGroups[0].Name)
				assert.Equal(t, "select", config.ProxyGroups[0].Type)
				assert.Equal(t, []string{"proxy-1"}, config.ProxyGroups[0].Proxies)
			},
		},
		{
			name: "single https proxy",
			proxies: []ValidProxy{
				{Host: "https://example.com:443", GeoInfo: &GeoInfo{Country: "US", IP: "1.2.3.4"}},
			},
			groupName: "myproxy",
			wantErr:   false,
			validate: func(t *testing.T, config *ClashConfig) {
				assert.Len(t, config.Proxies, 1)
				assert.Equal(t, "US-1", config.Proxies[0].Name)
				assert.Equal(t, "http", config.Proxies[0].Type)
				assert.Equal(t, "example.com", config.Proxies[0].Server)
				assert.Equal(t, 443, config.Proxies[0].Port)
				assert.True(t, config.Proxies[0].TLS)

				assert.Equal(t, "myproxy", config.ProxyGroups[0].Name)
			},
		},
		{
			name: "multiple proxies",
			proxies: []ValidProxy{
				{Host: "http://proxy1.com:80", GeoInfo: nil},
				{Host: "https://proxy2.com:443", GeoInfo: &GeoInfo{Country: "JP", IP: "1.2.3.4"}},
				{Host: "socks5://proxy3.com:1080", GeoInfo: &GeoInfo{Country: "KR", IP: "2001:db8::1"}},
			},
			groupName: "all",
			wantErr:   false,
			validate: func(t *testing.T, config *ClashConfig) {
				assert.Len(t, config.Proxies, 3)
				assert.Equal(t, "proxy-1", config.Proxies[0].Name)
				assert.Equal(t, "JP-2", config.Proxies[1].Name)
				assert.Equal(t, "KR-v6-3", config.Proxies[2].Name)

				assert.Equal(t, "http", config.Proxies[0].Type)
				assert.Equal(t, "http", config.Proxies[1].Type)
				assert.Equal(t, "socks5", config.Proxies[2].Type)

				assert.False(t, config.Proxies[0].TLS)
				assert.True(t, config.Proxies[1].TLS)
				assert.False(t, config.Proxies[2].TLS)

				assert.Len(t, config.ProxyGroups, 1)
				assert.Equal(t, []string{"proxy-1", "JP-2", "KR-v6-3"}, config.ProxyGroups[0].Proxies)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := GenerateClashConfig(tt.proxies, tt.groupName)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, config)
			if tt.validate != nil {
				tt.validate(t, config)
			}
		})
	}
}

func TestSaveClashConfig(t *testing.T) {
	config := &ClashConfig{
		Proxies: []ClashProxy{
			{
				Name:   "proxy1",
				Type:   "http",
				Server: "example.com",
				Port:   443,
				TLS:    true,
			},
		},
		ProxyGroups: []ClashProxyGroup{
			{
				Name:    "test",
				Type:    "select",
				Proxies: []string{"proxy1"},
			},
		},
	}

	// 创建临时目录
	tmpDir := t.TempDir()
	filename := filepath.Join(tmpDir, "test-clash.yaml")

	// 保存配置
	err := SaveClashConfig(config, filename)
	require.NoError(t, err)

	// 验证文件存在
	_, err = os.Stat(filename)
	require.NoError(t, err)

	// 读取并验证内容
	content, err := os.ReadFile(filename)
	require.NoError(t, err)

	assert.Contains(t, string(content), "proxies:")
	assert.Contains(t, string(content), "name: proxy1")
	assert.Contains(t, string(content), "type: http")
	assert.Contains(t, string(content), "server: example.com")
	assert.Contains(t, string(content), "port: 443")
	assert.Contains(t, string(content), "tls: true")
	assert.Contains(t, string(content), "proxy-groups:")
	assert.Contains(t, string(content), "name: test")
	assert.Contains(t, string(content), "type: select")
}

func TestFixURL(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		proxyType string
		expected  string
	}{
		{
			name:      "ip:port with port 80",
			input:     "192.168.1.1:80",
			proxyType: "auto",
			expected:  "http://192.168.1.1",
		},
		{
			name:      "ip:port with port 443",
			input:     "192.168.1.1:443",
			proxyType: "auto",
			expected:  "https://192.168.1.1",
		},
		{
			name:      "ip:port with custom port",
			input:     "192.168.1.1:8080",
			proxyType: "auto",
			expected:  "http://192.168.1.1:8080",
		},
		{
			name:      "ip:port with explicit socks5 type",
			input:     "192.168.1.1:1080",
			proxyType: "socks5",
			expected:  "socks5://192.168.1.1:1080",
		},
		{
			name:      "full http url",
			input:     "http://example.com:8080",
			proxyType: "auto",
			expected:  "http://example.com:8080",
		},
		{
			name:      "full https url with standard port",
			input:     "https://example.com:443",
			proxyType: "auto",
			expected:  "https://example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FixURL(tt.input, tt.proxyType)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestResponsePackage(t *testing.T) {
	// 创建模拟的 http.Response
	resp := &http.Response{
		Header: make(http.Header),
	}
	resp.Header.Set("Content-Type", "text/html")
	resp.Header.Set("Server", "nginx")

	body := []byte("Hello, World!")

	rp := NewResponsePackage(resp, body)

	// 测试 Header 方法
	assert.Equal(t, "text/html", rp.Header("Content-Type"))
	assert.Equal(t, "nginx", rp.Header("Server"))
	assert.Equal(t, "", rp.Header("Non-Existent"))

	// 测试 Body 方法
	assert.Equal(t, "Hello, World!", rp.Body())
}
