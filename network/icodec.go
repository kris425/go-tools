package network

type ICodec interface {
	Encode(data []byte) ([]byte, error)

	Decode(conn IConn) ([]byte, error)
}
