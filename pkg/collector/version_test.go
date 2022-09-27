package collector

import (
	"encoding/json"
	"testing"

	goversion "github.com/hashicorp/go-version"
)

const (
	defaultOutput   = "expected: %q, got: %q"
	expectedSuccess = "expected to succeed, failed instead: %v"
)

func TestVersionMarshalText(t *testing.T) {
	v := "1.2.3"

	version, err := NewVersion(v)
	if err != nil {
		t.Fatal(err)
	}

	actual, err := version.MarshalText()

	if err != nil {
		t.Fatal(err)
	}
	if string(actual) != v {
		t.Fatalf(defaultOutput, v, actual)
	}
}

func TestVersionString(t *testing.T) {
	v := "1.2.3"

	version, err := NewVersion(v)
	if err != nil {
		t.Fatal(err)
	}

	actual := version.String()

	if err != nil {
		t.Fatal(err)
	}
	if actual != v {
		t.Fatalf(defaultOutput, v, actual)
	}
}

func TestVersionStringNil(t *testing.T) {
	var version Version
	expected := ""

	actual := version.String()

	if actual != expected {
		t.Fatalf(defaultOutput, expected, actual)
	}
}

func TestNewVersion(t *testing.T) {
	expected := "1.2.3"

	v, err := NewVersion(expected)
	if err != nil {
		t.Fatal(err)
	}

	if v.String() != expected {
		t.Fatalf(defaultOutput, expected, v.String())
	}
}

func TestNewVersionEmpty(t *testing.T) {
	expected := ""

	_, err := NewVersion(expected)
	if err == nil {
		t.Fatalf("expected to fail with non sem-ver version")
	}
}

func TestVersionSet(t *testing.T) {
	expected := "1.2.3"
	v := Version{}

	err := v.Set(expected)
	if err != nil {
		t.Fatalf(expectedSuccess, err)
	}

	if v.String() != expected {
		t.Fatalf(defaultOutput, expected, v.String())
	}
}

func TestVersionNewFromGoVersion(t *testing.T) {
	expected := "1.2.3"
	goVer, err := goversion.NewVersion(expected)
	if err != nil {
		t.Fatalf(expectedSuccess, err)
	}

	v, err := NewFromGoVersion(goVer)
	if err != nil {
		t.Fatalf(expectedSuccess, err)
	}

	if v.String() != expected {
		t.Fatalf(defaultOutput, expected, v.String())
	}
}

func TestVersionUnmarshalText(t *testing.T) {
	expected := "1.2.3"
	source, err := NewVersion(expected)
	if err != nil {
		t.Fatal(err)
	}

	sourceJson, err := json.Marshal(source)
	if err != nil {
		t.Fatal(err)
	}

	v := Version{}
	if err := json.Unmarshal(sourceJson, &v); err != nil {
		t.Fatal(err)
	}

	if v.String() != expected {
		t.Fatalf(defaultOutput, expected, v.String())
	}
}
