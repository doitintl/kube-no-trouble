package printer

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/stepanstipl/kube-no-trouble/pkg/judge"
)

type TextPrinter struct {
}

type TextOpts struct {
}

func NewTextPrinter(opts *TextOpts) (*TextPrinter, error) {
	printer := &TextPrinter{}

	return printer, nil
}

func (c *TextPrinter) Print(results []judge.Result) error {

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

	ruleSet := ""
	w := tabwriter.NewWriter(os.Stdout, 10, 0, 3, ' ', 0)

	for _, r := range results {
		if ruleSet != r.RuleSet {
			ruleSet = r.RuleSet
			fmt.Fprintf(w, "%s\n", strings.Repeat("_", 90))
			fmt.Fprintf(w, ">>> %s <<<\n", ruleSet)
			fmt.Fprintf(w, "%s\n", strings.Repeat("-", 90))
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t\n", "KIND", "NAMESPACE", "NAME", "API_VERSION")
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t\n", r.Kind, r.Namespace, r.Name, r.ApiVersion)
	}
	w.Flush()
	return nil
}
