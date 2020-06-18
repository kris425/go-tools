package iface

type IHook interface {
	OnRecv(conn IConn, data []byte)

	OnNewConn(conn IConn)

	OnClosed(conn IConn)
}
