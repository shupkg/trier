package toml

import (
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/pelletier/go-toml"
)

var (
	Unmarshal  = toml.Unmarshal
	Marshal    = toml.Marshal
	NewEncoder = toml.NewEncoder
	NewDecoder = toml.NewDecoder
	Load       = loadTree

	ErrNotFound = errors.New("key is not found")
)

type Tree struct {
	tree *toml.Tree
}

func loadTree(fn string) (*Tree, error) {
	tree, err := toml.LoadFile(fn)
	if err != nil {
		return nil, err
	}
	return &Tree{tree: tree}, nil
}

func (tree *Tree) Unmarshal(v interface{}, keys ...string) error {
	if len(keys) > 0 {
		i := tree.tree.GetPath(keys)
		if i == nil {
			return ErrNotFound
		}
		tr, ok := i.(*toml.Tree)
		if ok {
			return tr.Unmarshal(v)
		}
		return fmt.Errorf("%v", tr)
	}
	return tree.tree.Unmarshal(v)
}

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
