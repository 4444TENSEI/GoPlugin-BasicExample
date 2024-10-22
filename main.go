package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"kzplugin/shared"

	hclog "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"gopkg.in/yaml.v2"
)

// PluginConfig 定义了插件的配置结构
type PluginConfig struct {
	Path string `yaml:"path"` // 插件的路径
}

// Config 定义了配置文件的结构
type Config struct {
	Plugins []PluginConfig `yaml:"plugins"` // 插件列表
}

func main() {
	// 从YAML文件中加载配置
	config := loadConfig("config.yml")

	// 创建一个hclog.Logger
	logger := hclog.New(&hclog.LoggerOptions{
		Name:   "plugin",    // 日志名称
		Output: os.Stdout,   // 输出位置
		Level:  hclog.Debug, // 日志级别
	})

	// 遍历插件配置并启动每个插件
	for _, pluginCfg := range config.Plugins {
		// 首先启动插件进程
		client := plugin.NewClient(&plugin.ClientConfig{
			HandshakeConfig: handshakeConfig,              // 握手配置
			Plugins:         pluginMap,                    // 插件映射
			Cmd:             exec.Command(pluginCfg.Path), // 插件的执行命令
			Logger:          logger,                       // 日志记录器
		})
		defer client.Kill() // 程序退出前杀死插件进程

		// 通过RPC连接
		rpcClient, err := client.Client()
		if err != nil {
			log.Fatalf("为插件[%s]创建RPC客户端时出错: %s", pluginCfg.Path, err)
		}

		// 请求插件
		raw, err := rpcClient.Dispense("greeter") // 所有插件类型都是greeter
		if err != nil {
			log.Fatalf("分发插件%s时出错: %s", pluginCfg.Path, err)
		}

		// 现在我们应该有一个Greeter实例！这看起来像一个正常的接口实现，
		// 但实际上是通过RPC连接实现的。
		greeter := raw.(shared.Greeter)
		fmt.Printf("\n插件%s说: %s\n", pluginCfg.Path, greeter.Greet())
	}
}

// loadConfig 从给定路径加载配置
func loadConfig(path string) Config {
	var config Config
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("读取配置文件时出错: %s", err)
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("解析配置文件时出错: %s", err)
	}

	return config
}

// handshakeConfig 用于插件和宿主之间的基本握手。如果握手失败，将显示用户友好的错误。
var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,              // 协议版本
	MagicCookieKey:   "BASIC_PLUGIN", // 魔法饼干键
	MagicCookieValue: "hello",        // 魔法饼干值
}

// pluginMap 是我们可以分配的插件映射。
var pluginMap = map[string]plugin.Plugin{
	"greeter": &shared.GreeterPlugin{}, // 所有插件类型都是greeter
	// 在这里添加其他插件
}
