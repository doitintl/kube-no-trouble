package judge

import (
	"context"
	"github.com/LeMyst/kube-no-trouble/pkg/rules"
	"github.com/open-policy-agent/opa/rego"
	"github.com/rs/zerolog/log"
)

type RegoJudge struct {
	preparedQuery rego.PreparedEvalQuery
}

type RegoOpts struct {
}

func NewRegoJudge(opts *RegoOpts, rules []rules.Rule) (*RegoJudge, error) {
	ctx := context.Background()

	r := rego.New(
		rego.Query("data[_].main"),
	)

	for _, info := range rules {
		rego.Module(info.Name, info.Rule)(r)
		log.Info().Str("name", info.Name).Msg("Loaded ruleset")
	}

	pq, err := r.PrepareForEval(ctx)
	if err != nil {
		return nil, err
	}

	judge := &RegoJudge{preparedQuery: pq}
	return judge, nil
}

func (j *RegoJudge) Eval(input []map[string]interface{}) ([]Result, error) {
	ctx := context.Background()

	log.Trace().Msgf("evaluating +%v", input)
	rs, err := j.preparedQuery.Eval(ctx, rego.EvalInput(input))
	if err != nil {
		return nil, err
	}

	results := []Result{}
	for _, r := range rs {
		for _, e := range r.Expressions {
			for _, i := range e.Value.([]interface{}) {
				m := i.(map[string]interface{})
				log.Trace().Msgf("parsing +%v", m)

				since, err := NewVersion(m["Since"].(string))
				if err != nil {
					log.Debug().Msgf("Failed to parse version: %s", err)
				}

				// shouldn't really happen - but if it does fix up and move on
				if m["Namespace"] == nil {
					log.Warn().Msgf("Object has invalid namespace: %s/%s %s", m["ApiVersion"].(string), m["Kind"].(string), m["Name"].(string))
					m["Namespace"] = "<undefined>"
				}

				results = append(results, Result{
					Name:        m["Name"].(string),
					Namespace:   m["Namespace"].(string),
					Kind:        m["Kind"].(string),
					ApiVersion:  m["ApiVersion"].(string),
					ReplaceWith: m["ReplaceWith"].(string),
					RuleSet:     m["RuleSet"].(string),
					Since:       since,
				})
			}
		}
	}

	return results, nil
}
