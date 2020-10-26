package trier

import (
	"strings"
	"text/template"
)

type defaulter interface {
	Default()
}
type validator interface {
	Valid() error
}

func ApplyDefault(v interface{}) {
	if i, ok := v.(defaulter); ok {
		i.Default()
	}
}

func ApplyValid(v interface{}) error {
	if i, ok := v.(validator); ok {
		return i.Valid()
	}
	return nil
}

func ReplaceTemplate(tpl string, data interface{}) (string, error) {
	te, err := template.New("").Parse(tpl)
	if err != nil {
		return "", err
	}
	w := &strings.Builder{}
	if err = te.Execute(w, data); err != nil {
		return "", err
	}
	return w.String(), nil
}

func ReplaceTemplateMust(tpl string, data interface{}) string {
	s, err := ReplaceTemplate(tpl, data)
	HandleMustError(err)
	return s
}
