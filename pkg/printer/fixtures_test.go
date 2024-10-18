package printer

import "github.com/doitintl/kube-no-trouble/pkg/judge"

func getTestResult(labels map[string]interface{}) []judge.Result {
	version, _ := judge.NewVersion("1.2.3")

	res := []judge.Result{{
		Name:        "Name",
		Namespace:   "Namespace",
		Kind:        "Kind",
		ApiVersion:  "1.2.3",
		RuleSet:     "Test",
		ReplaceWith: "4.5.6",
		Since:       version,
		Labels:      labels,
	}}
	return res
}
