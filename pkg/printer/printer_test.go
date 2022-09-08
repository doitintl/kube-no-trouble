package printer

import (
	"reflect"
	"testing"
)

func TestParsePrinter(t *testing.T) {
	for k, v := range printers {
		p, err := ParsePrinter(k)
		if err != nil {
			t.Fatalf("failed to parse printer %s: %v", k, err)
		}

		if reflect.ValueOf(p).Pointer() != reflect.ValueOf(v).Pointer() {
			t.Fatalf("expected to get function %p, got %p instead", p, p)
		}
	}
}

func TestParsePrinterInvalid(t *testing.T) {
	_, err := ParsePrinter("BAD")
	if err == nil {
		t.Fatal(err)
	}
}
