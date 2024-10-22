package shared

import (
	"net/rpc"

	"github.com/hashicorp/go-plugin"
)

// Greeter 是我们作为插件暴露的接口。
type Greeter interface {
	Greet() string
}

// 这是一个通过RPC通信的实现
type GreeterRPC struct {
	client *rpc.Client
}

// Greet 实现了Greeter接口的Greet方法。
func (g *GreeterRPC) Greet() string {
	var resp string
	err := g.client.Call("Plugin.Greet", new(interface{}), &resp)
	if err != nil {
		// 通常情况下，接口应该返回错误。如果它们不返回，
		// 这里没有太多其他选择。
		panic(err)
	}

	return resp
}

// 这是GreeterRPC与之通信的RPC服务器，符合net/rpc的要求
type GreeterRPCServer struct {
	// 这是真正的实现
	Impl Greeter
}

// Greet 实现了GreeterRPCServer的Greet方法。
func (s *GreeterRPCServer) Greet(args interface{}, resp *string) error {
	*resp = s.Impl.Greet()
	return nil
}

// 这是plugin.Plugin的实现，以便我们可以提供/使用这个插件
//
// 这个实现有两个方法：Server必须返回一个此插件类型的RPC服务器。
// 我们为此构造了一个GreeterRPCServer。
//
// Client必须返回一个实现我们接口的实例，该实例通过RPC客户端进行通信。
// 我们为此返回GreeterRPC。
//
// 忽略MuxBroker。这是用来在我们的插件连接上创建更多多路复用流的高级用例。
type GreeterPlugin struct {
	// Impl 注入
	Impl Greeter
}

// Server 实现了GreeterPlugin的Server方法。
func (p *GreeterPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &GreeterRPCServer{Impl: p.Impl}, nil
}

// Client 实现了GreeterPlugin的Client方法。
func (GreeterPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &GreeterRPC{client: c}, nil
}
