package encrypt

import "testing"

func TestMD5(t *testing.T) {
	s := "abc"
	md5Str := MD5(s, nil)
	if md5Str != "900150983cd24fb0d6963f7d28e17f72" {
		t.Failed()
	}
	t.Log(md5Str)
}

func TestSha1(t *testing.T) {
	s := "abc"
	str := Sha1([]byte(s), nil)
	if str != "a9993e364706816aba3e25717850c26c9cd0d89d" {
		t.Failed()
	}
	t.Log(str)
}

func TestBase64Encode(t *testing.T) {
	s := "abc"
	str := Base64Encode([]byte(s))
	if str != "YWJj" {
		t.Failed()
	}
	t.Log(str)
}

func TestBase64Decode(t *testing.T) {
	s := "YWJj"
	str, err := Base64Decode(s)
	if err != nil {
		t.Error(err)
		t.Failed()
		return
	}
	if string(str) != "abc" {
		t.Failed()
	}
	t.Log(str)
}

func TestAESEncrypt(t *testing.T) {
	s := "hello"
	str := Base64Encode(AESEncryptECB([]byte(s), []byte("123456")))
	if str != "yPILSUwmVwvFC+Y0QHMGtA==" {
		t.Error(str)
		t.Failed()
	}
	t.Log(str)
}

func TestAESDecrypt(t *testing.T) {
	s := "yPILSUwmVwvFC+Y0QHMGtA=="
	data, err := Base64Decode(s)
	if err != nil {
		t.Error(err)
		t.Failed()
		return
	}
	str := AESDecryptECB(data, []byte("123456"))
	if string(str) != "hello" {
		t.Error(str)
		t.Failed()
	}
	t.Log(string(str))
}

func TestDESEncryptECB(t *testing.T) {
	s := "hello"
	data, err := DESEncryptECB([]byte(s), []byte("12345678"))
	if err != nil {
		t.Error(err)
		t.Failed()
	}
	str := Base64Encode(data)
	if str != "uhbGoCVxJa8=" {
		t.Error(str)
		t.Failed()
		return
	}
	t.Log(str)
}

func TestDESDecryptECB(t *testing.T) {
	s := "uhbGoCVxJa8="
	data, err := Base64Decode(s)
	if err != nil {
		t.Error(err)
		t.Failed()
		return
	}
	str, err := DESDecryptECB(data, []byte("12345678"))
	if err != nil {
		t.Error(err)
		t.Failed()
	}
	if string(str) != "hello" {
		t.Error(str)
		t.Failed()
		return
	}
	t.Log(string(str))
}

func TestHexEncode(t *testing.T) {
	s := "abc"
	str := HexEncode([]byte(s))
	if str != "616263" {
		t.Error(str)
		t.Failed()
		return
	}
	t.Log(str)
}
