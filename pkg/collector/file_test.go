package collector

import (
	"io"
	"os"
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

func TestFileCollectorGetStdin(t *testing.T) {
	input := []string{"-"}
	inputFilename := "../../fixtures/deployment-v1beta1.yaml"
	expected := 1

	c, err := NewFileCollector(
		&FileOpts{Filenames: input},
	)
	if err != nil {
		t.Errorf("Expected to succeed for %s, failed: %s", input, err)
	}

	fakeStdinReader, fakeStdinWriter, err := os.Pipe()
	if err != nil {
		t.Errorf("Failed to create pipe for Stdin redirect: %s", err)
	}

	// override os.stdin
	origStdin := os.Stdin
	os.Stdin = fakeStdinReader
	defer func() {
		os.Stdin = origStdin
		fakeStdinReader.Close()
	}()

	f, err := os.Open(inputFilename)
	if err != nil {
		t.Errorf("Failed to open fixture file %s: %s", inputFilename, err)
	}
	defer func() { f.Close() }()

	// read file to fake stdin
	_, err = io.Copy(fakeStdinWriter, f)
	if err != nil {
		t.Errorf("Failed to read fixture file: %s", err)
	}
	fakeStdinWriter.Close()

	manifests, err := c.Get()
	if err != nil {
		t.Errorf("Expected to succeed for %s, failed: %s", input, err)
	} else if len(manifests) != expected {
		t.Errorf("Expected to get %d, got %d", expected, len(manifests))
	}
}
