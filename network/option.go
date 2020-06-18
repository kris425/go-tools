package network

import "github.com/kris425/go-tools/network/iface"

func NewConnOption(hook iface.IHook, packageMax uint32, codec iface.ICodec) *ConnOption {
	return &ConnOption{
		Hook:       hook,
		PackageMax: packageMax,
		Codec:      codec,
	}
}

type ConnOption struct {
	// 事件钩子接口
	Hook iface.IHook

	// 包最大长度
	PackageMax uint32

	// 编码器
	Codec iface.ICodec
}
type OptionFunc func(option *ConnOption)

// 设置Conn连接钩子
func SetHook(hook iface.IHook) OptionFunc {
	return func(option *ConnOption) {
		option.Hook = hook
	}
}

// 设置最大包长度
func SetPackageMax(len uint32) OptionFunc {
	return func(option *ConnOption) {
		option.PackageMax = len
	}
}

// 设置编码器
func SetCodec(codec iface.ICodec) OptionFunc {
	return func(option *ConnOption) {
		option.Codec = codec
	}
}
