package printer

import (
	"errors"
	"os"

	"github.com/doitintl/kube-no-trouble/pkg/judge"
	prom "github.com/prometheus/client_golang/prometheus"

	"github.com/prometheus/client_golang/prometheus/push"
)

type PrometheusPrinter struct {
	gauge  *prom.GaugeVec
	pusher *push.Pusher
	job    string
}

func newPrometheusPrinter(outputFileName string) (Printer, error) {
	gatewayUrl := os.Getenv("PUSHGATEWAY_URL")
	job := os.Getenv("PUSHGATEWAY_JOB")

	if gatewayUrl == "" || job == "" {
		return &PrometheusPrinter{}, errors.New("pushgateway URL or job not defined")
	}

	gauge := prom.NewGaugeVec(prom.GaugeOpts{
		Name: "kubent_deprecated_objects",
		Help: "Objects which have been detected by kubent as being pinned to deprecated APIs",
	}, []string{"name", "namespace", "kind", "api_version", "rule_set", "replace_with", "since"})
	pusher := push.New(gatewayUrl, job).Collector(gauge)

	return &PrometheusPrinter{
		gauge,
		pusher,
		job,
	}, nil
}

func (p *PrometheusPrinter) Print(results []judge.Result) error {
	p.gauge.Reset()

	for _, r := range results {
		p.gauge.WithLabelValues(r.Name, r.Namespace, r.Kind, r.ApiVersion, r.RuleSet, r.ReplaceWith, r.Since.String()).Set(1)
	}

	return nil
}

func (p *PrometheusPrinter) Close() error {
	return p.pusher.Push()
}
