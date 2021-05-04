package printer

import (
	"testing"
)

func TestParsePrinters(t *testing.T) {
	for k, v := range printers {
		printer, err := NewPrinter(k)
		if err != nil {
			t.Fatal(err)
		}
		if printer != v() {
			t.Fatal("Should be a valid printer")
		}
	}
}

func TestInvalidStringForParsePrinters(t *testing.T) {
	_, err := NewPrinter("BAD")
	if err == nil {
		t.Fatal(err)
	}
}
