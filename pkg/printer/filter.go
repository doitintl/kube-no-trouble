package printer

import (
	"github.com/doitintl/kube-no-trouble/pkg/judge"
)

func FilterNonRelevantResults(results []judge.Result, tv *judge.Version) ([]judge.Result, error) {
	if tv != nil {
		filtered := []judge.Result{}

		for i := range results {
			if results[i].Since == nil || results[i].Since.LessThanOrEqual(tv.Version) {
				filtered = append(filtered, results[i])
			}
		}

		return filtered, nil
	}

	return results, nil
}
