package flagx_test

import (
	"testing"

	"eagain.net/go/bearproxy/flagx"
)

func TestIPEmpty(t *testing.T) {
	var a flagx.IP
	if err := a.Set(""); err != nil {
		t.Fatalf("empty IP.Set failed: %v", err)
	}
	if a.IP != nil {
		t.Fatalf("empty IP is not nil: %v", a)
	}
}

func setIP(t testing.TB, value string) string {
	var a flagx.IP
	if err := a.Set(value); err != nil {
		t.Fatalf("IP.Set failed: %v", err)
	}
	return a.String()
}

func TestSetIPv4(t *testing.T) {
	if g, e := setIP(t, "192.0.2.42"), "192.0.2.42"; g != e {
		t.Errorf("unexpected IP: %q != %q", g, e)
	}
}

func TestSetIPv6(t *testing.T) {
	if g, e := setIP(t, "fe80::ba2a:a5ff:1260:a7f6"), "fe80::ba2a:a5ff:1260:a7f6"; g != e {
		t.Errorf("unexpected IP: %q != %q", g, e)
	}
}
