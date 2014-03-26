package bencode

import (
	"testing"
	"bytes"
)

func TestBoolMarshal(t *testing.T) {
	bt, b0 := true, []byte("i1e")
	bf, b1 := false, []byte("i0e")

	u0, err0 := Marshal(bt)
	if err0 != nil {
		t.Fatal(err0)
	}

	if bytes.Compare(u0, b0) != 0 {
		t.Errorf("invalid data from marshalling; got %s, expected %s", string(u0), string(b0))
	}

	u1, err1 := Marshal(bf)
	if err1 != nil {
		t.Fatal(err1)
	}

	if bytes.Compare(u1, b1) != 0 {
		t.Errorf("invalid data from marshalling; got %s, expected %s", string(u1), string(b1))
	}
}

func TestIntMarshal(t *testing.T) {
	i := -65536
	b := []byte("i-65536e")

	m, err := Marshal(i)
	if err != nil {
		t.Fatal(err)
	}

	if bytes.Compare(b, m) != 0 {
		t.Errorf("invalid data from marshalling; got %s, expected %s", string(m), string(b))
	}
}

func TestUintMarshal(t *testing.T) {
	i := 2048
	b := []byte("i2048e")

	m, err := Marshal(i)
	if err != nil {
		t.Fatal(err)
	}

	if bytes.Compare(b, m) != 0 {
		t.Errorf("invalid data from marshalling; got %s, expected %s", string(m), string(b))
	}
}

func TestArrayMarshal(t *testing.T) {
	i := []uint{1, 2, 3}
	b := []byte("li1ei2ei3ee")

	m, err := Marshal(i)
	if err != nil {
		t.Fatal(err)
	}

	if bytes.Compare(b, m) != 0 {
		t.Errorf("invalid data from marshalling; got %s, expected %s", string(m), string(b))
	}
}

func TestMapMarshal(t *testing.T) {
	i := map[string]uint{
		"a": 1,
		"b": 2,
		"c": 3,
	}
	b := []byte("d1:ai1e1:bi2e1:ci3ee")

	m, err := Marshal(i)
	if err != nil {
		t.Fatal(err)
	}

	if bytes.Compare(b, m) != 0 {
		t.Errorf("invalid data from marshalling; got %s, expected %s", string(m), string(b))
	}
}

type testSimpleStruct struct {
	A uint
	B *uint
	S0 testInnerStruct
	S1 *testInnerStruct
}

type testInnerStruct struct {
	C string
}

func TestSimpleStructMarshal(t *testing.T) {
	i := new(testSimpleStruct)
	i.A = 1
	i.B = new(uint)
	*i.B = 2
	i.S0.C = "a"
	i.S1 = new(testInnerStruct)
	i.S1.C = "b"

	b := []byte("d1:Ai1e1:Bi2e2:S0d1:C1:ae2:S1d1:C1:bee")

	m, err := Marshal(i)
	if err != nil {
		t.Fatal(err)
	}

	if bytes.Compare(b, m) != 0 {
		t.Errorf("invalid data from marshalling; got %s, expected %s", string(m), string(b))
	}
}

func TestStringMarshal(t *testing.T) {
	str := "foo"
	b := []byte("3:foo")

	m, err := Marshal(str)
	if err != nil {
		t.Fatal(err)
	}

	if bytes.Compare(m, b) != 0 {
		t.Errorf("invalid data from marshalling; got %s, expected %s", string(m), string(b))
	}
}
