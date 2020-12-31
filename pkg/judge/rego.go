package judge

import (
	"context"
	"io/ioutil"
	"os"

	"github.com/open-policy-agent/opa/rego"
	"github.com/rakyll/statik/fs"
	"github.com/rs/zerolog/log"

	_ "github.com/doitintl/kube-no-trouble/generated/statik"
)

type RegoJudge struct {
	preparedQuery rego.PreparedEvalQuery
}

type RegoOpts struct {
}

func NewRegoJudge(opts *RegoOpts) (*RegoJudge, error) {
	ctx := context.Background()

	r := rego.New(
		rego.Query("data[_].main"),
	)

	statikFS, err := fs.New()

	fs.Walk(statikFS, "/",
		func(path string, info os.FileInfo, err error) error {
			if !info.IsDir() {
				if err != nil {
					return err
				}
				f, err := statikFS.Open(path)
				if err != nil {
					return err
				}
				c, err := ioutil.ReadAll(f)
				if err != nil {
					return err
				}
				rego.Module(info.Name(), string(c))(r)
				log.Info().Str("name", info.Name()).Msg("Loaded ruleset")
			}
			return nil
		})

	pq, err := r.PrepareForEval(ctx)
	if err != nil {
		return nil, err
	}

	judge := &RegoJudge{preparedQuery: pq}
	return judge, nil
}

func (j *RegoJudge) Eval(input []map[string]interface{}) ([]Result, error) {
	ctx := context.Background()

	log.Debug().Msgf("evaluating +%v", input)
	rs, err := j.preparedQuery.Eval(ctx, rego.EvalInput(input))
	if err != nil {
		return nil, err
	}

	results := []Result{}
	for _, r := range rs {
		for _, e := range r.Expressions {
			for _, i := range e.Value.([]interface{}) {
				m := i.(map[string]interface{})
				results = append(results, Result{
					Name:        m["Name"].(string),
					Namespace:   m["Namespace"].(string),
					Kind:        m["Kind"].(string),
					ApiVersion:  m["ApiVersion"].(string),
					ReplaceWith: m["ReplaceWith"].(string),
					RuleSet:     m["RuleSet"].(string),
					Since:       m["Since"].(string),
				})
			}
		}
	}

	return results, nil
}
