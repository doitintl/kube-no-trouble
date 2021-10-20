package printer

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/doitintl/kube-no-trouble/pkg/judge"
)

type textPrinter struct {
}

func newTextPrinter() Printer {
	return &textPrinter{}
}

func (c *textPrinter) Print(results []judge.Result) error {

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
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s (%s)\n", "KIND", "NAMESPACE", "NAME", "API_VERSION", "REPLACE_WITH", "SINCE")
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s (%s)\n", r.Kind, r.Namespace, r.Name, r.ApiVersion, r.ReplaceWith, r.Since)
	}
	w.Flush()
	return nil
}
