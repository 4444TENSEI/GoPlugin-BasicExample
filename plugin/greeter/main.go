package main

import (
	"os"

	"kzplugin/shared"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
)

// GreeterHello 是 Greeter 接口的实际实现
type GreeterHello struct {
	logger hclog.Logger
}

// Greet 方法实现了 Greeter 接口，返回一个问候语
func (g *GreeterHello) Greet() string {
	g.logger.Debug("来自 GreeterHello.Greet 的消息")
	return "\n插件载入成功喵!\n"
}

// handshakeConfig 用于插件和宿主之间进行基本握手。
// 如果握手失败，将显示一个用户友好的错误。
// 这样可以防止用户执行错误的插件或执行一个插件目录。
// 这是一个用户体验特性，而不是一个安全特性。
var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,              // 握手协议版本
	MagicCookieKey:   "BASIC_PLUGIN", // 握手魔法饼干键
	MagicCookieValue: "hello",        // 握手魔法饼干值
}

func main() {
	// 创建一个日志记录器，配置为输出 JSON 格式的调试信息到标准错误输出
	logger := hclog.New(&hclog.LoggerOptions{
		Level:      hclog.Trace,
		Output:     os.Stderr,
		JSONFormat: true,
	})

	// 实例化 GreeterHello 并传递日志记录器
	greeter := &GreeterHello{
		logger: logger,
	}

	// pluginMap 是我们可以提供的插件映射
	var pluginMap = map[string]plugin.Plugin{
		"greeter": &shared.GreeterPlugin{Impl: greeter}, // 将 greeter 实现映射到 greeter 插件
	}

	// 记录一条调试信息
	logger.Debug("来自插件的消息", "foo", "bar")

	// 启动插件服务
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshakeConfig, // 插件握手配置
		Plugins:         pluginMap,       // 插件映射
	})
}
