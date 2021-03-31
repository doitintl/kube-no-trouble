package printer

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/doitintl/kube-no-trouble/pkg/judge"
)

type JSONPrinter struct {
	opts *JSONOpts
}

type JSONOpts struct {
	// OnePerLine enables JSONL support
	OnePerLine bool
}

func NewJSONPrinter(opts *JSONOpts) (*JSONPrinter, error) {
	printer := &JSONPrinter{opts: opts}

	return printer, nil
}

func (c *JSONPrinter) Print(results []judge.Result) error {
	buffer := new(bytes.Buffer)

	if !c.opts.OnePerLine {
		encoder := json.NewEncoder(buffer)
		encoder.SetIndent("", "\t")

		err := encoder.Encode(results)
		if err != nil {
			return err
		}
	} else {
		for _, entry := range results {
			b, err := json.Marshal(entry)
			if err != nil {
				return err
			}

			_, err = buffer.Write(b)
			if err != nil {
				return err
			}

			err = buffer.WriteByte('\n')
			if err != nil {
				return err
			}
		}
	}

	fmt.Printf("%s", buffer.String())

	return nil
}
