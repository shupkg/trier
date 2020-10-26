package strs

import (
	"bytes"
	"testing"
)

func TestUnsafe(t *testing.T) {
	var s = "我是中国人"
	b := UnsafeBytes(s)
	if !bytes.Equal(b, []byte(s)) {
		t.Fail()
	}
	s2 := UnsafeString(b)
	if s2 != s {
		t.Fail()
	}
}
