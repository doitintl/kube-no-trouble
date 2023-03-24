package printer

import (
	"context"
	"encoding/csv"
	"fmt"
	"sort"

	"github.com/LeMyst/kube-no-trouble/pkg/judge"
)

type csvPrinter struct {
	*commonPrinter
}

// newCSVPrinter creates new CSV printer that prints to given output file
func newCSVPrinter(outputFileName string) (Printer, error) {
	cp, err := newCommonPrinter(outputFileName)
	if err != nil {
		return nil, fmt.Errorf("failed to create new common printer: %w", err)
	}
	return &csvPrinter{
		commonPrinter: cp,
	}, nil
}

// Close will free resources used by the printer
func (c *csvPrinter) Close() error {
	return c.commonPrinter.Close()
}

// Print will print results in CSV format
func (c *csvPrinter) Print(results []judge.Result, ctx context.Context) error {

	sort.Slice(results, func(i, j int) bool {
		return results[i].Name < results[j].Name
	})
	sort.Slice(results, func(i, j int) bool {
		return results[i].Namespace < results[j].Namespace
	})
	sort.Slice(results, func(i, j int) bool {
		return results[i].Kind < results[j].Kind
	})
	sort.Slice(results, func(i, j int) bool {
		return results[i].RuleSet < results[j].RuleSet
	})

	w := csv.NewWriter(c.commonPrinter.outputFile)

	fields := []string{
		"api_version",
		"kind",
		"namespace",
		"name",
		"replace_with",
		"since",
		"rule_set",
	}

	labels, err := shouldShowLabels(ctx)
	if err != nil {
		return fmt.Errorf("failed to get labels: %w", err)
	}

	if labels != nil && *labels {
		fields = append(fields, "labels")
	}

	w.Write(fields)

	for _, r := range results {
		row := []string{
			r.ApiVersion,
			r.Kind,
			r.Namespace,
			r.Name,
			r.ReplaceWith,
			r.Since.String(),
			r.RuleSet,
		}

		if labels != nil && *labels {
			row = append(row, mapToCommaSeparatedString(r.Labels))
		}

		w.Write(row)
	}

	w.Flush()
	return nil
}
