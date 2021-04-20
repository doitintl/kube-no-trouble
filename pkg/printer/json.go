package printer

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/doitintl/kube-no-trouble/pkg/judge"
)

type jsonPrinter struct {
}

func newJSONPrinter() Printer {
	return &jsonPrinter{}
}

func (c *jsonPrinter) populateOutput(results []judge.Result) (string, error) {
	buffer := new(bytes.Buffer)
	encoder := json.NewEncoder(buffer)
	encoder.SetIndent("", "\t")

	err := encoder.Encode(results)
	if err != nil {
		return "", err
	}
	return buffer.String(), nil
}

func (c *jsonPrinter) Print(results []judge.Result) error {
	output, err := c.populateOutput(results)
	if err != nil {
		return err
	}
	fmt.Printf("%s", output)

	return nil
}
