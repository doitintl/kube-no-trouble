package collector

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/ghodss/yaml"
	"helm.sh/helm/v3/pkg/releaseutil"
)

type FileCollector struct {
	*commonCollector
	filenames []string
}

type FileOpts struct {
	Filenames []string
}

func NewFileCollector(opts *FileOpts) (*FileCollector, error) {
	collector := &FileCollector{
		commonCollector: &commonCollector{name: "File"},
		filenames:       opts.Filenames,
	}

	return collector, nil
}

func (c *FileCollector) Get() ([]interface{}, error) {

	var manifest map[string]interface{}
	var results []interface{}

	for _, f := range c.filenames {
		input, err := ioutil.ReadFile(f)
		if err != nil {
			return nil, fmt.Errorf("failed to read file %s: %v", f, err)
		}

		// try to parse JSON
		err = json.Unmarshal(input, &manifest)
		if err == nil {
			results = append(results, manifest)
		}

		// let's try YAML too
		if err != nil {
			manifests := releaseutil.SplitManifests(string(input))
			for _, m := range manifests {
				err := yaml.Unmarshal([]byte(m), &manifest)
				if err != nil {
					return nil, fmt.Errorf("failed to parse file %s: %v", f, err)
				}

				results = append(results, manifest)
			}
		}

	}

	return results, nil
}
