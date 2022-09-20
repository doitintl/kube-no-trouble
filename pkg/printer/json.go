package printer

import (
	"bufio"
	"encoding/json"
	"fmt"

	"github.com/doitintl/kube-no-trouble/pkg/judge"
)

type jsonPrinter struct {
	*commonPrinter
}

// newJSONPrinter creates new JSON printer that prints to given output file
func newJSONPrinter(outputFileName string) (Printer, error) {
	cp, err := newCommonPrinter(outputFileName)
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
	writer := bufio.NewWriter(c.commonPrinter.outputFile)
	defer writer.Flush()

	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "\t")

	err := encoder.Encode(results)
	if err != nil {
		return err
	}

	return nil
}
