package printer

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/doitintl/kube-no-trouble/pkg/judge"
)

var findingsJsonTesting []judge.Result = []judge.Result{
	{
		Name:        "testName1",
		Kind:        "testKind1",
		Namespace:   "testNamespace1",
		ApiVersion:  "v1",
		RuleSet:     "testRuleset1",
		ReplaceWith: "testReplaceWith1",
		Since:       "testSince1",
	},
	{
		Name:        "testName2",
		Kind:        "testKind2",
		Namespace:   "testNamespace2",
		ApiVersion:  "v1",
		RuleSet:     "testRuleset2",
		ReplaceWith: "testReplaceWith2",
		Since:       "testSince2",
	},
}

func captureOutput(f func([]judge.Result) error, r []judge.Result) (string, error) {
	reader, writer, err := os.Pipe()
	if err != nil {
		panic(err)
	}
	stdout := os.Stdout
	stderr := os.Stderr
	defer func() {
		os.Stdout = stdout
		os.Stderr = stderr
		log.SetOutput(os.Stderr)
	}()
	os.Stdout = writer
	os.Stderr = writer
	log.SetOutput(writer)
	out := make(chan string)
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		var buf bytes.Buffer
		wg.Done()
		io.Copy(&buf, reader)
		out <- buf.String()
	}()
	wg.Wait()
	err = f(r)
	writer.Close()
	return <-out, err
}

func TestNewJsonPrinter(t *testing.T) {
	jsonPrinter, _ := NewJSONPrinter(&JSONOpts{})
	if jsonPrinter.opts.OnePerLine != false {
		t.Errorf("OnePerLine should default to false")
	}

	jsonPrinter, _ = NewJSONPrinter(&JSONOpts{OnePerLine: true})
	if jsonPrinter.opts.OnePerLine != true {
		t.Errorf("OnePerLine should default to true")
	}
}

func TestPrintJson(t *testing.T) {
	jsonPrinter, _ := NewJSONPrinter(&JSONOpts{OnePerLine: false})
	output, _ := captureOutput(jsonPrinter.Print, findingsJsonTesting)

	var j []judge.Result
	err := json.Unmarshal([]byte(output), &j)
	if err != nil {
		t.Fatal(err)
	}
	if len(j) != 2 {
		t.Error("wrong number of results")
	}
}

func TestPrintJsonL(t *testing.T) {
	x, _ := NewJSONPrinter(&JSONOpts{OnePerLine: true})
	output, _ := captureOutput(x.Print, findingsJsonTesting)

	if len(strings.Split(output, "\n")) != 3 {
		t.Errorf("execting two lines and got %d", len(strings.Split(output, "\n")))
	}

	for _, line := range strings.Split(output, "\n") {
		if len(line) > 0 {
			var j judge.Result
			err := json.Unmarshal([]byte(line), &j)
			if err != nil {
				t.Fatal(err)
			}

			if !strings.HasPrefix(j.Name, "testName") {
				t.Errorf("unexpected name")
			}
		}
	}
}
