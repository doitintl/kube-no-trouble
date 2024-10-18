package test

import (
	"testing"

	"github.com/LeMyst/kube-no-trouble/pkg/collector"
	"github.com/LeMyst/kube-no-trouble/pkg/judge"
	"github.com/LeMyst/kube-no-trouble/pkg/rules"

	"k8s.io/apimachinery/pkg/runtime/schema"
)

func TestRegoCustom(t *testing.T) {
	manifestName := "../fixtures/issuer-v1alpha2.yaml"
	expectedKind := "Issuer"
	additionalKinds := []schema.GroupVersionKind{
		{
			Group:   "cert-manager.io",
			Version: "v1alpha2",
			Kind:    "Issuer",
		},
	}

	c, err := collector.NewFileCollector(
		&collector.FileOpts{Filenames: []string{manifestName}},
	)

	if err != nil {
		t.Errorf("expected to succeed for %s, failed: %s", manifestName, err)
	}

	manifests, err := c.Get()
	if err != nil {
		t.Errorf("expected to succeed for %s, failed: %s", manifestName, err)
	} else if len(manifests) == 0 {
		t.Errorf("expected to get some manifests, got %d", len(manifests))
	}

	loadedRules, err := rules.FetchRegoRules(additionalKinds)
	if err != nil {
		t.Errorf("failed to load rules: %s", err)
	}

	judge, err := judge.NewRegoJudge(&judge.RegoOpts{}, loadedRules)
	if err != nil {
		t.Errorf("failed to initialize judge: %s", err)
	}

	results, err := judge.Eval(manifests)
	if err != nil {
		t.Errorf("failed to evaluate rules: %s", err)
	}

	if len(results) == 0 {
		t.Errorf("expected to get some manifests, got %d", len(results))
	}

	if results[0].Kind != expectedKind {
		t.Errorf("expected to get %s result, got %s", expectedKind, results[0].Kind)
	}
}
