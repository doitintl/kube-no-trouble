package printer

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/doitintl/kube-no-trouble/pkg/judge"
)

type JSONPrinter struct {
}

type JSONOpts struct {
}

func NewJSONPrinter(opts *JSONOpts) (*JSONPrinter, error) {
	printer := &JSONPrinter{}

	return printer, nil
}

func (c *JSONPrinter) Print(results []judge.Result) error {

	buffer := new(bytes.Buffer)
	encoder := json.NewEncoder(buffer)
	encoder.SetIndent("", "\t")

	err := encoder.Encode(results)
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", buffer.String())

	return nil
}
