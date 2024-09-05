package cfg

import (
	"log"
	"os"

	"github.com/ghodss/yaml"
)

type Opts struct {
	Opts []Opt `json:"opts"`
}

type OptType string

const (
	OptTypeBool OptType = "bool"
	OptTypeList OptType = "list"
)

type Opt struct {
	Name  string `json:"name"`
	Short string `json:"short"`
	// Allowed: bool, list
	Type         OptType `json:"type"`
	DefaultValue string  `json:"defaultValue"`
	Optional     bool    `json:"optional"`
	Desc         string  `json:"desc"`
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

	return opts
}
