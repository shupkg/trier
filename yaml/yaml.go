package yaml

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

var (
	Unmarshal  = yaml.Unmarshal
	Marshal    = yaml.Marshal
	NewEncoder = yaml.NewEncoder
	NewDecoder = yaml.NewDecoder
)

func ReadFromFile(path string, i interface{}) (err error) {
	if v, err := ioutil.ReadFile(path); err != nil {
		return err
	} else {
		return Unmarshal(v, i)
	}
}

func WriteToFile(path string, i interface{}) (err error) {
	if v, err := Marshal(i); err != nil {
		return err
	} else {
		return ioutil.WriteFile(path, v, 0644)
	}
}
