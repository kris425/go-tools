package config

import (
	"testing"
	"time"
)

func TestNewViperConfig(t *testing.T) {
	v := NewViperConfig(".", "app")
	t.Log(v.GetString("app.Env"))
	v.AddNotifyFunc(func(name string) {
		t.Log(name)
	})
	time.Sleep(time.Second * 10)
}
