package json

import "io/ioutil"

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

