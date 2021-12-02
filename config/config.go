package config

import (
	"github.com/hashicorp/hcl"
	"github.com/pkg/errors"
	"io/ioutil"
)

type Configuration struct {
	Token string `hcl:"token"`
}

func LoadConfigHCL(filePath string, target interface{}) (err error) {
	defer func() { err = errors.Wrap(err, "config.LoadConfigHCL") }()
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		err = errors.Wrap(err, `config: read`)
		return
	}
	err = errors.Wrap(hcl.Unmarshal(b, target), `config: parse`)
	return
}
