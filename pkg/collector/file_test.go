package collector

import (
	"strings"
	"testing"
)

func TestNewFileCollectorEmpty(t *testing.T) {
	input := []string{}
	expected := "empty"

	_, err := NewFileCollector(
		&FileOpts{Filenames: input},
	)

	if err == nil {
		t.Errorf("Expected error with empty file list")
	} else if !strings.Contains(err.Error(), expected) {
		t.Errorf("Expected error message with %s, got %s", expected, err.Error())
	}
}

func TestFileCollectorGet(t *testing.T) {
	testCases := []struct {
		name     string
		input    []string // file list
		expected int      // number of manifests
	}{
		{"yaml", []string{"../../fixtures/deployment-v1beta1.yaml"}, 1},
		{"yamlMulti", []string{"../../fixtures/deployment-v1beta1-and-ingress-v1beta1.yaml"}, 2},
		{"json", []string{"../../fixtures/deployment-v1beta1.json"}, 1},
		{"mixed", []string{"../../fixtures/deployment-v1beta1.json", "../../fixtures/deployment-v1beta1.yaml"}, 2},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c, err := NewFileCollector(
				&FileOpts{Filenames: tc.input},
			)

			if err != nil {
				t.Errorf("Expected to succeed for %s, failed: %s", tc.input, err)
			}

			manifests, err := c.Get()
			if err != nil {
				t.Errorf("Expected to succeed for %s, failed: %s", tc.input, err)
			} else if len(manifests) != tc.expected {
				t.Errorf("Expected to get %d, got %d", tc.expected, len(manifests))
			}
		})
	}
}

func TestFileCollectorGetUnknown(t *testing.T) {
	input := []string{"../../fixtures/meow.txt"}
	expected := "failed to parse"

	c, err := NewFileCollector(
		&FileOpts{Filenames: input},
	)

	if err != nil {
		t.Errorf("Expected to succeed for %s, failed: %s", input, err)
	}

	_, err = c.Get()

	if err == nil {
		t.Errorf("Expected error with unknown file type")
	} else if !strings.Contains(err.Error(), expected) {
		t.Errorf("Expected error message with %s, got %s", expected, err.Error())
	}
}

func TestFileCollectorGetNonExistent(t *testing.T) {
	input := []string{"../../fixtures/does-not-exist"}
	expected := "failed to read"

	c, err := NewFileCollector(
		&FileOpts{Filenames: input},
	)

	if err != nil {
		t.Errorf("Expected to succeed for %s, failed: %s", input, err)
	}

	_, err = c.Get()

	if err == nil {
		t.Errorf("Expected error with non-existent file type")
	} else if !strings.Contains(err.Error(), expected) {
		t.Errorf("Expected error message with %s, got %s", expected, err.Error())
	}
}
