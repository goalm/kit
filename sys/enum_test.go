package sys

import "testing"

func TestNewEnum(t *testing.T) {
	enum := NewEnum()
	if enum == nil {
		t.Error("NewEnum() returned nil")
	}
}

func TestEnum_Add(t *testing.T) {
	enum := NewEnum()
	enum.Add(1, "name", "desc")
	if enum.Size() != 1 {
		t.Error("Size() should return 1")
	}

	if idx, ok := enum.StrToInt("name"); !ok || idx != 1 {
		t.Error("StrToInt() failed")
	}

	if name, ok := enum.IntToStr(1); !ok || name != "name" {
		t.Error("IntToStr() failed")
	}
}
