package bencode

import "testing"

func TestIntUnmarshal(t *testing.T) {
	b := []byte("i65536e")
	var i int64 = 0

	err := Unmarshal(b, &i)
	if err != nil {
		t.Fatal(err)
	}

	if i != 65536 {
		t.Errorf("Unmarshal failed for 65536, got: %d", i)
	}
}

func TestUintUnmarshal(t *testing.T) {
	b := []byte("i2048e")
	var i uint64 = 0

	err := Unmarshal(b, &i)
	if err != nil {
		t.Fatal(err)
	}

	if i != 2048 {
		t.Errorf("Unmarshal failed for 65536, got: %d", i)
	}
}

func TestStringUnmarshal(t *testing.T) {
	b := []byte("4:spam")
	str := ""

	err := Unmarshal(b, &str)
	if err != nil {
		t.Fatal(err)
	}

	if str != "spam" {
		t.Errorf("Unmarshal failed for spam, got: %s", str)
	}
}

func TestArrayUnmarshal(t *testing.T) {
	b := []byte("li1ei2ei3ee")
	var vals []uint64

	err := Unmarshal(b, &vals)
	if err != nil {
		t.Fatal(err)
	}

	if len(vals) != 3 {
		t.Fatalf("invalid array length, got %d, expected 3\n", len(vals))
	}

	for i := 0; i < 3; i++ {
		if vals[i] != uint64(i+1) {
			t.Errorf("invalid array value, got %u, expected %u", vals[i], i+1)
		}
	}
}

func TestMapUnmarshal(t *testing.T) {
	b := []byte("d1:ai1e1:bi2e1:ci3ee")
	control := make(map[string]uint)
	control["a"] = 1
	control["b"] = 2
	control["c"] = 3
	var m map[string]uint

	err := Unmarshal(b, &m)
	if err != nil {
		t.Fatal(err)
	}

	for k, v := range control {
		if val, ok := m[k]; !ok {
			t.Errorf("map does not have the value %d at key %s", v, k)
		} else {
			if val != v {
				t.Errorf("invalid array in map, got %d, expected %d at %s", val, v, k)
			}
		}
	}
}

type testUnmarshalStruct struct {
	A int
	B string
}

func TestStructUnmarshal(t *testing.T) {
	b := []byte("d1:ai1e1:b2:abe")
	s := new(testUnmarshalStruct)

	err := Unmarshal(b, s)
	if err != nil {
		t.Fatal(err)
	}

	if s.A != 1 {
		t.Errorf("invalid value in struct, got %d, expected 1", s.A)
	}

	if s.B != "ab" {
		t.Errorf("invalid value in struct, got %s, expected ab", s.B)
	}
}

type testUnmarshalStructAggregated struct {
	A int
	B string
}

type testUnmarshalComplexStruct struct {
	testUnmarshalStruct
	C []uint
	D testUnmarshalStructAggregated
	E *testUnmarshalStructAggregated
}

func TestComplexStructUnmarshal(t *testing.T) {
	b := []byte("d1:ai1e1:b2:ab1:cli1ei2ei3ee1:dd1:ai1e1:b2:abe1:ed1:ai1e1:b2:abee")
	s := new(testUnmarshalComplexStruct)

	err := Unmarshal(b, s)
	if err != nil {
		t.Fatal(err)
	}

	if s.A != 1 {
		t.Errorf("invalid value in struct, got %d, expected 1", s.A)
	}

	if s.B != "ab" {
		t.Errorf("invalid value in struct, got %s, expected ab", s.B)
	}

	if len(s.C) != 3 {
		t.Fatalf("invalid array length, got %d, expected 3\n", len(s.C))
	}

	for i := 0; i < 3; i++ {
		if s.C[i] != uint(i + 1) {
			t.Errorf("invalid array value, got %u, expected %u", s.C[i], i+1)
		}
	}

	if s.D.A != 1 {
		t.Errorf("invalid value in struct, got %d, expected 1", s.A)
	}

	if s.D.B != "ab" {
		t.Errorf("invalid value in struct, got %s, expected ab", s.B)
	}

	if s.E != nil {
		if s.E.A != 1 {
			t.Errorf("invalid value in struct, got %d, expected 1", s.A)
		}

		if s.E.B != "ab" {
			t.Errorf("invalid value in struct, got %s, expected ab", s.B)
		}
	} else {
		t.Error("pointer field is nil in struct")
	}
}

type testNonstandardNameStruct struct {
	AA int
	BB int
}

func TestNonstandardNameStruct(t *testing.T) {
	b := []byte("d3:a-ai1e3:b bi2ee")
	var s testNonstandardNameStruct

	err := Unmarshal(b, &s)
	if err != nil {
		t.Fatal(err)
	}

	if s.AA != 1 {
		t.Errorf("invalid value in struct, got %d, expected 1", s.AA)
	}

	if s.BB != 2 {
		t.Errorf("invalid value in struct, git %d, expected 2", s.BB)
	}
}
