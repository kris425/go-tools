package ws

import (
	"testing"
)

func TestWSServer(t *testing.T) {
	ws := NewServer(":8881", nil)
	ws.Serve()
}
