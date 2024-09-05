package cfg

import (
	"log"
	"os"
	"sort"

	"github.com/ghodss/yaml"
)

type Opts struct {
	Opts map[string]Opt `json:"opts"`
}

type OptType string

const (
	OptTypeBool OptType = "bool"
	OptTypeList OptType = "list"
)

type Opt struct {
	Short string `json:"short"`
	// Allowed: bool, list
	Type         OptType `json:"type"`
	DefaultValue string  `json:"defaultValue"`
	Optional     bool    `json:"optional"`
	Desc         string  `json:"desc"`
}

func sortedOptsData(opts Opts) map[string]Opt {
	// Extract keys from map
	keys := make([]string, 0, len(opts.Opts))
	for k := range opts.Opts {
		keys = append(keys, k)
	}

	// Sort keys
	sort.Strings(keys)

	result := map[string]Opt{}
	for _, k := range keys {
		result[k] = opts.Opts[k]
	}
	return result
}

func ReadInput(filename string) Opts {
	b, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	opts := Opts{}
	err = yaml.Unmarshal(b, &opts)
	if err != nil {
		log.Fatal(err)
	}

	sortedKeys := sortedOptsData(opts)
	return Opts{
		Opts: sortedKeys,
	}
}
