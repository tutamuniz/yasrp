package reverseproxy

import (
	"testing"

	"github.com/tutamuniz/yasrp/minihttp/reverseproxy/configtypes"
)

func TestNewReverseProxy(t *testing.T) {
	_, err := NewReverseProxy("127.0.0.1", 8080)
	if err != nil {
		t.Errorf("Error creating ReverseProxy %s", err.Error())
	}

}

func TestRemoveLastSlash(t *testing.T) {
	paths := [][]string{
		{"/home/", "/home"},
		{"/home//////", "/home"},
		{"/", "/"},
	}

	for i, path := range paths {
		if ret := removeLastSlash(path[0]); ret != path[1] {
			t.Errorf("(index %d) -> Used %s expected %s got %s", i, path[0], path[1], ret)
		}
	}

}
func TestReverseProxyAddLocation(t *testing.T) {
	rp, err := NewReverseProxy("127.0.0.1", 8080)
	if err != nil {
		t.Errorf("Error creating ReverseProxy %s", err.Error())
	}
	locs := []configtypes.Location{
		{
			Path:   "/home",
			Target: "http://www.uol.com.br",
		},
		{
			Path:   "/stats/",
			Target: "https://www.teste.com",
		},
		{
			Path:   "/cloud/",
			Target: "HTTPs:\\wwwteste.com",
		},
	}

	rp.AddLocation(locs[0])

	if err := rp.AddLocation(locs[1]); err != nil {
		t.Errorf("%s", err.Error())
	}

	// wrong host
	if err := rp.AddLocation(locs[2]); err != nil {
		t.Errorf("%s", err.Error())
	}

}
