package main

import (
	"encoding/json"
	"net"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestFilterIPs(t *testing.T) {
	ips := []net.IP{
		net.ParseIP("192.0.2.1"),
		net.ParseIP("2001:db8::1"),
		net.ParseIP("192.0.2.2"),
		net.ParseIP("2001:db8::2"),
	}

	v4 := filterIPs(ips, false)
	if want := []string{"192.0.2.1", "192.0.2.2"}; !equal(v4, want) {
		t.Errorf("filterIPs(v4) = %v, want %v", v4, want)
	}

	v6 := filterIPs(ips, true)
	if want := []string{"2001:db8::1", "2001:db8::2"}; !equal(v6, want) {
		t.Errorf("filterIPs(v6) = %v, want %v", v6, want)
	}
}

func TestFilterIPsEmpty(t *testing.T) {
	if out := filterIPs(nil, false); out != nil {
		t.Errorf("filterIPs(nil) = %v, want nil", out)
	}
}

func TestResultMarshalOmitsEmptyFields(t *testing.T) {
	result := Result{Host: "example.com", A: []string{"192.0.2.1"}}

	jsonOut, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("json.Marshal: %v", err)
	}
	want := `{"host":"example.com","a":["192.0.2.1"]}`
	if string(jsonOut) != want {
		t.Errorf("json.Marshal = %s, want %s", jsonOut, want)
	}

	yamlOut, err := yaml.Marshal(result)
	if err != nil {
		t.Fatalf("yaml.Marshal: %v", err)
	}
	var roundTrip Result
	if err := yaml.Unmarshal(yamlOut, &roundTrip); err != nil {
		t.Fatalf("yaml.Unmarshal: %v", err)
	}
	if roundTrip.AAAA != nil || roundTrip.CNAME != "" || roundTrip.MX != nil {
		t.Errorf("yaml round-trip had unexpected non-empty fields: %+v", roundTrip)
	}
	if len(roundTrip.A) != 1 || roundTrip.A[0] != "192.0.2.1" {
		t.Errorf("yaml round-trip A = %v, want [192.0.2.1]", roundTrip.A)
	}
}

func TestLookupHostUnknownType(t *testing.T) {
	if _, err := lookupHost("example.com", []string{"BOGUS"}); err == nil {
		t.Error("lookupHost with unknown type = nil error, want error")
	}
}

func equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
