package printer

import (
	"context"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	ctxKey "github.com/doitintl/kube-no-trouble/pkg/context"
)

func TestTypePrinterPrint(t *testing.T) {
	tests := []struct {
		name       string
		labels     map[string]interface{}
		withLabels bool
	}{
		{
			name:       "WithLabels",
			labels:     map[string]interface{}{"app": "version1"},
			withLabels: true,
		},
		{
			name:       "NoLabels",
			labels:     map[string]interface{}{},
			withLabels: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpFile, err := ioutil.TempFile(os.TempDir(), tempFilePrefix)
			if err != nil {
				t.Fatalf(tempFileCreateFailureMessage, err)
			}
			defer os.Remove(tmpFile.Name())

			tp := &csvPrinter{
				commonPrinter: &commonPrinter{tmpFile},
			}

			results := getTestResult(tt.labels)

			ctx := context.WithValue(context.Background(), ctxKey.LABELS_CTX_KEY, &tt.withLabels)
			if err := tp.Print(results, ctx); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			fi, _ := tmpFile.Stat()
			if fi.Size() == 0 {
				t.Fatalf("expected non-zero size output file: %v", err)
			}
		})
	}
}

func TestMapToCommaSeparatedString(t *testing.T) {
	tests := []struct {
		input    map[string]interface{}
		expected string
	}{
		{
			input:    map[string]interface{}{},
			expected: "",
		},
		{
			input:    map[string]interface{}{"key1": "value1"},
			expected: "key1:value1",
		},
		{
			input:    map[string]interface{}{"key1": "value1", "key2": "value2"},
			expected: "key1:value1, key2:value2",
		},
		{
			input:    map[string]interface{}{"key1": 123, "key2": true, "key3": 45.67},
			expected: "key1:123, key2:true, key3:45.67",
		},
	}

	for _, test := range tests {
		result := mapToCommaSeparatedString(test.input)
		parsedResult := parseCSVString(result)
		parsedExpected := parseCSVString(test.expected)

		// Compare the parsed maps
		if !compareMaps(parsedResult, parsedExpected) {
			t.Errorf("For input %v, expected %s but got %s", test.input, test.expected, result)
		}
	}
}

func parseCSVString(s string) map[string]string {
	result := make(map[string]string)
	if s == "" {
		return result
	}

	pairs := strings.Split(s, ", ")
	for _, pair := range pairs {
		kv := strings.SplitN(pair, ":", 2)
		if len(kv) == 2 {
			result[kv[0]] = kv[1]
		}
	}
	return result
}

func compareMaps(a, b map[string]string) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if bv, ok := b[k]; !ok || v != bv {
			return false
		}
	}
	return true
}
