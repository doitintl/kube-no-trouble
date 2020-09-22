package collector

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ghodss/yaml"
	"helm.sh/helm/v3/pkg/releaseutil"
	"io/ioutil"
	"os"
	"sort"
)

type FileCollector struct {
	*commonCollector
	filenames []string
}

type FileOpts struct {
	Filenames []string
}

func NewFileCollector(opts *FileOpts) (*FileCollector, error) {

	if len(opts.Filenames) == 0 {
		return nil, errors.New("file list can't be empty")
	}

	collector := &FileCollector{
		commonCollector: &commonCollector{name: "File"},
		filenames:       opts.Filenames,
	}

	return collector, nil
}

func (c *FileCollector) Get() ([]map[string]interface{}, error) {

	var results []map[string]interface{}

	for _, f := range c.filenames {
		var input []byte
		var err error
		if f == "-" {
			input, err = ioutil.ReadAll(os.Stdin)
		} else {
			input, err = ioutil.ReadFile(f)
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read file %s: %v", f, err)
		}

		// try to parse JSON
		var manifest map[string]interface{}
		err = json.Unmarshal(input, &manifest)
		if err == nil {
			results = append(results, manifest)
		}

		// let's try YAML too
		if err != nil {
			manifests := releaseutil.SplitManifests(string(input))

			// keep output stable
			var keys []string
			for key := range manifests {
				keys = append(keys, key)
			}
			sort.Sort(releaseutil.BySplitManifestsOrder(keys))

			for _, k := range keys {
				var manifest map[string]interface{}
				err := yaml.Unmarshal([]byte(manifests[k]), &manifest)
				if err != nil {
					return nil, fmt.Errorf("failed to parse file %s: %v", f, err)
				}

				results = append(results, manifest)
			}
		}

	}

	return results, nil
}
