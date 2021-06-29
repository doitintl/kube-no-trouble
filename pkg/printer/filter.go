package printer

import (
	"github.com/doitintl/kube-no-trouble/pkg/judge"

	goversion "github.com/hashicorp/go-version"
)

func FilterNonRelevantResults(results []judge.Result, tv *goversion.Version) ([]judge.Result, error) {
	if tv != nil {
		filtered := []judge.Result{}

		for i := range results {
			if results[i].Since.LessThanOrEqual(tv) {
				filtered = append(filtered, results[i])
			}
		}

		return filtered, nil
	}

	return results, nil
}
