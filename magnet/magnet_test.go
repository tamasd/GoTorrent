package magnet

import "testing"

func TestBasic(t *testing.T) {
	m, err := Parse("magnet:?xt=urn:btih:1234567890123456789012345678901234567890&dn=foo&tr=bar.baz")
	if err != nil {
		t.Fatal(err)
	}

	if m.InfoHash != "1234567890123456789012345678901234567890" {
		t.Error("Failed to extract infohash")
	}

	if m.Name != "foo" {
		t.Error("Failed to extract name")
	}

	if m.Tracker != "bar.baz" {
		t.Error("Failed to extract tracker")
	}
}
