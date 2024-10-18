package printer

import (
	"bufio"
	"encoding/json"
	"fmt"

	"github.com/LeMyst/kube-no-trouble/pkg/judge"
)

type jsonPrinter struct {
	*commonPrinter
}

// newJSONPrinter creates new JSON printer that prints to given output file
func newJSONPrinter(options *PrinterOptions) (Printer, error) {
	cp, err := newCommonPrinter(options)
	if err != nil {
		return nil, fmt.Errorf("failed to create new common printer: %w", err)
	}
	return &jsonPrinter{
		commonPrinter: cp,
	}, nil
}

// Close will free resources used by the printer
func (c *jsonPrinter) Close() error {
	return c.commonPrinter.Close()
}

// Print will print results in text format
func (c *jsonPrinter) Print(results []judge.Result) error {
	writer := bufio.NewWriter(c.commonPrinter.options.outputFile)
	defer writer.Flush()

	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "\t")

	if !c.commonPrinter.options.showLabels {
		removeLabels(results)
	}

	err := encoder.Encode(results)
	if err != nil {
		return err
	}

	return nil
}

func removeLabels(results []judge.Result) {
	for i := range results {
		if results[i].Labels != nil {
			results[i].Labels = map[string]interface{}{}
		}
	}
}
